package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v9"
	"github.com/yangkequn/saavuu/rCtx"
)

func loadOSEnv(key string, panicString string) (value string) {
	if value = os.Getenv(key); panicString != "" && value == "" {
		fmt.Println(panicString)
		os.Exit(1)
	}
	return value
}

func LoadConfigFromEnv() (err error) {

	//try load from env
	Cfg.RedisAddressParam = loadOSEnv("REDIS_ADDR_PARAM", "Error: Can not load REDIS_ADDR_PARAM from env")
	Cfg.RedisPasswordParam = loadOSEnv("REDIS_PASSWORD_PARAM", "")
	if Cfg.RedisDbParam, err = strconv.Atoi(loadOSEnv("REDIS_DB_PARAM", "Error: REDIS_DB_PARAM is not set ")); err != nil {
		fmt.Println("Error: REDIS_DB_PARAM is not a number")
		os.Exit(1)
	}

	Cfg.RedisAddressData = loadOSEnv("REDIS_ADDR_DATA", "Error: Can not load REDIS_ADDR_DATA from env")
	Cfg.RedisPasswordData = loadOSEnv("REDIS_PASSWORD_DATA", "")
	if Cfg.RedisDbData, err = strconv.Atoi(loadOSEnv("REDIS_DB_DATA", "Error: REDIS_DB_DATA is not set ")); err != nil {
		fmt.Println("Error: REDIS_DB_DATA is not a number")
		os.Exit(1)
	}
	Cfg.JwtSecret = loadOSEnv("JWT_SECRET", "Error: JWT_SECRET Can not load from env")
	if Cfg.JwtIgnoreFields = loadOSEnv("JWT_IGNORE_FIELDS", ""); Cfg.JwtIgnoreFields != "" {
		Cfg.JwtIgnoreFields = strings.ToLower(Cfg.JwtIgnoreFields)
	}
	if Cfg.MaxBufferSize, err = strconv.ParseInt(loadOSEnv("MAX_BUFFER_SIZE", "Error: MAX_BUFFER_SIZE is not set "), 10, 64); err != nil {
		fmt.Println("Error: MAX_BUFFER_SIZE is not a number")
		os.Exit(1)
	}

	UseConfig()
	SaveConfigToRedis(ParamRds, loadOSEnv("SAAVUU_CONFIG_KEY", "Error: Can not load SAAVUU_CONFIG_KEY from env"))
	return nil
}

func LoadConfigFromRedis(ParamServer *redis.Client, keyName string) (err error) {
	// 保存到 ParamServer
	rc := rCtx.DataCtx{Ctx: context.Background(), Rds: ParamServer}
	if err = rc.Get(keyName, &Cfg); err != nil {
		return err
	}
	UseConfig()
	return nil
}
func SaveConfigToRedis(ParamServer *redis.Client, keyName string) (err error) {
	// 保存到 ParamServer
	rc := rCtx.DataCtx{Ctx: context.Background(), Rds: ParamServer}
	if err = rc.Set(keyName, &Cfg, -1); err != nil {
		return err
	}
	return nil
}
