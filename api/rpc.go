package api

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/specification"
)

// create Api context.
// This New function is for the case the API is defined outside of this package.
// If the API is defined in this package, use Api() instead.
func Rpc[i any, o any](options ...Option) (retf func(InParam i, callAt ...time.Time) (ret o, err error)) {
	var (
		db     *redis.Client
		ok     bool
		ctx             = context.Background()
		option *Options = optionsMerge(options...)
	)

	if len(option.ServiceName) > 0 {
		option.ServiceName = specification.ApiName(option.ServiceName)
	}
	if len(option.ServiceName) == 0 {
		option.ServiceName = specification.ApiNameByType((*i)(nil))
	}
	if len(option.ServiceName) == 0 {
		log.Error().Str("service misnamed", option.ServiceName).Send()
	}

	if db, ok = config.Rds[option.DbName]; !ok {
		log.Info().Str("DBName not defined in enviroment", option.DbName).Send()
		return nil
	}

	retf = func(InParam i, callAt ...time.Time) (out o, err error) {
		var (
			b       []byte
			results []string
			cmd     *redis.StringCmd
			Values  []string
		)
		if b, err = specification.MarshalApiInput(InParam); err != nil {
			return out, err
		}
		if len(callAt) > 0 {
			timeAt := callAt[0]
			Values = []string{"timeAt", strconv.FormatInt(timeAt.UnixMilli(), 10), "data", string(b)}
		} else {
			Values = []string{"data", string(b)}
		}
		args := &redis.XAddArgs{Stream: option.ServiceName, Values: Values, MaxLen: 4096}
		if cmd = db.XAdd(ctx, args); cmd.Err() != nil {
			log.Info().AnErr("Do XAdd", cmd.Err()).Send()
			return out, cmd.Err()
		}
		if len(callAt) == 0 {
			return out, nil
		}

		//BLPop 返回结果 [key1,value1,key2,value2]
		//cmd.Val() is the stream id, the result will be poped from the list with this id
		if results, err = db.BLPop(ctx, time.Second*6, cmd.Val()).Result(); err != nil {
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
	return retf
}
