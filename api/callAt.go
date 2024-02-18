package api

import (
	"context"
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
func CallAt[i any, o any](f func(InParam i) (ret o, err error), timeAt time.Time) (retf func(InParam i) (err error)) {
	var (
		db     *redis.Client
		ok     bool
		ctx             = context.Background()
		option *Options = &Options{}
	)
	if apiInfo, ok := fun2ApiInfo.Load(&f); !ok {
		log.Fatal().Str("service function should be defined By Api or Rpc before used in CallAt", specification.ApiNameByType((*i)(nil))).Send()
	} else {
		_apiInfo := apiInfo.(*ApiInfo)
		option.ApiName = _apiInfo.ApiName
		option.DbName = _apiInfo.DbName
	}

	if db, ok = config.Rds[option.DbName]; !ok {
		log.Info().Str("DBName not defined in enviroment", option.DbName).Send()
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
		Values = []string{"timeAt", strconv.FormatInt(timeAt.UnixMilli(), 10), "data", string(b)}
		args := &redis.XAddArgs{Stream: option.ApiName, Values: Values, MaxLen: 4096}
		if cmd = db.XAdd(ctx, args); cmd.Err() != nil {
			log.Info().AnErr("Do XAdd", cmd.Err()).Send()
			return cmd.Err()
		}
		return nil

	}
	return retf
}
