package config

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/go-redis/redis/v8"
	. "github.com/yangkequn/saavuu/redisContext"
)

func LoadFromTomlOrEnviroment(configureFile string, redisConfigurationKey string) (err error) {
	if err = LoadConfigFromToml(configureFile); err != nil {
		if err = LoadConfigFromEnv(); err != nil {
			return err
		}
	}
	UseConfig()
	SaveConfigToRedis(ParamRedis, redisConfigurationKey)
	return nil
}
func LoadConfigFromToml(configFile string) (err error) {
	//read config file to buffer tomlData
	var tomlData []byte
	//if configFile does not exist, use "./config/config.toml" as default
	if _, err = os.Stat(configFile); os.IsNotExist(err) {
		configFile = "./config/config.toml"
	}
	if _, err = os.Stat(configFile); err == nil {
		tomlData, err = os.ReadFile(configFile)
		toml.Decode(string(tomlData), &Cfg)
		fmt.Println("Redis configuration loaded. restart service if you configuration is modified")
	}
	return err
}

func LoadConfigFromEnv() (err error) {

	//try load from env
	if Cfg.RedisAddressParam = os.Getenv("REDIS_ADDR_PARAM"); Cfg.RedisAddressParam == "" {
		return fmt.Errorf("Error: Can not load REDIS_ADDR_PARAM from env")
	}
	if Cfg.RedisPasswordParam = os.Getenv("REDIS_PASSWORD_PARAM"); Cfg.RedisPasswordParam == "" {
		return fmt.Errorf("Error: Can not load REDIS_PASSWORD_PARAM from env")
	}
	if Cfg.RedisDbParam, err = strconv.Atoi(os.Getenv("REDIS_DB_PARAM")); err != nil {
		return fmt.Errorf("Error: Can not load REDIS_DB_PARAM from env")
	}
	if Cfg.RedisAddressData = os.Getenv("REDIS_ADDR_DATA"); Cfg.RedisAddressData == "" {
		return fmt.Errorf("Error: Can not load REDIS_ADDR_DATA from env")
	}
	if Cfg.RedisPasswordData = os.Getenv("REDIS_PASSWORD_DATA"); Cfg.RedisPasswordData == "" {
		return fmt.Errorf("Error: Can not load REDIS_PASSWORD_DATA from env")
	}
	if Cfg.RedisDbData, err = strconv.Atoi(os.Getenv("REDIS_DB_DATA")); err != nil {
		return fmt.Errorf("Error: Can not load REDIS_DB_DATA from env")
	}
	return nil
}

func LoadConfigFromRedis(ParamServer *redis.Client, keyName string) (err error) {
	rc := RedisContext{Ctx: context.Background(), RdsClient: ParamServer}
	if err = rc.Get(keyName, &Cfg); err != nil {
		return err
	}
	UseConfig()
	return nil
}
func SaveConfigToRedis(ParamServer *redis.Client, keyName string) (err error) {
	rc := RedisContext{Ctx: context.Background(), RdsClient: ParamServer}
	if err = rc.Set(keyName, &Cfg, -1); err != nil {
		return err
	}
	return nil
}
