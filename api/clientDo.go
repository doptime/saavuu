package api

import (
	"errors"
	"reflect"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/specification"
)

// RedisCall: 1.use RPush to push data to redis. 2.use BLPop to pop data from selected channel
// return: error
func (ac *Ctx[i, o]) do(paramIn i, timeAt *time.Time) (out o, err error) {
	var (
		b       []byte
		results []string
		cmd     *redis.StringCmd
		Values  []string
	)
	if b, err = specification.MarshalApiInput(paramIn); err != nil {
		return out, err
	}

	if timeAt != nil {
		Values = []string{"timeAt", strconv.FormatInt(timeAt.UnixMilli(), 10), "data", string(b)}
	} else {
		Values = []string{"data", string(b)}
	}
	args := &redis.XAddArgs{Stream: ac.ServiceName, Values: Values, MaxLen: 4096}
	if cmd = ac.Rds.XAdd(ac.Ctx, args); cmd.Err() != nil {
		log.Info().AnErr("Do XAdd", cmd.Err()).Send()
		return out, cmd.Err()
	}
	if timeAt != nil {
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
func (ac *Ctx[i, o]) DoAt(paramIn i, timeAt *time.Time) (err error) {
	_, err = ac.do(paramIn, timeAt)
	return err
}

func (ac *Ctx[i, o]) Do(paramIn i) (out o, err error) {
	if ac.Func != nil {
		return ac.Func(paramIn)
	}
	return ac.do(paramIn, nil)
}
