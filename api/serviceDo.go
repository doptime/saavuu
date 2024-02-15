package api

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/specification"
)

func receiveJobs() {
	var (
		cmd     *redis.XStreamSliceCmd
		apiName string
		err     error
		strs    []string
		rds     *redis.Client = config.RdsDefaultClient()
	)
	c := context.Background()
	//create group if none exists, with consumer saavuu
	for _, serviceName := range apiServiceNames() {
		if err = XGroupCreateOne(c, serviceName); err != nil {
			time.Sleep(time.Second)
		}
	}

	//deprecate using list command LRange, to avoid continually query consumption
	//use xreadgroup to receive data ,2023-01-31
	for args := defaultXReadGroupArgs(); ; {
		if cmd = rds.XReadGroup(c, args); cmd.Err() == redis.Nil {
			continue
		} else if cmd.Err() != nil {
			time.Sleep(time.Second)
			log.Info().AnErr("receiveApiJobs", cmd.Err()).Send()
			//2023-10-18T05:39:41Z INF receiveApiJobs=NOGROUP No such key 'api:skillSearch' or consumer group 'group0' in XREADGROUP with GROUP option
			if strs = strings.Split(cmd.Err().Error(), "NOGROUP No such key '"); len(strs) < 2 {
				continue
			}
			if apiName = strings.Split(strs[1], "'")[0]; len(apiName) > 0 {
				if cmd := rds.Del(c, apiName); cmd.Err() == nil {
					log.Info().Str("Recreate group completed", apiName).Send()
					XGroupCreateOne(c, apiName)
				} else {
					log.Info().AnErr("Recreate group err", cmd.Err()).Send()
				}
			}
		}

		for _, stream := range cmd.Val() {
			apiName = stream.Stream
			for _, message := range stream.Messages {
				bytesValue := message.Values["data"].(string)
				//the delay calling will lost if the app is down
				if timeAtStr, ok := message.Values["timeAt"]; ok {
					go delayTaskAddOne(apiName, timeAtStr.(string), bytesValue)
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
		rds           *redis.Client = config.RdsDefaultClient()
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
	pipline := rds.Pipeline()
	pipline.RPush(ctx, BackToID, msgPackResult)
	pipline.Expire(ctx, BackToID, time.Second*20)
	_, err = pipline.Exec(ctx)
	return err
}

func CallByHTTP(ServiceName string, paramIn map[string]interface{}) (ret interface{}, err error) {
	var (
		fuc *ApiInfo
		ok  bool
		buf []byte
	)
	if ServiceName = specification.ApiName(ServiceName); len(ServiceName) == 0 {
		return nil, fmt.Errorf("service misnamed %s", ServiceName)
	}
	var rpc = Rpc[interface{}, interface{}](OpName(ServiceName))
	//if function is stored locally, call it directly. This is alias monolithic mode
	if fuc, ok = ApiServices.Get(ServiceName); ok {
		if buf, err = specification.MarshalApiInput(paramIn); err != nil {
			return nil, err
		}
		return fuc.ApiFuncWithMsgpackedParam(buf)
	} else {
		//if function is not stored locally, call it remotely (RPC). This is alias microservice mode
		return rpc(paramIn)
	}

}
