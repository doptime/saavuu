package rCtx

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/logger"
)

type ParamCtx struct {
	Ctx context.Context
	Rds *redis.Client
}

// RedisCall: 1.use RPush to push data to redis. 2.use BLPop to pop data from selected channel
// return: error
func (sc *ParamCtx) RdsApiBasic(ServiceKey string, paramIn interface{}) (result []byte, err error) {
	var (
		b       []byte
		BackTo  string
		results []string
	)
	//ensure ServiceKey start with "svc:"
	if ServiceKey[:4] != "svc:" {
		ServiceKey = "svc:" + ServiceKey
	}

	if b, err = msgpack.Marshal(paramIn); err != nil {
		return nil, err
	}
	args := &redis.XAddArgs{Stream: ServiceKey, Values: []string{"data", string(b)}}
	if cmd := sc.Rds.XAdd(sc.Ctx, args); cmd.Err() != nil {
		logger.Lshortfile.Println(cmd.Err())
		return nil, cmd.Err()
	} else {
		BackTo = cmd.Val()
	}

	//BLPop 返回结果 [key1,value1,key2,value2]
	if results, err = sc.Rds.BLPop(sc.Ctx, time.Second*20, BackTo).Result(); err != nil {
		return nil, err
	}
	return []byte(results[1]), nil
}
func (sc *ParamCtx) RdsApi(ServiceKey string, structIn interface{}, out interface{}) (err error) {
	var b []byte
	if b, err = sc.RdsApiBasic(ServiceKey, structIn); err != nil {
		return err
	}
	return msgpack.Unmarshal(b, out)
}
