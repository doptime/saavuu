package config

import (
	"github.com/go-redis/redis/v8"
)

type Configuration struct {
	Rds           *redis.Client
	JwtSecret     string
	MaxBufferSize int64
	CORS          string
}

var Cfg Configuration = Configuration{
	Rds:           nil,
	JwtSecret:     "",
	MaxBufferSize: 32 << 20,
	CORS:          "*",
}
