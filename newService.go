package saavuu

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v5"
)

type fn func(paramIn map[string]interface{}) (out map[string]interface{}, err error)

var ServiceMap map[string]fn = map[string]fn{}

func PrintServicesNames() {
	// all keys of ServiceMap to []string serviceNames
	var serviceNames []string = make([]string, 0, len(ServiceMap))
	for k := range ServiceMap {
		serviceNames = append(serviceNames, k)
	}

	fmt.Println("ServiceMap has", len(ServiceMap), "services:", serviceNames)
}

func CounterResetEveryMinute() func() (newMinute bool, counter int) {
	var minute int = time.Now().Minute()
	var _counter int = 0
	return func() (newMinute bool, counter int) {
		now := time.Now()
		if newMinute = (now.Minute() != minute); newMinute {
			minute = now.Minute()
			_counter = 0
		}
		_counter++
		return newMinute, _counter
	}
}

var ErrBackTo = fmt.Errorf("param[\"backTo\"] is not a string")

func NewService(_rds *redis.Client, _serviceName string, f fn) {
	ServiceMap[_serviceName] = f
	rds := _rds
	var batch_size int64 = 128
	serviceName := _serviceName
	cnt := CounterResetEveryMinute()

	ProcessOneJob := func(s []byte) (err error) {
		var (
			backTo         string
			out            interface{}
			marshaledBytes []byte
			param          map[string]interface{} = map[string]interface{}{}
			ok             bool
		)
		if err = msgpack.Unmarshal(s, &param); err != nil || param["backTo"] == nil {
			return err
		}
		if backTo, ok = param["backTo"].(string); !ok {
			return ErrBackTo
		}
		delete(param, "backTo")
		//process one job
		if out, err = f(param); err != nil {
			return err
		}
		//Post Back
		if marshaledBytes, err = msgpack.Marshal(out); err != nil {
			return err
		}
		ctx := context.Background()
		pipline := rds.Pipeline()
		pipline.RPush(ctx, backTo, marshaledBytes)
		pipline.Expire(ctx, backTo, 6)
		_, err = pipline.Exec(ctx)
		return err
	}
	loop := func() {
		var data []string
		for true {
			//fetch datas from redis
			c := context.Background()
			pipline := rds.Pipeline()
			pipline.LRange(c, serviceName, 0, batch_size-1)
			pipline.LTrim(c, serviceName, batch_size, -1)
			cmd, err := pipline.Exec(c)
			if err != nil || len(cmd) < 2 {
				rlt := rds.BLPop(c, time.Minute, serviceName)
				if rlt.Err() != nil || len(rlt.Val()) == 0 {
					continue
				}
				data = rlt.Val()
			} else {
				data = cmd[0].(*redis.StringSliceCmd).Val()
			}
			for _, s := range data {
				go ProcessOneJob([]byte(s))
			}
			if log, num := cnt(); log {
				fmt.Print(time.Now().String()[11:19] + " service " + serviceName + " rcved " + strconv.Itoa(num) + " tasks")
			}
		}
	}
	go loop()
}
