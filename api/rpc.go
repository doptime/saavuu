package api

import (
	"context"
	"errors"
	"reflect"
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
func Rpc[i any, o any](options ...*ApiOption) (retf func(InParam i) (ret o, err error)) {
	var (
		db     *redis.Client
		ok     bool
		ctx               = context.Background()
		option *ApiOption = &ApiOption{}
	)
	if len(options) > 0 {
		option = options[0]
	}

	if len(option.Name_) > 0 {
		option.Name_ = specification.ApiName(option.Name_)
	}
	if len(option.Name_) == 0 {
		option.Name_ = specification.ApiNameByType((*i)(nil))
	}
	if len(option.Name_) == 0 {
		log.Error().Str("service misnamed", option.Name_).Send()
	}

	if db, ok = config.Rds[option.DataSource_]; !ok {
		log.Info().Str("DataSource not defined in enviroment", option.DataSource_).Send()
		return nil
	}

	retf = func(InParam i) (out o, err error) {
		var (
			b       []byte
			results []string
			cmd     *redis.StringCmd
			Values  []string
		)
		if b, err = specification.MarshalApiInput(InParam); err != nil {
			return out, err
		}
		Values = []string{"data", string(b)}
		// if hashCallAt {
		// 	Values = []string{"timeAt", strconv.FormatInt(ops.CallAt.UnixMilli(), 10), "data", string(b)}
		// } else {
		// 	Values = []string{"data", string(b)}
		// }
		args := &redis.XAddArgs{Stream: option.Name_, Values: Values, MaxLen: 4096}
		if cmd = db.XAdd(ctx, args); cmd.Err() != nil {
			log.Info().AnErr("Do XAdd", cmd.Err()).Send()
			return out, cmd.Err()
		}
		// if hashCallAt {
		// 	return out, nil
		// }

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
	rpcInfo := &ApiInfo{
		DataSource: option.DataSource_,
		Name:       option.Name_,
		WithHeader: HeaderFieldsUsed(new(i)),
	}
	funcPtr := reflect.ValueOf(retf).Pointer()
	fun2ApiInfoMap.Store(funcPtr, rpcInfo)
	APIGroupByDataSource.Upsert(option.DataSource_, []string{}, func(exist bool, valueInMap, newValue []string) []string {
		return append(valueInMap, option.Name_)
	})
	return retf
}
