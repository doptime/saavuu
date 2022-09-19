package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/go-redis/redis/v8"
)

func LoadConfigFromEnv(configFile string) (err error) {
	//read config file to buffer tomlData
	var tomlData []byte
	//if configFile does not exist, use "./config/config.toml" as default
	if _, err = os.Stat(configFile); os.IsNotExist(err) {
		configFile = "./config/config.toml"
	}
	if _, err = os.Stat(configFile); !os.IsNotExist(err) {
		tomlData, err = os.ReadFile(configFile)
	}
	if err != nil {
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
	}

	toml.Decode(string(tomlData), &Cfg)
	ParamRedisOption := &redis.Options{
		Addr:     Cfg.RedisAddressParam,
		Password: Cfg.RedisPasswordParam, // no password set
		DB:       Cfg.RedisDbParam,       // use default DB
	}
	ParamRedis = redis.NewClient(ParamRedisOption)

	DataRedisOption := &redis.Options{
		Addr:     Cfg.RedisAddressData,
		Password: Cfg.RedisPasswordData, // no password set
		DB:       Cfg.RedisDbData,       // use default DB
	}
	DataRedis = redis.NewClient(DataRedisOption)
	fmt.Println("Redis configuration loaded. restart service if you configuration is modified")
	return nil
}
