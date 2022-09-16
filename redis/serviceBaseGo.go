package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v5"
)

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

type fn func(paramIn string) (out interface{}, BackTo string, err error)

func RedisServe(_rds *redis.Client, _serviceName string, f fn) func() {
	rds := _rds
	var batch_size int64 = 128
	serviceName := _serviceName
	cnt := CounterResetEveryMinute()

	ProcessOneJob := func(s string) (err error) {
		var (
			backTo         string
			out            interface{}
			marshaledBytes []byte
		)
		//process one job
		if out, backTo, err = f(s); err != nil {
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
				go ProcessOneJob(s)
			}
			if log, num := cnt(); log {
				fmt.Print(time.Now().String()[11:19] + " service " + serviceName + " rcved " + strconv.Itoa(num) + " tasks")
			}
		}
	}
	go loop()
	return func() {
		fmt.Print("service " + serviceName + "is running")
	}
}
