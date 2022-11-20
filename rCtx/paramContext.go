package rCtx

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

type ParamCtx struct {
	Ctx context.Context
	Rds *redis.Client
}

// RedisCall: 1.use RPush to push data to redis. 2.use BLPop to pop data from selected channel
// return: error
func (sc *ParamCtx) RdsApiBasic(ServiceKey string, paramIn map[string]interface{}) (result []byte, err error) {
	var (
		b       []byte
		BackTo  string = fmt.Sprintf("%x", rand.Int63())
		results []string
	)
	//ensure ServiceKey start with "svc:"
	if ServiceKey[:4] != "svc:" {
		ServiceKey = "svc:" + ServiceKey
	}
	paramIn["BackTo"] = BackTo

	if b, err = msgpack.Marshal(paramIn); err != nil {
		return nil, err
	}
	ppl := sc.Rds.Pipeline()
	ppl.RPush(sc.Ctx, ServiceKey, b)
	//长期不执行的任务，抛弃
	ppl.Expire(sc.Ctx, ServiceKey, time.Second*60)
	if _, err := ppl.Exec(sc.Ctx); err != nil {
		return nil, err
	}
	//BLPop 返回结果 [key1,value1,key2,value2]
	if results, err = sc.Rds.BLPop(sc.Ctx, time.Second*20, BackTo).Result(); err != nil {
		return nil, err
	}
	return []byte(results[1]), nil
}
func (sc *ParamCtx) RdsApi(ServiceKey string, structIn interface{}, out interface{}) (err error) {
	var (
		paramIn = map[string]interface{}{}
		ok      bool
	)
	//convert struct to map, to allow add field "BackTo"
	if paramIn, ok = structIn.(map[string]interface{}); !ok {
		b, err := msgpack.Marshal(structIn)
		if err != nil {
			return err
		}
		if err = msgpack.Unmarshal(b, &paramIn); err != nil {
			return err
		}
	}

	result, err := sc.RdsApiBasic(ServiceKey, paramIn)
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(result, out)
}
