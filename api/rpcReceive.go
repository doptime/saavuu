package api

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/config"
)

// ensure all apis can be called by rpc
// because rpc receive needs to know all api names to create stream reading
var ApiStartingWaiter func() = func() func() {
	LastApiCnt := -1
	//if the count of apis is not changing, then all apis are loaded
	Checker := func() {
		//if ApiServices.Count() no longer changed, then all apis are loaded
		for _cnt := ApiServices.Count(); _cnt == 0 || LastApiCnt != _cnt; _cnt = ApiServices.Count() {
			time.Sleep(time.Millisecond * 30)
			LastApiCnt = _cnt
		}
	}
	return Checker
}()

func rpcReceive() {
	var (
		rds      *redis.Client
		services []string
		ok       bool
	)
	for _, dataSource := range APIGroupByDataSourceName.Keys() {
		if services, ok = APIGroupByDataSourceName.Get(dataSource); !ok {
			log.Error().Str("dataSource", dataSource).Msg("dataSource not found")
			continue
		}

		if rds, ok = config.Rds[dataSource]; !ok {
			log.Error().Str("dataSource", dataSource).Msg("dataSource not found")
			continue
		}
		go rpcReceiveOneDatasource(services, rds)
	}
}
func rpcReceiveOneDatasource(serviceNames []string, rds *redis.Client) {
	var (
		apiName, data string
		cmd           *redis.XStreamSliceCmd
	)

	//wait for all rpc services ready, so that rpc results can be received
	ApiStartingWaiter()

	c := context.Background()
	XGroupEnsureCreated(c, serviceNames, rds)

	//deprecate using list command LRange, to avoid continually query consumption
	//use xreadgroup to receive data ,2023-01-31
	for args := defaultXReadGroupArgs(serviceNames); ; {
		if cmd = rds.XReadGroup(c, args); cmd.Err() == redis.Nil {
			continue
		} else if cmd.Err() != nil {
			time.Sleep(time.Second)
			log.Error().AnErr("rpcReceive", cmd.Err()).Send()
		}

		for _, stream := range cmd.Val() {
			apiName = stream.Stream
			for _, message := range stream.Messages {
				timeAtStr, atOk := message.Values["timeAt"]
				//skip case of placeholder stream while not atOk
				//but if timeAt is setted, then empty data is allowed, used to clear the task
				if data = message.Values["data"].(string); len(data) == 0 && !atOk {
					continue
				}
				//the delay calling will lost if the app is down
				if atOk {
					if len(data) == 0 {
						rpcCallAtTaskRemoveOne(apiName, timeAtStr.(string))
					} else {
						rpcCallAtTaskAddOne(apiName, timeAtStr.(string), data)
					}
				} else {
					go CallApiLocallyAndSendBackResult(apiName, message.ID, []byte(data))
				}
				apiCounter.Add(apiName, 1)
			}
		}
	}
}
func CallApiLocallyAndSendBackResult(apiName, BackToID string, s []byte) (err error) {
	var (
		msgPackResult []byte
		ret           interface{}
		service       *ApiInfo
		ok            bool
		rds           *redis.Client
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
	if rds, ok = config.Rds[service.DataSource]; !ok {
		return fmt.Errorf("DataSourceName not defined in enviroment %s", service.DataSource)
	}
	pipline := rds.Pipeline()
	pipline.RPush(ctx, BackToID, msgPackResult)
	pipline.Expire(ctx, BackToID, time.Second*20)
	_, err = pipline.Exec(ctx)
	return err
}
