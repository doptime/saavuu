package saavuu

import (
	"github.com/go-redis/redis/v8"
)

type Configuration struct {
	rds           *redis.Client
	JwtSecret     string
	MaxBufferSize int64
}

var Config Configuration = Configuration{
	rds:           nil,
	JwtSecret:     "",
	MaxBufferSize: 32 << 20,
}
