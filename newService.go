package saavuu

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/logger"
	"github.com/yangkequn/saavuu/rCtx"
)

type fn func(dc *rCtx.DataCtx, pc *rCtx.ParamCtx, paramIn map[string]interface{}) (out map[string]interface{}, err error)

var ServiceMap map[string]fn = map[string]fn{}

var counter Counter = Counter{}

func PrintServiceStates() {
	// all keys of ServiceMap to []string serviceNames
	var serviceNames []string = make([]string, 0, len(ServiceMap))
	for k := range ServiceMap {
		serviceNames = append(serviceNames, k)
	}
	logger.Lshortfile.Println("ServiceMap has", len(ServiceMap), "services:", serviceNames)
	for {
		time.Sleep(time.Second * 60)
		now := time.Now().String()[11:19]
		for _, serviceName := range serviceNames {
			num, _ := counter.Get(serviceName)
			logger.Lshortfile.Println(now + " service " + serviceName + " proccessed " + strconv.Itoa(int(num)) + " tasks")
			counter.DeleteAndGetLastValue(serviceName)
		}
	}
}

var ErrBackTo = errors.New("param[\"backTo\"] is not a string")

func NewService(serviceName string, DataRcvBatchSize int64, f fn) {
	serviceName = "svc:" + serviceName
	//check configureation is loaded
	if config.DataRds == nil {
		logger.Lshortfile.Panic("config.DataRedis is nil. you should call config.LoadConfigFromRedis first")
	}
	if config.ParamRds == nil {
		logger.Lshortfile.Panic("config.ParamRedis is nil. you should call config.LoadConfigFromRedis first")
	}

	ServiceMap[serviceName] = f
	counter.DeleteAndGetLastValue(serviceName)

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
		dc := &rCtx.DataCtx{Ctx: context.Background(), Rds: config.DataRds}
		pc := &rCtx.ParamCtx{Ctx: context.Background(), Rds: config.ParamRds}
		if out, err = f(dc, pc, param); err != nil {
			return err
		}
		//Post Back
		if marshaledBytes, err = msgpack.Marshal(out); err != nil {
			return err
		}
		ctx := context.Background()
		pipline := config.ParamRds.Pipeline()
		pipline.RPush(ctx, BackTo, marshaledBytes)
		pipline.Expire(ctx, BackTo, time.Second*6)
		_, err = pipline.Exec(ctx)
		return err
	}
	loop := func() {
		var data []string
		c := context.Background()
		for {
			//fetch datas from redis,using LRange
			pipline := config.ParamRds.Pipeline()
			pipline.LRange(c, serviceName, 0, DataRcvBatchSize-1)
			pipline.LTrim(c, serviceName, DataRcvBatchSize, -1)
			cmd, err := pipline.Exec(c)
			if data = cmd[0].(*redis.StringSliceCmd).Val(); err != nil {
				time.Sleep(time.Millisecond * 100)
				continue
			}

			//try to receive another 1 data, using BLPop
			// 10秒阻塞连接
			if len(data) == 0 {
				rlt := config.ParamRds.BLPop(c, time.Second*10, serviceName)
				if data = rlt.Val(); rlt.Err() == nil && len(data) == 2 {
					data = data[1:]
				}
			}

			for _, s := range data {
				go ProcessOneJob([]byte(s))
				counter.Add(serviceName, 1)
			}
		}
	}
	go loop()
}
