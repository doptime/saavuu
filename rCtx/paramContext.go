package rCtx

import (
	"context"
	"strconv"
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
func (sc *ParamCtx) RdsApiBasic(ServiceKey string, paramIn interface{}, dueTime int64) (result []byte, err error) {
	var (
		b       []byte
		results []string
		cmd     *redis.StringCmd
		Values  []string
		ok      bool
	)
	//ensure ServiceKey start with "svc:"
	if ServiceKey[:4] != "svc:" {
		ServiceKey = "svc:" + ServiceKey
	}

	//ensure the paramIn is a map or struct
	//later, paramIn will be deocded to a map in newService
	if paramIn, ok = paramIn.(map[string]interface{}); !ok {
		if b, err = msgpack.Marshal(paramIn); err != nil {
			return nil, err
		}
		if err = msgpack.Unmarshal(b, &paramIn); err != nil {
			logger.Lshortfile.Println("RdsApiBasic param should be a map or struct")
			return nil, err
		}
	}

	if b, err = msgpack.Marshal(paramIn); err != nil {
		return nil, err
	}
	if dueTime != 0 {
		Values = []string{"dueTime", strconv.FormatInt(dueTime, 10), "data", string(b)}
	} else {
		Values = []string{"data", string(b)}
	}
	args := &redis.XAddArgs{Stream: ServiceKey, Values: Values, MaxLen: 4096}
	if cmd = sc.Rds.XAdd(sc.Ctx, args); cmd.Err() != nil {
		logger.Lshortfile.Println(cmd.Err())
		return nil, cmd.Err()
	}
	if dueTime != 0 {
		return nil, nil
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
	if b, err = sc.RdsApiBasic(ServiceKey, structIn, 0); err != nil {
		return err
	}
	return msgpack.Unmarshal(b, out)
}
