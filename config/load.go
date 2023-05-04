package config

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
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

	Cfg.RedisAddress = loadOSEnv("REDIS_ADDR", "Error: Can not load REDIS_ADDR from env")
	Cfg.RedisPassword = loadOSEnv("REDIS_PASSWORD", "")
	json.Unmarshal([]byte(loadOSEnv("REDIS_DB", "Error: REDIS_DB is not set ")), &Cfg.RedisDb)

	//try load from env
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

	//apply configuration
	redisOption := &redis.Options{
		Addr:     Cfg.RedisAddress,
		Password: Cfg.RedisPassword, // no password set
		DB:       Cfg.RedisDb,       // use default DB
	}
	ParamRds = redis.NewClient(redisOption)

	logger.Lshortfile.Println("Load config from env success")
}
