package config

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Netflix/go-env"
	"github.com/go-ping/ping"
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

var Rds *redis.Client = nil

func init() {
	log.Info().Msg("App Start! load config from OS env")

	if _, err := env.UnmarshalFromEnviron(&Cfg); err != nil {
		log.Fatal().Err(err).Msg("Load config from env failed")
	}

	if Cfg.JwtFieldsKept != "" {
		Cfg.JwtFieldsKept = strings.ToLower(Cfg.JwtFieldsKept)
	}
	redisAddress := Cfg.RedisAddress
	//if redisAddress  is not of format address:port , add default port 6379
	if !strings.Contains(redisAddress, ":") {
		redisAddress = redisAddress + ":6379"
	}
	jsBytes, _ := json.Marshal(Cfg)
	log.Info().Str("Current Envs:", string(jsBytes)).Msg("Load config from env success")

	address := strings.Split(redisAddress, ":")[0]
	if len(address) == 0 {
		log.Fatal().Msg("RedisAddress is empty")
	}
	log.Info().Str("Start checking redis connection", address)
	//ping the address of redisAddress, if failed, print to log
	go pingServer(address)
	//apply configuration
	redisOption := &redis.Options{
		Addr:        Cfg.RedisAddress,
		Password:    Cfg.RedisPassword, // no password set
		DB:          Cfg.RedisDb,       // use default DB
		PoolSize:    200,
		DialTimeout: time.Second * 10,
	}
	rds := redis.NewClient(redisOption)
	//test connection
	if _, err := rds.Ping(context.Background()).Result(); err != nil {
		log.Fatal().Err(err).Msg("Redis connection failed: " + redisAddress)
	}
	Rds = rds

}
func pingServer(domain string) {
	pinger, err := ping.NewPinger(domain)
	if err != nil {
		log.Info().AnErr("ERROR Ping", err)
	}
	pinger.Count = 4
	pinger.Timeout = time.Second * 10

	pinger.OnRecv = func(pkt *ping.Packet) {
		fmt.Printf("Ping Received packet from %s: icmp_seq=%d time=%v\n",
			pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d Ping packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("Ping round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	fmt.Printf("start pinging %s", domain)
	if err := pinger.Run(); err != nil {
		log.Info().AnErr("ERROR Ping", err)
	}
}
