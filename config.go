package saavuu

import (
	"github.com/go-redis/redis/v8"
)

type Configuration struct {
	Port          int
	rds           *redis.Client
	JwtSecret     string
	MaxBufferSize int64
}

var Config Configuration = Configuration{
	Port:          8080,
	rds:           nil,
	JwtSecret:     "",
	MaxBufferSize: 32 << 20,
}
