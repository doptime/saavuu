package config

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-ping/ping"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type ConfigHttp struct {
	CORES  string `env:"CORES,default=*"`
	Port   int64  `env:"Port,default=80"`
	Path   string `env:"Path,default=/"`
	Enable bool   `env:"Enable,default=false"`
	//MaxBufferSize is the max size of a task in bytes, default 10M
	MaxBufferSize int64 `env:"MaxBufferSize,default=10485760"`
}
type ConfigRedis struct {
	Username string `env:"Username"`
	Password string `env:"Password"`
	Host     string `env:"Host,required=true"`
	Port     string `env:"Port,required=true"`
	DB       int64  `env:"DB,required=true"`
}
type ConfigJWT struct {
	Secret string `env:"Secret"`
	Fields string `env:"Fields"`
}
type ConfigAPI struct {
	RPCFirst bool `env:"RPCFirst,default=false"`
	//AutoPermission should never be true in production
	AutoPermission bool `env:"AutoPermission,default=false"`
	//ServiceBatchSize is the number of tasks that a service can read from redis at the same time
	ServiceBatchSize int64 `env:"ServiceBatchSize,default=64"`
}

type Configuration struct {
	//redis server, format: username:password@address:port/db
	Redis ConfigRedis `env:"REDIS,required=true"`
	Jwt   ConfigJWT   `env:"JWT"`
	Http  ConfigHttp  `env:"HTTP"`
	Api   ConfigAPI   `env:"API"`
	//{"DebugLevel": 0,"InfoLevel": 1,"WarnLevel": 2,"ErrorLevel": 3,"FatalLevel": 4,"PanicLevel": 5,"NoLevel": 6,"Disabled": 7	  }
	LogLevel int8 `env:"LogLevel,default=1"`
}

var Cfg Configuration = Configuration{}

var Rds *redis.Client = nil

func init() {
	log.Info().Msg("Step1.0: App Start! load config from OS env")

	if err := viper.Unmarshal(&Cfg); err != nil {
		log.Fatal().Err(err).Msg("Load config from env failed")
	}
	zerolog.SetGlobalLevel(zerolog.Level(Cfg.LogLevel))

	if Cfg.Jwt.Fields != "" {
		Cfg.Jwt.Fields = strings.ToLower(Cfg.Jwt.Fields)
	}
	log.Info().Any("Step1.1 Current Envs:", Cfg).Msg("Load config from env success")

	log.Info().Str("Step1.2 Checking Redis", "Start").Send()
	//test connection
	if _, err := Rds.Ping(context.Background()).Result(); err != nil {
		log.Fatal().Err(err).Any("Step1.3 Redis server not rechable", Cfg.Redis).Send()
		return //if redis server is not valid, exit
	}
	log.Info().Str("Step1.3 Redis Load ", "Success").Any("RedisUsername", Cfg.Redis.Username).Any("RedisPassword", Cfg.Redis.Password).Any("RedisHost", Cfg.Redis.Host).Any("RedisPort", Cfg.Redis.Port).Send()
	timeCmd := Rds.Time(context.Background())
	log.Info().Any("Step1.4 Redis server time: ", timeCmd.Val().String()).Send()
	//ping the address of redisAddress, if failed, print to log
	go pingServer(Cfg.Redis.Host)

	log.Info().Msg("Step1.E: App loaded done")

}
func pingServer(domain string) {
	pinger, err := ping.NewPinger(domain)
	if err != nil {
		log.Info().AnErr("Step1.5 ERROR NewPinger", err).Send()
	}
	pinger.Count = 4
	pinger.Timeout = time.Second * 10

	pinger.OnRecv = func(pkt *ping.Packet) {
		//fmt.Printf("Ping Received packet from %s: icmp_seq=%d time=%v\n",pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		// fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		log.Info().Str("Step1.5 Ping ", fmt.Sprintf("--- %s ping statistics ---", stats.Addr)).Send()
		// fmt.Printf("%d Ping packets transmitted, %d packets received, %v%% packet loss\n",
		// 	stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		log.Info().Str("Step1.5 Ping", fmt.Sprintf("%d/%d/%v%%", stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)).Send()

		// fmt.Printf("Ping round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
		// 	stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
		log.Info().Str("Step1.5 Ping", fmt.Sprintf("%v/%v/%v/%v", stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)).Send()
	}

	if err := pinger.Run(); err != nil {
		log.Info().AnErr("Step1.5 ERROR Ping", err).Send()
	}
}
