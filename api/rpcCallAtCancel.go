package api

import (
	"context"
	"reflect"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/specification"
)

func CallAtCancel[i any, o any](f func(InParam i) (ret o, err error), timeAt time.Time) (ok bool) {
	var (
		Rds     *redis.Client
		apiInfo *ApiInfo
		Values  []string
		//cmd     *redis.IntCmd
	)
	funcPtr := reflect.ValueOf(f).Pointer()
	if _apiInfo, ok := fun2ApiInfoMap.Load(funcPtr); !ok {
		log.Fatal().Str("service function should be defined By Api or Rpc before used in CallAt", specification.ApiNameByType((*i)(nil))).Send()
	} else {
		apiInfo = _apiInfo.(*ApiInfo)
	}
	if Rds, ok = config.Rds[apiInfo.DbName]; !ok {
		log.Info().Str("DBName not defined in enviroment", apiInfo.DbName).Send()
		return false
	}
	Values = []string{"timeAt", strconv.FormatInt(timeAt.UnixNano(), 10), "data", ""}
	args := &redis.XAddArgs{Stream: apiInfo.ApiName, Values: Values, MaxLen: 4096}
	//use Rds.XAdd rather than Rds.HSet, to prevent Hset before receiing the result of  XAdd
	if cmd := Rds.XAdd(context.Background(), args); cmd.Err() != nil {
		log.Info().AnErr("Do XAdd", cmd.Err()).Send()
		return false
	}
	return true
	// if cmd =; cmd.Err() != nil {
	// 	log.Info().AnErr("Do XAdd", cmd.Err()).Send()
	// 	return out, cmd.Err()
	// }
	// if cmd = Rds.HSet(context.Background(), apiInfo.ApiName+":delay", timeAtStr, "abc123"); cmd.Err() != nil {
	// 	log.Info().Str("CallAtCancel failed. key", apiInfo.ApiName+":delay").Str("field", timeAtStr).Send()
	// 	log.Info().AnErr("CallAtCancel:Do HSet", cmd.Err()).Send()
	// 	return false
	// }
	// cmd1 := Rds.HGet(context.Background(), apiInfo.ApiName+":delay", timeAtStr)
	// log.Info().Str("CallAtCancel done. key", apiInfo.ApiName+":delay").Str("field", timeAtStr).Str("value", cmd1.Val()).Send()
	// return true
}
