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
		results []string
		cmd     *redis.StringCmd
	)
	//ensure ServiceKey start with "svc:"
	if ServiceKey[:4] != "svc:" {
		ServiceKey = "svc:" + ServiceKey
	}

	if b, err = msgpack.Marshal(paramIn); err != nil {
		return nil, err
	}
	args := &redis.XAddArgs{Stream: ServiceKey, Values: []string{"data", string(b)}, MaxLen: 4096}
	if cmd = sc.Rds.XAdd(sc.Ctx, args); cmd.Err() != nil {
		logger.Lshortfile.Println(cmd.Err())
		return nil, cmd.Err()
	}

	//BLPop 返回结果 [key1,value1,key2,value2]
	//cmd.Val() is the stream id, the result will be poped from the list with this id
	if results, err = sc.Rds.BLPop(sc.Ctx, time.Second*20, cmd.Val()).Result(); err != nil {
		return nil, err
	}
	return []byte(results[1]), nil
}
func (sc *ParamCtx) RdsApi(ServiceKey string, structIn interface{}, out interface{}) (err error) {
	var (
		b []byte
	)
	if b, err = sc.RdsApiBasic(ServiceKey, structIn); err != nil {
		return err
	}
	return msgpack.Unmarshal(b, out)
}
