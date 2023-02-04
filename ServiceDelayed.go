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
		cmd                                        *redis.StringCmd
		bytesValue                                 string
		dueTimeUnixMilliSecond, nowUnixMilliSecond int64
		err                                        error
	)
	nowUnixMilliSecond = time.Now().UnixMilli()
	if dueTimeUnixMilliSecond, err = strconv.ParseInt(dueTimeStr, 10, 64); err != nil {
		logger.Lshortfile.Println(err)
		return
	}
	time.Sleep(time.Duration(dueTimeUnixMilliSecond-nowUnixMilliSecond) * time.Millisecond)
	if cmd = config.ParamRds.HGet(context.Background(), serviceName+":delayed", dueTimeStr); cmd.Err() != nil {
		logger.Lshortfile.Println(cmd.Err())
		return
	}
	bytesValue = cmd.Val()
	config.ParamRds.HDel(context.Background(), serviceName+":delayed", dueTimeStr)
	services[serviceName].ServiceFunc(serviceName, []byte(bytesValue))
}

func LoadDelayedServiceTask() {
	var (
		services    = serviceNames()
		dueTimeStrs []string
		err         error
	)
	pc := NewParamContext(context.Background())
	for _, service := range services {
		if dueTimeStrs, err = pc.HKeys(service); err != nil {
			continue
		}
		for _, dueTimeStr := range dueTimeStrs {
			go DelayedServiceStartOne(service, dueTimeStr)
		}
	}
}
