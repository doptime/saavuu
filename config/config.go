package config

import (
	"encoding/json"
	"strings"

	"github.com/Netflix/go-env"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Configuration struct {
	RPCFirst bool `env:"RPCFirst,default=false"`
	//on redis server
	RedisAddress  string `env:"RedisAddress,required=true"`
	RedisPassword string `env:"RedisPassword"`
	RedisDb       int    `env:"RedisDb,required=true"`

	//AutoPermission should never be true in production
	AutoPermission bool   `env:"AutoPermission"`
	JWTSecret      string `env:"JWTSecret"`
	JwtFieldsKept  string `env:"JwtFieldsKept"`
	MaxBufferSize  int64  `env:"MaxBufferSize,default=10*1024*1024"`
	CORS           string `env:"CORS,default=*"`

	ServerPort int64  `env:"ServerPort,default=8000"`
	ServerPath string `env:"ServerPath,default=/"`

	//ServiceBatchSize is the number of tasks that a service can read from redis at the same time
	ServiceBatchSize int64 `env:"ServiceBatchSize"`
}

var Cfg Configuration = Configuration{}

var Rds *redis.Client

func init() {
	log.Info().Msg("App Start! load config from OS env")

	if _, err := env.UnmarshalFromEnviron(&Cfg); err != nil {
		log.Fatal().Err(err).Msg("Load config from env failed")
	}

	if Cfg.JwtFieldsKept != "" {
		Cfg.JwtFieldsKept = strings.ToLower(Cfg.JwtFieldsKept)
	}
	//apply configuration
	redisOption := &redis.Options{
		Addr:     Cfg.RedisAddress,
		Password: Cfg.RedisPassword, // no password set
		DB:       Cfg.RedisDb,       // use default DB
		PoolSize: 200,
	}
	Rds = redis.NewClient(redisOption)

	jsBytes, _ := json.Marshal(Cfg)
	log.Info().Str("Current Envs:", string(jsBytes)).Msg("Load config from env success")
}
