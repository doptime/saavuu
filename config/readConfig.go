package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

func LoadConfigFromEnv() {
	REDIS_ADDRESS_PARAM := os.Getenv("REDIS_ADDRESS_PARAM")
	if len(REDIS_ADDRESS_PARAM) == 0 {
		panic("REDIS_ADDRESS_PARAM is not set")
	}

	REDIS_PASSWORD_PARAM := os.Getenv("REDIS_PASSWORD_PARAM")
	REDIS_DB_PARAM := os.Getenv("REDIS_DB_PARAM")
	if len(REDIS_DB_PARAM) == 0 {
		fmt.Println("REDIS_DB_PARAM is not set, use default 0")
		REDIS_DB_PARAM = "0"
	}
	REDIS_DB_PARAM_INT, err := strconv.Atoi(REDIS_DB_PARAM)
	if err != nil {
		panic("REDIS_DB_DATA is not a number")
	}
	ParamRedisOption := &redis.Options{
		Addr:     REDIS_ADDRESS_PARAM,
		Password: REDIS_PASSWORD_PARAM, // no password set
		DB:       REDIS_DB_PARAM_INT,   // use default DB
	}
	Cfg.ParamRedis = redis.NewClient(ParamRedisOption)

	//read redis data config
	REDIS_ADDRESS_DATA := os.Getenv("REDIS_ADDRESS_DATA")
	if len(REDIS_ADDRESS_DATA) == 0 {
		panic("REDIS_ADDRESS_DATA is not set")
	}
	REDIS_PASSWORD_DATA := os.Getenv("REDIS_PASSWORD_DATA")
	REDIS_DB_DATA := os.Getenv("REDIS_DB_DATA")
	if len(REDIS_DB_DATA) == 0 {
		fmt.Println("REDIS_DB_DATA is not set, use default 0")
		REDIS_DB_DATA = "0"
	}
	REDIS_DB_DATA_INT, err := strconv.Atoi(REDIS_DB_DATA)
	if err != nil {
		panic("REDIS_DB_DATA is not a number")
	}
	DataRedisOption := &redis.Options{
		Addr:     REDIS_ADDRESS_DATA,
		Password: REDIS_PASSWORD_DATA, // no password set
		DB:       REDIS_DB_DATA_INT,   // use default DB
	}
	Cfg.DataRedis = redis.NewClient(DataRedisOption)

	if v := os.Getenv("JWT_SECRET"); len(v) > 0 {
		Cfg.JwtSecret = v
	} else {
		panic("JWT_SECRET is not set")
	}
	if v := os.Getenv("MAX_BUFFER_SIZE"); len(v) > 0 {
		Cfg.MaxBufferSize, _ = strconv.ParseInt(v, 10, 64)
	}
	if v := os.Getenv("CORS"); len(v) > 0 {
		Cfg.CORS = v
	}
}
