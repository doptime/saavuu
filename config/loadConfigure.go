package config

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/yangkequn/saavuu/logger"
)

func loadOSEnv(key string, panicString string) (value string) {
	if value = os.Getenv(key); panicString != "" && value == "" {
		logger.Lshortfile.Panicln(panicString)
	}
	return value
}

func init() {
	var err error
	logger.Std.Println("App Start! load config from OS env")

	Cfg.RedisAddressParam = loadOSEnv("REDIS_ADDR_PARAM", "Error: Can not load REDIS_ADDR_PARAM from env")
	Cfg.RedisPasswordParam = loadOSEnv("REDIS_PASSWORD_PARAM", "")
	json.Unmarshal([]byte(loadOSEnv("REDIS_DB_PARAM", "Error: REDIS_DB_PARAM is not set ")), &Cfg.RedisDbParam)

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

	logger.Lshortfile.Println("Load config from env success")
}
