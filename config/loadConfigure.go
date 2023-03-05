package config

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/yangkequn/saavuu/logger"
	"github.com/yangkequn/saavuu/rds"
)

func loadOSEnv(key string, panicString string) (value string) {
	if value = os.Getenv(key); panicString != "" && value == "" {
		logger.Lshortfile.Panicln(panicString)
	}
	return value
}
func loadParamRds() (err error) {
	Cfg.RedisAddressParam = loadOSEnv("REDIS_ADDR_PARAM", "Error: Can not load REDIS_ADDR_PARAM from env")
	Cfg.RedisPasswordParam = loadOSEnv("REDIS_PASSWORD_PARAM", "")
	return json.Unmarshal([]byte(loadOSEnv("REDIS_DB_PARAM", "Error: REDIS_DB_PARAM is not set ")), &Cfg.RedisDbParam)
}

func initialAppFramework() (err error) {
	//try load from env
	Cfg.RedisAddressData = loadOSEnv("REDIS_ADDR_DATA", "Error: Can not load REDIS_ADDR_DATA from env")
	Cfg.RedisPasswordData = loadOSEnv("REDIS_PASSWORD_DATA", "")
	if Cfg.RedisDbData, err = strconv.Atoi(loadOSEnv("REDIS_DB_DATA", "Error: REDIS_DB_DATA is not set ")); err != nil {
		logger.Lshortfile.Panicln("Error: REDIS_DB_DATA is not a number")
	}
	Cfg.JwtSecret = loadOSEnv("JWT_SECRET", "Error: JWT_SECRET Can not load from env")
	if Cfg.JwtFieldsKept = loadOSEnv("JWT_FIELDS_KEPT", ""); Cfg.JwtFieldsKept != "" {
		Cfg.JwtFieldsKept = strings.ToLower(Cfg.JwtFieldsKept)
	}
	if Cfg.MaxBufferSize, err = strconv.ParseInt(loadOSEnv("MAX_BUFFER_SIZE", "Error: MAX_BUFFER_SIZE is not set "), 10, 64); err != nil {
		logger.Lshortfile.Panicln("Error: MAX_BUFFER_SIZE is not a number")
	}
	if dev := loadOSEnv("DEVELOP_MODE", ""); len(dev) > 0 {
		if Cfg.DevelopMode, err = strconv.ParseBool(dev); err != nil {
			logger.Lshortfile.Println("Error: bad string of env: DEVELOP_MODE")
		}
		logger.Std.Println("DEVELOP_MODE is set to ", Cfg.DevelopMode)
	}
	if ServiceBatchSize := loadOSEnv("SERVICE_BATCH_SIZE", ""); len(ServiceBatchSize) > 0 {
		if Cfg.ServiceBatchSize, err = strconv.ParseInt(ServiceBatchSize, 10, 64); err != nil {
			logger.Lshortfile.Println("Error: bad string of env: SERVICE_BATCH_SIZE")
		}
		logger.Std.Println("SERVICE_BATCH_SIZE is set to ", Cfg.ServiceBatchSize)
	}

	useParamRedis()
	useDataRedis()

	// 保存到 ParamServer
	if err = rds.Set(context.Background(), ParamRds, redisConfigKey, &Cfg, -1); err != nil {
		return err
	}

	logger.Lshortfile.Println("Load config from env success")
	return nil
}

const redisConfigKey = "saavuu_config"

func initialAppService() (err error) {
	ParamRds = redis.NewClient(&redis.Options{
		Addr:     Cfg.RedisAddressParam,
		Password: Cfg.RedisPasswordParam, // no password set
		DB:       Cfg.RedisDbParam,       // use default DB
	})
	if err = rds.Get(context.Background(), ParamRds, redisConfigKey, &Cfg); err != nil {
		logger.Lshortfile.Panicln(err)
	}
	useParamRedis()
	useDataRedis()
	if DataRds == nil {
		logger.Lshortfile.Panicln("config.DataRedis is nil. ")
	}
	logger.Lshortfile.Println("ApiInitial success")
	return nil
}

func init() {
	logger.Std.Println("App Start! load config from OS env")
	loadParamRds()
	AppMode = loadOSEnv("APP_MODE", "APP_MODE , should be either "+AppModeFRAMEWROK+" as www server or "+AppModeSERVICE+" to run services")
	logger.Lshortfile.Println("Using AppMode " + AppMode)

	if AppMode == AppModeFRAMEWROK {
		initialAppFramework()
	} else if AppMode == AppModeSERVICE {
		initialAppService()
	} else {
		logger.Lshortfile.Panicln("APP_MODE , should be either FRAMWORK as www server or SERVICE to run services")
	}

}
