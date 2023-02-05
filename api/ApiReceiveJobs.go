package api

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/logger"
)

type ApiInfo struct {
	// ApiName is the name of the service
	ApiName string
	// ApiFunc is the function of the service
	ApiFunc func(backTo string, s []byte) (err error)
}

var apiServices map[string]*ApiInfo = map[string]*ApiInfo{}

func apiServiceNames() (serviceNames []string) {
	for _, serviceInfo := range apiServices {
		serviceNames = append(serviceNames, serviceInfo.ApiName)
	}
	return serviceNames
}
func defaultXReadGroupArgs() *redis.XReadGroupArgs {
	var streams []string
	services := apiServiceNames()
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
	for _, serviceName := range apiServiceNames() {
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

func apiReceiveJobs() {
	var (
		cmd     *redis.XStreamSliceCmd
		apiName string
	)
	c := context.Background()
	//create group if none exists
	for err := XGroupCreate(c); err != nil; err = XGroupCreate(c) {
		logger.Lshortfile.Println("receiveApiJobs error:", err)
		time.Sleep(time.Second)
	}

	//deprecate using list command LRange, to avoid continually query consumption
	//use xreadgroup to receive data ,2023-01-31
	for args := defaultXReadGroupArgs(); ; {
		if cmd = config.ParamRds.XReadGroup(c, args); cmd.Err() == redis.Nil {
			continue
		} else if cmd.Err() != nil {
			logger.Lshortfile.Println("receiveApiJobs error:", cmd.Err())
			time.Sleep(time.Second)
			continue
		}

		for _, stream := range cmd.Val() {
			apiName = stream.Stream
			for _, message := range stream.Messages {
				bytesValue := message.Values["data"].(string)
				//the delay calling will lost if the app is down
				if dueTimeStr, ok := message.Values["dueTime"]; ok {
					go PendingApiAddOne(apiName, dueTimeStr.(string), bytesValue)
				} else {
					go apiServices[apiName].ApiFunc(message.ID, []byte(bytesValue))
				}
				apiCounter.Add(apiName, 1)
			}
		}
	}
}
