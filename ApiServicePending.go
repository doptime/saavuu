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
func PendingApiAddOne(serviceName string, dueTimeStr string, bytesValue string) {
	if cmd := config.ParamRds.HSet(context.Background(), serviceName+":pending", dueTimeStr, bytesValue); cmd.Err() != nil {
		logger.Lshortfile.Println(cmd.Err())
		return
	}
	go PendingApiRunOne(serviceName, dueTimeStr)
}
func PendingApiRunOne(serviceName, dueTimeStr string) {
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
	pipeline.HGet(context.Background(), serviceName+":pending", dueTimeStr)
	pipeline.HDel(context.Background(), serviceName+":pending", dueTimeStr)
	if cmd, err = pipeline.Exec(context.Background()); err != nil {
		logger.Lshortfile.Println(err)
		return
	}
	bytesValue = cmd[0].(*redis.StringCmd).Val()
	apiServices[serviceName].ApiFunc(serviceName, []byte(bytesValue))
}

func PendingApiFromRedisToLoal() {
	var (
		services    = apiServiceNames()
		dueTimeStrs []string
		cmd         []redis.Cmder
		err         error
	)
	pipeline := config.ParamRds.Pipeline()
	for _, service := range services {
		pipeline.HKeys(context.Background(), service+":pending")
	}
	if cmd, err = pipeline.Exec(context.Background()); err != nil {
		logger.Lshortfile.Println("err LoadPendingApiTask, ", err)
		return
	}
	for i, service := range services {
		dueTimeStrs = cmd[i].(*redis.StringSliceCmd).Val()
		for _, dueTimeStr := range dueTimeStrs {
			go PendingApiRunOne(service, dueTimeStr)
		}
	}
}
