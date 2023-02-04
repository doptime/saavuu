package saavuu

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/logger"
)

type ServiceInfo struct {
	// ServiceName is the name of the service
	ServiceName string
	// ServiceFunc is the function of the service
	ServiceFunc func(backTo string, s []byte) (err error)
}

var services map[string]*ServiceInfo = map[string]*ServiceInfo{}

func serviceNames() (serviceNames []string) {
	for _, serviceInfo := range services {
		serviceNames = append(serviceNames, serviceInfo.ServiceName)
	}
	return serviceNames
}
func defaultXReadGroupArgs() *redis.XReadGroupArgs {
	var streams []string
	services := serviceNames()
	streams = append(streams, services...)
	//from services to ServiceInfos
	for i := 0; i < len(services); i++ {
		//append default stream id
		streams = append(streams, ">")
	}
	args := &redis.XReadGroupArgs{Streams: streams, Block: time.Second * 20, Count: config.Cfg.ServiceBatchSize, NoAck: true, Group: "group0", Consumer: "saavuu"}
	return args
}
func XGroupCreate(c context.Context) (err error) {
	//if there is no group, create a group, and create a consumer
	for _, serviceName := range serviceNames() {
		//continue if the group already exists
		if cmd := config.ParamRds.XInfoGroups(c, serviceName); cmd.Err() == nil || len(cmd.Val()) > 0 {
			continue
		}
		//create a group if none exists
		if cmd := config.ParamRds.XGroupCreateMkStream(c, serviceName, "group0", "$"); cmd.Err() != nil {
			return cmd.Err()
		}
	}
	return nil
}
func StartWidthDelay(delay time.Duration, serviceName string, ID string, bytesValue string) {
	time.Sleep(delay)
	services[serviceName].ServiceFunc(ID, []byte(bytesValue))
}
func receiveServiceTask() {
	var (
		cmd         *redis.XStreamSliceCmd
		serviceName string
		delay       time.Duration
		err         error
	)
	c := context.Background()
	//create group if none exists
	for err := XGroupCreate(c); err != nil; err = XGroupCreate(c) {
		logger.Lshortfile.Println("pipingServiceTask error:", err)
		time.Sleep(time.Second)
	}

	//deprecate using list command LRange, to avoid continually query consumption
	//use xreadgroup to receive data ,2023-01-31
	for args := defaultXReadGroupArgs(); ; {
		if cmd = config.ParamRds.XReadGroup(c, args); cmd.Err() == redis.Nil {
			continue
		} else if cmd.Err() != nil {
			logger.Lshortfile.Println("pipingServiceTask error:", cmd.Err())
			time.Sleep(time.Second)
			continue
		}

		for _, stream := range cmd.Val() {
			serviceName = stream.Stream
			for _, message := range stream.Messages {
				bytesValue := message.Values["data"].(string)
				//the delay calling will lost if the app is down
				if delayStr, ok := message.Values["delay"]; ok {
					//print message infomation
					//logger.Lshortfile.Println("pipingServiceTask delay:", delayStr.(string), "serviceName:", serviceName, "messageID:", message.ID)
					if delay, err = time.ParseDuration(delayStr.(string)); err != nil {
						logger.Lshortfile.Println("pipingServiceTask error:", err)
						break
					}
					go StartWidthDelay(delay, serviceName, message.ID, bytesValue)
				} else {
					go services[serviceName].ServiceFunc(message.ID, []byte(bytesValue))
				}
				serviceCounter.Add(serviceName, 1)
			}
		}
	}
}
