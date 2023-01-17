package config

import (
	"github.com/go-redis/redis/v9"
)

type Configuration struct {
	RedisAddressParam  string
	RedisPasswordParam string
	RedisDbParam       int

	RedisAddressData  string
	RedisPasswordData string
	RedisDbData       int

	//DevelopMode should never be true in production
	DevelopMode     bool
	JwtSecret       string
	JwtIgnoreFields string
	MaxBufferSize   int64
	CORS            string
}

var Cfg Configuration = Configuration{
	JwtSecret:     "",
	MaxBufferSize: 32 << 20,
	CORS:          "*",
}

// Parameter Server Should be Memory Only server, with high bandwidth and low latency.
// All parameter from web client are post to this redis server first
var ParamRds *redis.Client

// DataRds usually slower But with Flash Storage support ,such as Pikadb, and later may be KeyDB or DragonflyDB
// Default redis server to read data from and write data to web client
var DataRds *redis.Client
