package api

import (
	"context"
	"fmt"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/config"
)

type ApiInfo struct {
	// ApiName is the name of the service
	ApiName string
	// ApiFuncWithMsgpackedParam is the function of the service
	ApiFuncWithMsgpackedParam func(s []byte) (ret interface{}, err error)
}

var ApiServices cmap.ConcurrentMap[string, *ApiInfo] = cmap.New[*ApiInfo]()

func apiServiceNames() (serviceNames []string) {
	for _, serviceInfo := range ApiServices.Items() {
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

	//ServiceBatchSize is the number of tasks that a service can read from redis at the same time
	args := &redis.XReadGroupArgs{Streams: streams, Block: time.Second * 20, Count: config.Cfg.ServiceBatchSize, NoAck: true, Group: "group0", Consumer: "saavuu"}
	return args
}
func XGroupCreate(c context.Context) (err error) {
	//if there is no group, create a group, and create a consumer
	for _, serviceName := range apiServiceNames() {
		//continue if the group already exists
		if cmd := config.Rds.XInfoGroups(c, serviceName); cmd.Err() == nil || len(cmd.Val()) > 0 {
			continue
		}
		//create a group if none exists
		if cmd := config.Rds.XGroupCreateMkStream(c, serviceName, "group0", "$"); cmd.Err() != nil {
			return cmd.Err()
		}
	}
	return nil
}

func receiveJobs() {
	var (
		cmd     *redis.XStreamSliceCmd
		apiName string
	)
	c := context.Background()
	//create group if none exists
	for err := XGroupCreate(c); err != nil; err = XGroupCreate(c) {
		log.Info().Str("receiveApiJobs error:", err.Error()).Send()
		log.Info().AnErr("XGroupCreate", cmd.Err()).Send()
		time.Sleep(time.Second)
	}

	//deprecate using list command LRange, to avoid continually query consumption
	//use xreadgroup to receive data ,2023-01-31
	for args := defaultXReadGroupArgs(); ; {
		if cmd = config.Rds.XReadGroup(c, args); cmd.Err() == redis.Nil {
			continue
		} else if cmd.Err() != nil {
			log.Info().AnErr("receiveApiJobs", cmd.Err()).Send()
			time.Sleep(time.Second)
			continue
		}

		for _, stream := range cmd.Val() {
			apiName = stream.Stream
			for _, message := range stream.Messages {
				bytesValue := message.Values["data"].(string)
				//the delay calling will lost if the app is down
				if dueTimeStr, ok := message.Values["dueTime"]; ok {
					go delayTaskAddOne(apiName, dueTimeStr.(string), bytesValue)
				} else {
					go DoOneJob(apiName, apiName, []byte(bytesValue))
				}
				apiCounter.Add(apiName, 1)
			}
		}
	}
}
func DoOneJob(apiName, BackToID string, s []byte) (err error) {
	var (
		msgPackResult []byte
		ret           interface{}
		service       *ApiInfo
		ok            bool
	)
	if service, ok = ApiServices.Get(apiName); !ok {
		return fmt.Errorf("service %s not found", apiName)
	}
	if ret, err = service.ApiFuncWithMsgpackedParam(s); err != nil {
		return err
	}
	if msgPackResult, err = msgpack.Marshal(ret); err != nil {
		return
	}
	ctx := context.Background()
	pipline := config.Rds.Pipeline()
	pipline.RPush(ctx, BackToID, msgPackResult)
	pipline.Expire(ctx, BackToID, time.Second*20)
	_, err = pipline.Exec(ctx)
	return err
}
