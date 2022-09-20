package config

import "github.com/go-redis/redis/v8"

func UseConfig() {
	ParamRedisOption := &redis.Options{
		Addr:     Cfg.RedisAddressParam,
		Password: Cfg.RedisPasswordParam, // no password set
		DB:       Cfg.RedisDbParam,       // use default DB
	}
	ParamRds = redis.NewClient(ParamRedisOption)

	DataRedisOption := &redis.Options{
		Addr:     Cfg.RedisAddressData,
		Password: Cfg.RedisPasswordData, // no password set
		DB:       Cfg.RedisDbData,       // use default DB
	}
	DataRds = redis.NewClient(DataRedisOption)
}
