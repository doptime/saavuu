package api

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/config"
)

// put parameter to redis ,make it persistent
func delayTaskAddOne(serviceName string, dueTimeStr string, bytesValue string) {
	var (
		rds *redis.Client = config.RdsClientDefault()
	)
	if cmd := rds.HSet(context.Background(), serviceName+":delay", dueTimeStr, bytesValue); cmd.Err() != nil {
		log.Info().Err(cmd.Err()).Send()
		return
	}
	go delayTaskDoOne(serviceName, dueTimeStr)
}
func delayTaskDoOne(serviceName, dueTimeStr string) {
	var (
		ret                                        interface{}
		bytes, msgPackResult                       []byte
		dueTimeUnixMilliSecond, nowUnixMilliSecond int64
		err                                        error
		cmd                                        []redis.Cmder
		service                                    *ApiInfo
		ok                                         bool
		rds                                        *redis.Client = config.RdsClientDefault()
	)
	nowUnixMilliSecond = time.Now().UnixMilli()
	if dueTimeUnixMilliSecond, err = strconv.ParseInt(dueTimeStr, 10, 64); err != nil {
		log.Info().Err(err).Send()
		return
	}
	time.Sleep(time.Duration(dueTimeUnixMilliSecond-nowUnixMilliSecond) * time.Millisecond)
	pipeline := rds.Pipeline()
	pipeline.HGet(context.Background(), serviceName+":delay", dueTimeStr)
	pipeline.HDel(context.Background(), serviceName+":delay", dueTimeStr)
	if cmd, err = pipeline.Exec(context.Background()); err != nil {
		log.Info().Err(err).Send()
		return
	}
	if bytes, err = cmd[0].(*redis.StringCmd).Bytes(); err == nil {
		if service, ok = ApiServices.Get(serviceName); !ok {
			log.Info().Err(err).Send()
			return
		}
		if ret, err = service.ApiFuncWithMsgpackedParam(bytes); err != nil {
			log.Info().Err(err).Send()
			return
		}
		if msgPackResult, err = msgpack.Marshal(ret); err != nil {
			return
		}

		//Post Back
		ctx := context.Background()
		pipline := rds.Pipeline()
		var BackToID = serviceName
		pipline.RPush(ctx, BackToID, msgPackResult)
		pipline.Expire(ctx, BackToID, time.Second*20)
		pipline.Exec(ctx)
	}
}

func delayTasksLoad() {
	var (
		services    = apiServiceNames()
		dueTimeStrs []string
		cmd         []redis.Cmder
		err         error
		rds         *redis.Client = config.RdsClientDefault()
	)
	log.Info().Msg("delayTasksLoading started")
	pipeline := rds.Pipeline()
	for _, service := range services {
		pipeline.HKeys(context.Background(), service+":delay")
	}
	if cmd, err = pipeline.Exec(context.Background()); err != nil {
		log.Info().AnErr("err LoadDelayApiTask, ", err).Send()
		return
	}
	for i, service := range services {
		dueTimeStrs = cmd[i].(*redis.StringSliceCmd).Val()
		for _, dueTimeStr := range dueTimeStrs {
			go delayTaskDoOne(service, dueTimeStr)
		}
	}
	log.Info().Msg("delayTasksLoading completed")
}
