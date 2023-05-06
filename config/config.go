package config

import (
	"github.com/redis/go-redis/v9"
)

type Configuration struct {
	RPCFirst bool
	//on redis server
	RedisAddress  string
	RedisPassword string
	RedisDb       int

	//AutoPermission should never be true in production
	AutoPermission bool
	JWTSecret      string
	JwtFieldsKept  string
	MaxBufferSize  int64
	CORS           string

	ServerPort int64
	ServerPath string

	//ServiceBatchSize is the number of tasks that a service can read from redis at the same time
	ServiceBatchSize int64
}

var Cfg Configuration = Configuration{
	RPCFirst:      false,
	JWTSecret:     "",
	MaxBufferSize: 32 << 20,
	CORS:          "*",
	ServerPort:    8000,
	ServerPath:    "/",
}

var Rds *redis.Client
