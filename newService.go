package saavuu

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/config"
	. "github.com/yangkequn/saavuu/redisContext"
)

type fn func(rc *RedisContext, paramIn map[string]interface{}) (out map[string]interface{}, err error)

var ServiceMap map[string]fn = map[string]fn{}

var counter Counter = Counter{}

func PrintServiceStates() {
	// all keys of ServiceMap to []string serviceNames
	var serviceNames []string = make([]string, 0, len(ServiceMap))
	for k := range ServiceMap {
		serviceNames = append(serviceNames, k)
	}
	fmt.Println("ServiceMap has", len(ServiceMap), "services:", serviceNames)
	for true {
		time.Sleep(time.Second * 60)
		now := time.Now().String()[11:19]
		for _, serviceName := range serviceNames {
			num, _ := counter.Get(serviceName)
			fmt.Println(now + " service " + serviceName + " proccessed " + strconv.Itoa(int(num)) + " tasks")
			counter.DeleteAndGetLastValue(serviceName)
		}
	}
}

var ErrBackTo = fmt.Errorf("param[\"backTo\"] is not a string")

func NewService(_serviceName string, f fn) {
	//check configureation is loaded
	if config.DataRds == nil {
		panic("config.DataRedis is nil. you should call config.LoadConfigFromRedis first")
	}
	if config.ParamRds == nil {
		panic("config.ParamRedis is nil. you should call config.LoadConfigFromRedis first")
	}

	ServiceMap[_serviceName] = f
	counter.DeleteAndGetLastValue(_serviceName)

	var batch_size int64 = 128
	serviceName := _serviceName
	ProcessOneJob := func(s []byte) (err error) {
		var (
			BackTo         string
			out            interface{}
			marshaledBytes []byte
			param          map[string]interface{} = map[string]interface{}{}
			ok             bool
		)
		if err = msgpack.Unmarshal(s, &param); err != nil || param["BackTo"] == nil {
			return err
		}
		if BackTo, ok = param["BackTo"].(string); !ok {
			return ErrBackTo
		}
		delete(param, "BackTo")
		//process one job
		if out, err = f(&RedisContext{Ctx: context.Background(), DataRds: config.DataRds}, param); err != nil {
			return err
		}
		//Post Back
		if marshaledBytes, err = msgpack.Marshal(out); err != nil {
			return err
		}
		ctx := context.Background()
		pipline := config.ParamRds.Pipeline()
		pipline.RPush(ctx, BackTo, marshaledBytes)
		pipline.Expire(ctx, BackTo, 6)
		_, err = pipline.Exec(ctx)
		return err
	}
	loop := func() {
		var data []string
		for true {
			//fetch datas from redis
			c := context.Background()
			pipline := config.ParamRds.Pipeline()
			pipline.LRange(c, serviceName, 0, batch_size-1)
			pipline.LTrim(c, serviceName, batch_size, -1)
			cmd, err := pipline.Exec(c)
			if err != nil || len(cmd) < 2 {
				rlt := config.ParamRds.BLPop(c, time.Minute, serviceName)
				if rlt.Err() != nil || len(rlt.Val()) == 0 {
					time.Sleep(time.Millisecond * 10)
					continue
				}
				data = rlt.Val()
			} else {
				data = cmd[0].(*redis.StringSliceCmd).Val()
			}
			for _, s := range data {
				go ProcessOneJob([]byte(s))
				counter.Add(_serviceName, 1)
			}
		}
	}
	go loop()
}
