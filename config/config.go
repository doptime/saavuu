package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-ping/ping"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

// set default values
var Cfg Configuration = Configuration{
	Redis:    ConfigRedis{Username: "", Password: "", Host: "", Port: "6379", DB: 0},
	Jwt:      ConfigJWT{Secret: "", Fields: ""},
	Http:     ConfigHttp{CORES: "*", Port: 80, Path: "/", Enable: false, MaxBufferSize: 10485760},
	Api:      ConfigAPI{RPCFirst: false, AutoPermission: false, ServiceBatchSize: 64},
	LogLevel: 1,
}

var Rds *redis.Client = nil

func LoadConfig() (err error) {
	// Load and parse Redis config
	redisEnv, jwtEnv, httpEnv, apiEnv, logLevelEnv := os.Getenv("Redis"), os.Getenv("Jwt"), os.Getenv("Http"), os.Getenv("Api"), os.Getenv("LogLevel")
	if redisEnv == "" || jwtEnv == "" || apiEnv == "" {
		return fmt.Errorf("Step1.0 Load config from env failed")
	}

	if err := json.Unmarshal([]byte(redisEnv), &Cfg.Redis); err != nil {
		log.Fatal().Err(err).Msg("Step1.0 Load config from Redis env failed")
	}

	// Load and parse JWT config
	if err := json.Unmarshal([]byte(jwtEnv), &Cfg.Jwt); err != nil {
		log.Fatal().Err(err).Msg("Step1.0 Load config from JWT env failed")
	}

	// Load and parse HTTP config
	Cfg.Http.Enable, Cfg.Http.Path, Cfg.Http.CORES = true, "/", "*"
	if len(httpEnv) > 0 {
		if err := json.Unmarshal([]byte(httpEnv), &Cfg.Http); err != nil {
			log.Fatal().Err(err).Msg("Step1.0 Load config from HTTP env failed")
		}
	}

	// Load and parse API config

	if err := json.Unmarshal([]byte(apiEnv), &Cfg.Api); err != nil {
		log.Fatal().Err(err).Msg("Step1.0 Load config from API env failed")
	}

	// Load LogLevel
	if len(logLevelEnv) > 0 {
		if logLevel, err := strconv.ParseInt(logLevelEnv, 10, 8); err == nil {
			Cfg.LogLevel = int8(logLevel)
		}
	}
	return nil
}
func init() {
	log.Info().Msg("Step1.0: App Start! load config from OS env")
	if err := LoadConfig(); err != nil {
		log.Info().AnErr("Step1.0 ERROR LoadConfig", err).Send()
		log.Info().Msg("saavuu data & api will no be able to be used. please check your env and restart the app if you want to use itß")
		return
	}
	zerolog.SetGlobalLevel(zerolog.Level(Cfg.LogLevel))

	if Cfg.Jwt.Fields != "" {
		Cfg.Jwt.Fields = strings.ToLower(Cfg.Jwt.Fields)
	}
	log.Info().Any("Step1.1 Current Envs:", Cfg).Msg("Load config from env success")

	log.Info().Str("Step1.2 Checking Redis", "Start").Send()

	//apply configuration
	redisOption := &redis.Options{
		Addr:         Cfg.Redis.Host + ":" + Cfg.Redis.Port,
		Username:     Cfg.Redis.Username,
		Password:     Cfg.Redis.Password, // no password set
		DB:           int(Cfg.Redis.DB),  // use default DB
		PoolSize:     200,
		DialTimeout:  time.Second * 10,
		ReadTimeout:  time.Second * 300,
		WriteTimeout: time.Second * 300,
	}
	Rds = redis.NewClient(redisOption)
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
