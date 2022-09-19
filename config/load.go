package config

import (
	"fmt"
	"os"

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
	tomlData, err = os.ReadFile(configFile)
	if err != nil {
		fmt.Println("read config file error:", err)
		return err
	}

	toml.Decode(string(tomlData), &Cfg)
	ParamRedisOption := &redis.Options{
		Addr:     Cfg.RedisAddressParam,
		Password: Cfg.RedisPasswordParam, // no password set
		DB:       Cfg.RedisDbParam,       // use default DB
	}
	ParamRedis = redis.NewClient(ParamRedisOption)

	DataRedisOption := &redis.Options{
		Addr:     Cfg.RedisAdressData,
		Password: Cfg.RedisPasswordData, // no password set
		DB:       Cfg.RedisDbData,       // use default DB
	}
	DataRedis = redis.NewClient(DataRedisOption)
	fmt.Println("Redis configuration loaded. restart service if you configuration is modified")
	return err
}
