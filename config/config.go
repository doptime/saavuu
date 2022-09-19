package config

import (
	"github.com/go-redis/redis/v8"
)

type Configuration struct {
	// Parameter Server Should be Memory Only server, with high bandwidth and low latency.
	// All parameter from web client are post to this redis server first
	ParamRedis *redis.Client
	// DataRedis usually slower But with Flash Storage support ,such as Pikadb, and later may be KeyDB or DragonflyDB
	// Default redis server to read data from and write data to web client
	DataRedis     *redis.Client
	JwtSecret     string
	MaxBufferSize int64
	CORS          string
}

var Cfg Configuration = Configuration{
	ParamRedis:    nil,
	JwtSecret:     "",
	MaxBufferSize: 32 << 20,
	CORS:          "*",
}
