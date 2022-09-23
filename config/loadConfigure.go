package config

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/yangkequn/saavuu/redisContext"
)

func LoadConfigFromEnv() {
	var (
		ConfigKey string
		err       error
	)

	//try load from env
	if Cfg.RedisAddressParam = os.Getenv("REDIS_ADDR_PARAM"); Cfg.RedisAddressParam == "" {
		panic("Error: Can not load REDIS_ADDR_PARAM from env")
	}
	if Cfg.RedisPasswordParam = os.Getenv("REDIS_PASSWORD_PARAM"); Cfg.RedisPasswordParam == "" {
	}
	if Cfg.RedisDbParam, err = strconv.Atoi(os.Getenv("REDIS_DB_PARAM")); err != nil {
		panic("Error: Can not load REDIS_DB_PARAM from env")
	}
	if Cfg.RedisAddressData = os.Getenv("REDIS_ADDR_DATA"); Cfg.RedisAddressData == "" {
		panic("Error: Can not load REDIS_ADDR_DATA from env")
	}
	if Cfg.RedisPasswordData = os.Getenv("REDIS_PASSWORD_DATA"); Cfg.RedisPasswordData == "" {
	}
	if Cfg.RedisDbData, err = strconv.Atoi(os.Getenv("REDIS_DB_DATA")); err != nil {
		panic("Error: Can not load REDIS_DB_DATA from env")
	}
	if Cfg.JwtSecret = os.Getenv("JWT_SECRET"); Cfg.JwtSecret == "" {
		panic("Error: Can not load JWT_SECRET from env")
	}
	if Cfg.JwtIgnoreFields = os.Getenv("JWT_IGNORE_FIELDS"); Cfg.JwtIgnoreFields != "" {
		//to lower
		Cfg.JwtIgnoreFields = strings.ToLower(Cfg.JwtIgnoreFields)
	}
	if Cfg.MaxBufferSize, err = strconv.ParseInt(os.Getenv("MAX_BUFFER_SIZE"), 10, 64); err != nil {
		panic("Error: Can not load MAX_BUFFER_SIZE from env")
	}
	if ConfigKey = os.Getenv("SAAVUU_CONFIG_KEY"); ConfigKey == "" {
		panic("Error: Can not load SAAVUU_CONFIG_KEY from env")
	}
	UseConfig()
	SaveConfigToRedis(ParamRds, ConfigKey)
}

func LoadConfigFromRedis(ParamServer *redis.Client, keyName string) (err error) {
	// 保存到 ParamServer
	rc := redisContext.DataCtx{Ctx: context.Background(), Rds: ParamServer}
	if err = rc.Get(keyName, &Cfg); err != nil {
		return err
	}
	UseConfig()
	return nil
}
func SaveConfigToRedis(ParamServer *redis.Client, keyName string) (err error) {
	// 保存到 ParamServer
	rc := redisContext.DataCtx{Ctx: context.Background(), Rds: ParamServer}
	if err = rc.Set(keyName, &Cfg, -1); err != nil {
		return err
	}
	return nil
}
