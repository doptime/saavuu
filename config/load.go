package config

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/yangkequn/saavuu/logger"
)

func loadOSEnv(key string, des interface{}, defaultValue interface{}) {
	var value string
	if value = os.Getenv(key); defaultValue == nil && value == "" {
		logger.Lshortfile.Panicln("Panic:", key, " Can not load from env")
	} else if defaultValue != nil && value == "" {
		//des is a pointer, so we can set it to defaultValue using reflection
		reflect.ValueOf(des).Elem().Set(reflect.ValueOf(defaultValue))
		return
	}
	if err := json.Unmarshal([]byte(value), des); err == nil {
	} else if _, ok := des.(*string); ok {
		*des.(*string) = value
	} else {
		logger.Lshortfile.Panicln("Error: bad type of env: ", key)
	}
	logger.Std.Println("env", key, "is set to ", reflect.ValueOf(des).Elem().Interface())
}

func init() {
	logger.Std.Println("App Start! load config from OS env")

	loadOSEnv("RedisAddress", &Cfg.RedisAddress, nil)
	loadOSEnv("RedisPassword", &Cfg.RedisPassword, "")
	loadOSEnv("RedisDb", &Cfg.RedisDb, nil)
	loadOSEnv("JWTSecret", &Cfg.JWTSecret, "")
	//try load from env
	if loadOSEnv("JwtFieldsKept", &Cfg.JwtFieldsKept, ""); Cfg.JwtFieldsKept != "" {
		Cfg.JwtFieldsKept = strings.ToLower(Cfg.JwtFieldsKept)
	}
	loadOSEnv("CORS", &Cfg.CORS, "*")
	loadOSEnv("MaxBufferSize", &Cfg.MaxBufferSize, 5*1024*1024)
	loadOSEnv("ServiceBatchSize", &Cfg.ServiceBatchSize, 256)
	loadOSEnv("AutoPermission", &Cfg.AutoPermission, false)
	loadOSEnv("ServerPort", &Cfg.ServerPort, 8000)
	loadOSEnv("ServerPath", &Cfg.ServerPath, "/")

	//apply configuration
	redisOption := &redis.Options{
		Addr:     Cfg.RedisAddress,
		Password: Cfg.RedisPassword, // no password set
		DB:       Cfg.RedisDb,       // use default DB
	}
	Rds = redis.NewClient(redisOption)

	logger.Lshortfile.Println("Load config from env success")
}
