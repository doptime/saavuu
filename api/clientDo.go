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
func (ac *Ctx[i, o]) do(paramIn i, dueTime *time.Time) (out o, err error) {
	var (
		b       []byte
		results []string
		cmd     *redis.StringCmd
		Values  []string
	)
	//if F is not nil, call F directly
	if ac.LocalFunc != nil && dueTime == nil {
		return ac.LocalFunc(paramIn)
	}
	//ensure the paramIn is a map or struct
	paramType := reflect.TypeOf(paramIn)
	if paramType.Kind() == reflect.Struct {
	} else if paramType.Kind() == reflect.Map {
	} else if paramType.Kind() == reflect.Ptr && (paramType.Elem().Kind() == reflect.Struct || paramType.Elem().Kind() == reflect.Map) {
	} else {
		logger.Lshortfile.Println("RdsApiBasic param should be a map or struct")
		return out, err
	}

	if b, err = msgpack.Marshal(paramIn); err != nil {
		return out, err
	}
	if dueTime != nil {
		Values = []string{"dueTime", strconv.FormatInt(dueTime.UnixMilli(), 10), "data", string(b)}
	} else {
		Values = []string{"data", string(b)}
	}
	args := &redis.XAddArgs{Stream: ac.ServiceName, Values: Values, MaxLen: 4096}
	if cmd = ac.Rds.XAdd(ac.Ctx, args); cmd.Err() != nil {
		logger.Lshortfile.Println(cmd.Err())
		return out, cmd.Err()
	}
	if dueTime != nil {
		return out, nil
	}

	//BLPop 返回结果 [key1,value1,key2,value2]
	//cmd.Val() is the stream id, the result will be poped from the list with this id
	if results, err = ac.Rds.BLPop(ac.Ctx, time.Second*6, cmd.Val()).Result(); err != nil {
		return out, err
	}

	if len(results) != 2 {
		return out, errors.New("BLPop result length error")
	}
	b = []byte(results[1])

	oType := reflect.TypeOf((*o)(nil)).Elem()
	//if o type is a pointer, use reflect.New to create a new pointer
	if oType.Kind() == reflect.Ptr {
		out = reflect.New(oType.Elem()).Interface().(o)
		return out, msgpack.Unmarshal(b, out)
	}
	oValueWithPointer := reflect.New(oType).Interface().(*o)
	return *oValueWithPointer, msgpack.Unmarshal(b, oValueWithPointer)
}
func (ac *Ctx[i, o]) DoAt(paramIn i, dueTime *time.Time) (err error) {
	_, err = ac.do(paramIn, dueTime)
	return err
}

func (ac *Ctx[i, o]) Do(paramIn i) (out o, err error) {
	return ac.do(paramIn, nil)
}
