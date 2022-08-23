package saavuu

import (
	"github.com/go-redis/redis/v8"
)

type Configuration struct {
	Port          int
	rds           *redis.Client
	JwtToken      string
	MaxBufferSize int64
}

var Config Configuration = Configuration{
	Port:          8080,
	rds:           nil,
	JwtToken:      "",
	MaxBufferSize: 32 << 20,
}
