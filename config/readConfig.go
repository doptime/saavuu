package config

import (
	"context"
	"fmt"
	"os"
	"saavuu/redisService"
	"strconv"

	"github.com/go-redis/redis/v8"
)

func LoadConfigFromEnv() (err error) {
	ctx := context.Background()
	os.Setenv("REDIS_ADDRESS_PARAM", "docker.vm:6379")
	os.Setenv("REDIS_PASSWORD_PARAM", "")
	os.Setenv("REDIS_DB_PARAM", "0")
	os.Setenv("REDIS_CONFIG_KEY", "config_saavuu_services")
	REDIS_ADDRESS_PARAM := os.Getenv("REDIS_ADDRESS_PARAM")
	if len(REDIS_ADDRESS_PARAM) == 0 {
		panic("env REDIS_ADDRESS_PARAM is not set")
	}

	REDIS_PASSWORD_PARAM := os.Getenv("REDIS_PASSWORD_PARAM")
	REDIS_DB_PARAM := os.Getenv("REDIS_DB_PARAM")
	if len(REDIS_DB_PARAM) == 0 {
		fmt.Println("env REDIS_DB_PARAM is not set, use default 0")
		REDIS_DB_PARAM = "0"
	}
	REDIS_DB_PARAM_INT, err := strconv.Atoi(REDIS_DB_PARAM)
	if err != nil {
		panic("env REDIS_DB_DATA is not a number")
	}
	ParamRedisOption := &redis.Options{
		Addr:     REDIS_ADDRESS_PARAM,
		Password: REDIS_PASSWORD_PARAM, // no password set
		DB:       REDIS_DB_PARAM_INT,   // use default DB
	}
	Cfg.ParamRedis = redis.NewClient(ParamRedisOption)
	ConfigKey := os.Getenv("REDIS_CONFIG_KEY")
	if len(ConfigKey) == 0 {
		panic("env REDIS_CONFIG_KEY is not set")
	}

	rds := Cfg.ParamRedis
	REDIS_ADDRESS_DATA := REDIS_ADDRESS_PARAM
	if err = redisService.HGet(ctx, rds, ConfigKey, "REDIS_ADDRESS_DATA", &REDIS_ADDRESS_DATA); err != nil {
		redisService.HSet(ctx, rds, ConfigKey, "REDIS_ADDRESS_DATA", REDIS_ADDRESS_PARAM)
		fmt.Println("REDIS_ADDRESS_DATA is not set, use default REDIS_ADDRESS_PARAM")
	}
	redisService.HSet(ctx, rds, ConfigKey, "REDIS_ADDRESS_DATA_annotation", "http service get method will fetch data from this redis directly")

	REDIS_PASSWORD_DATA := REDIS_PASSWORD_PARAM
	if err = redisService.HGet(ctx, rds, ConfigKey, "REDIS_PASSWORD_DATA", &REDIS_PASSWORD_DATA); err != nil {
		redisService.HSet(ctx, rds, ConfigKey, "REDIS_PASSWORD_DATA", REDIS_PASSWORD_PARAM)
		fmt.Println("REDIS_PASSWORD_DATA is not set, use default REDIS_PASSWORD_PARAM")
	}
	REDIS_DB_DATA := REDIS_DB_PARAM
	if err = redisService.HGet(ctx, rds, ConfigKey, "REDIS_DB_DATA", &REDIS_DB_DATA); err != nil {
		redisService.HSet(ctx, rds, ConfigKey, "REDIS_DB_DATA", REDIS_DB_PARAM_INT)
		fmt.Println("REDIS_DB_DATA is not set, use default REDIS_DB_PARAM")
	}
	REDIS_DB_DATA_INT, err := strconv.Atoi(REDIS_DB_DATA)
	if err != nil {
		panic("REDIS_DB_DATA is not a number, check redis key: " + ConfigKey + "with field REDIS_DB_DATA  in redis param server. ")
	}
	DataRedisOption := &redis.Options{
		Addr:     REDIS_ADDRESS_DATA,
		Password: REDIS_PASSWORD_DATA, // no password set
		DB:       REDIS_DB_DATA_INT,   // use default DB
	}
	Cfg.DataRedis = redis.NewClient(DataRedisOption)
	//read jwt secret
	if err = redisService.HGet(ctx, rds, ConfigKey, "JWT_SECRET", &Cfg.JwtSecret); err != nil {
		redisService.HSet(ctx, rds, ConfigKey, "JWT_SECRET", "")
		fmt.Println("JWT_SECRET is not set, use default empty string. Restart service after you modify this value")
	}
	//read CORS
	if err = redisService.HGet(ctx, rds, ConfigKey, "CORS", &Cfg.CORS); err != nil {
		redisService.HSet(ctx, rds, ConfigKey, "CORS", "*")
		fmt.Println("CORS is not set, use default *")
	}
	return err
}
