package api

import (
	"reflect"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/logger"
)

// RedisCall: 1.use RPush to push data to redis. 2.use BLPop to pop data from selected channel
// return: error
func (sc *Ctx) Api(ServiceKey string, paramIn interface{}, out interface{}, dueTime int64) (err error) {
	var (
		b       []byte
		results []string
		cmd     *redis.StringCmd
		Values  []string
	)
	//ensure ServiceKey start with "api:"
	if ServiceKey[:4] != "api:" {
		ServiceKey = "api:" + ServiceKey
	}

	//ensure the paramIn is a map or struct

	paramType := reflect.TypeOf(paramIn)
	if paramType.Kind() == reflect.Struct {
	} else if paramType.Kind() == reflect.Map {
	} else if paramType.Kind() == reflect.Ptr && (paramType.Elem().Kind() == reflect.Struct || paramType.Elem().Kind() == reflect.Map) {
	} else {
		logger.Lshortfile.Println("RdsApiBasic param should be a map or struct")
		return err
	}

	if b, err = msgpack.Marshal(paramIn); err != nil {
		return err
	}
	if dueTime != 0 {
		Values = []string{"dueTime", strconv.FormatInt(dueTime, 10), "data", string(b)}
	} else {
		Values = []string{"data", string(b)}
	}
	args := &redis.XAddArgs{Stream: ServiceKey, Values: Values, MaxLen: 4096}
	if cmd = sc.Rds.XAdd(sc.Ctx, args); cmd.Err() != nil {
		logger.Lshortfile.Println(cmd.Err())
		return cmd.Err()
	}
	if dueTime != 0 {
		return nil
	}

	//BLPop 返回结果 [key1,value1,key2,value2]
	//cmd.Val() is the stream id, the result will be poped from the list with this id
	if results, err = sc.Rds.BLPop(sc.Ctx, time.Second*20, cmd.Val()).Result(); err != nil {
		return err
	} else if out != nil && len(results) == 2 {
		b = []byte(results[1])
		return msgpack.Unmarshal(b, out)
	}
	return nil
}
