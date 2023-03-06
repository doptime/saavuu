package api

import (
	"errors"
	"reflect"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/logger"
)

// RedisCall: 1.use RPush to push data to redis. 2.use BLPop to pop data from selected channel
// return: error
func (ac *Ctx) do(paramIn interface{}, out interface{}, dueTime *time.Time) (err error) {
	var (
		b       []byte
		results []string
		cmd     *redis.StringCmd
		Values  []string
	)
	//if service name is for redis, return error
	if ac.ServiceName == "api:redis" {
		return errors.New("api:redis not allowed to call")
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
	if dueTime != nil {
		Values = []string{"dueTime", strconv.FormatInt(dueTime.UnixMilli(), 10), "data", string(b)}
	} else {
		Values = []string{"data", string(b)}
	}
	args := &redis.XAddArgs{Stream: ac.ServiceName, Values: Values, MaxLen: 4096}
	if cmd = ac.Rds.XAdd(ac.Ctx, args); cmd.Err() != nil {
		logger.Lshortfile.Println(cmd.Err())
		return cmd.Err()
	}
	if dueTime != nil {
		return nil
	}

	//BLPop 返回结果 [key1,value1,key2,value2]
	//cmd.Val() is the stream id, the result will be poped from the list with this id
	if results, err = ac.Rds.BLPop(ac.Ctx, time.Second*20, cmd.Val()).Result(); err != nil {
		return err
	} else if out != nil && len(results) == 2 {
		b = []byte(results[1])
		return msgpack.Unmarshal(b, out)
	}
	return nil
}
func (ac *Ctx) DoAt(paramIn interface{}, dueTime *time.Time) (err error) {
	return ac.do(paramIn, nil, dueTime)
}

func (ac *Ctx) Do(paramIn interface{}, out interface{}) (err error) {
	return ac.do(paramIn, out, nil)
}
