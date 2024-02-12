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

// crate ApiFun. the created Api will be called at a specific time in future:
// cautions: there should be a API of same function f in local or remote service, other wise the job will no be executed
//
//	timeAt: the time to execute the job
//	f := func(InParam *InDemo) (ret string, err error) , this is logic function
//	options. there are 2 poosible options:
//		1. api.Name("ServiceName")  //set the ServiceName of the Api.  default value: the name of the InParameter type but with "In" removed
//		2. api.DB("RedisDatabaseName")  //set the DB name of the job. default value: the name of the function
func CallAt[i any, o any](timeAt time.Time, f func(InParam i) (ret o, err error), options ...Options) (retf func(InParam i) (err error)) {
	var (
		db     *redis.Client
		ok     bool
		err    error
		option Options
	)
	if option = optionsMerge(options...); len(option.ServiceName) == 0 {
		option.ServiceName = specification.TypeName((*i)(nil))
	}
	if option.ServiceName, err = specification.ApiName(option.ServiceName); err != nil {
		log.Fatal().AnErr("ApiName", err).Send()
	}

	if db, ok = config.Rds[option.DBName]; !ok {
		log.Info().Str("DBName not defined in enviroment", option.DBName).Send()
		return nil
	}
	if _, ok := ApiServices.Get(option.ServiceName); !ok {
		log.Debug().Str("there should be a API function f in local or remote service, other wise the job will no be executed", option.ServiceName).Send()
	}

	// do At function execute the job at a specific time in futuren , so the returned value can not be returned immediately
	DoAt := func(paramIn i) (err error) {
		var (
			b      []byte
			cmd    *redis.StringCmd
			Values []string
		)
		if b, err = specification.MarshalApiInput(paramIn); err != nil {
			return err
		}

		Values = []string{"timeAt", strconv.FormatInt(timeAt.UnixMilli(), 10), "data", string(b)}
		args := &redis.XAddArgs{Stream: option.ServiceName, Values: Values, MaxLen: 4096}
		if cmd = db.XAdd(context.Background(), args); cmd.Err() != nil {
			log.Info().AnErr("Do XAdd", cmd.Err()).Send()
			return cmd.Err()
		}
		return nil
	}
	//return Api context
	return DoAt
}
