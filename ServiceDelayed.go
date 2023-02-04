package saavuu

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/logger"
)

// put parameter to redis ,make it persistent
func DelayedServiceEnque(serviceName string, dueTimeStr string, bytesValue string) {
	if cmd := config.ParamRds.HSet(context.Background(), serviceName+":delayed", dueTimeStr, bytesValue); cmd.Err() != nil {
		logger.Lshortfile.Println(cmd.Err())
		return
	}
	go DelayedServiceStartOne(serviceName, dueTimeStr)
}
func DelayedServiceStartOne(serviceName, dueTimeStr string) {
	var (
		bytesValue                                 string
		dueTimeUnixMilliSecond, nowUnixMilliSecond int64
		err                                        error
		cmd                                        []redis.Cmder
	)
	nowUnixMilliSecond = time.Now().UnixMilli()
	if dueTimeUnixMilliSecond, err = strconv.ParseInt(dueTimeStr, 10, 64); err != nil {
		logger.Lshortfile.Println(err)
		return
	}
	time.Sleep(time.Duration(dueTimeUnixMilliSecond-nowUnixMilliSecond) * time.Millisecond)
	pipeline := config.ParamRds.Pipeline()
	pipeline.HGet(context.Background(), serviceName+":delayed", dueTimeStr)
	pipeline.HDel(context.Background(), serviceName+":delayed", dueTimeStr)
	if cmd, err = pipeline.Exec(context.Background()); err != nil {
		logger.Lshortfile.Println(err)
		return
	}
	bytesValue = cmd[0].(*redis.StringCmd).Val()
	services[serviceName].ServiceFunc(serviceName, []byte(bytesValue))
}

func LoadDelayedServiceTask() {
	var (
		services    = serviceNames()
		dueTimeStrs []string
		cmd         []redis.Cmder
		err         error
	)
	pipeline := config.ParamRds.Pipeline()
	for _, service := range services {
		pipeline.HKeys(context.Background(), service+":delayed")
	}
	if cmd, err = pipeline.Exec(context.Background()); err != nil {
		logger.Lshortfile.Println("err LoadDelayedServiceTask, ", err)
		return
	}
	for _, service := range services {
		dueTimeStrs = cmd[i].(*redis.StringSliceCmd).Val()
		for _, dueTimeStr := range dueTimeStrs {
			go DelayedServiceStartOne(service, dueTimeStr)
		}
	}
}
