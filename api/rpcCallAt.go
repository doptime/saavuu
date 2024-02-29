package api

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/specification"
)

// create Api context.
// This New function is for the case the API is defined outside of this package.
// If the API is defined in this package, use Api() instead.
// timeAt is ID of the task. if you want's to cancel the task, you should provide the same timeAt
func CallAt[i any, o any](f func(InParam i) (ret o, err error), timeAt time.Time) (retf func(InParam i) (err error)) {
	var (
		db     *redis.Client
		ok     bool
		ctx               = context.Background()
		option *ApiOption = &ApiOption{}
	)
	funcPtr := reflect.ValueOf(f).Pointer()
	if apiInfo, ok := fun2ApiInfoMap.Load(funcPtr); !ok {
		log.Fatal().Str("service function should be defined By Api or Rpc before used in CallAt", specification.ApiNameByType((*i)(nil))).Send()
	} else {
		_apiInfo := apiInfo.(*ApiInfo)
		option.Name_ = _apiInfo.Name
		option.DataSource_ = _apiInfo.DataSource
	}

	if db, ok = config.Rds[option.DataSource_]; !ok {
		log.Info().Str("DataSource not defined in enviroment", option.DataSource_).Send()
		return nil
	}

	retf = func(InParam i) (err error) {
		var (
			b      []byte
			cmd    *redis.StringCmd
			Values []string
		)
		if b, err = specification.MarshalApiInput(InParam); err != nil {
			return err
		}
		fmt.Println("CallAt", option.Name_, timeAt.UnixNano())
		Values = []string{"timeAt", strconv.FormatInt(timeAt.UnixNano(), 10), "data", string(b)}
		args := &redis.XAddArgs{Stream: option.Name_, Values: Values, MaxLen: 4096}
		if cmd = db.XAdd(ctx, args); cmd.Err() != nil {
			log.Info().AnErr("Do XAdd", cmd.Err()).Send()
			return cmd.Err()
		}
		return nil

	}
	return retf
}
