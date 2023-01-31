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
	// ServiceBatchSize is the number of data to be fetched from redis at one time
	ServiceBatchSize int64
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
	args := &redis.XReadGroupArgs{Streams: streams, Block: time.Second * 20, Count: 256, NoAck: true, Group: "group0", Consumer: "saavuu"}
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

func receiveServiceTask() {
	var (
		cmd *redis.XStreamSliceCmd
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
			serviceName := stream.Stream
			for _, message := range stream.Messages {
				bytesValue := message.Values["data"].(string)
				go services[serviceName].ServiceFunc(message.ID, []byte(bytesValue))
				serviceCounter.Add(serviceName, 1)
			}
		}
	}
}
