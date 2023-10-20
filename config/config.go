package config

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Netflix/go-env"
	"github.com/go-ping/ping"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Configuration struct {
	//redis server, format: username:password@address:port/db
	Redis         string `env:"Redis,required=true"`
	RedisUsername string
	RedisPassword string
	//RedisAddress is the address of redis server, format: address:port
	RedisAddress string
	RedisDB      int64

	JWTSecret     string `env:"JWTSecret"`
	JwtFieldsKept string `env:"JwtFieldsKept"`
	//MaxBufferSize is the max size of a task in bytes, default 10M
	MaxBufferSize int64  `env:"MaxBufferSize,default=10485760"`
	CORS          string `env:"CORS,default=*"`

	//http server, format: address:port/path.
	//if HTTP is empty, http service will not be started
	HTTP     string `env:"HTTP,default=0.0.0.0:80/"`
	HTTPPort int64
	HTTPPath string

	//AutoPermission should never be true in production
	RPCFirst       bool `env:"RPCFirst,default=false"`
	AutoPermission bool `env:"AutoPermission,default=false"`

	//{"DebugLevel": 0,"InfoLevel": 1,"WarnLevel": 2,"ErrorLevel": 3,"FatalLevel": 4,"PanicLevel": 5,"NoLevel": 6,"Disabled": 7	  }
	LogLevel int8 `env:"LogLevel,default=1"`
	//ServiceBatchSize is the number of tasks that a service can read from redis at the same time
	ServiceBatchSize int64 `env:"ServiceBatchSize,default=64"`
}

func (c *Configuration) GetReidsClient() (Rds *redis.Client, err error) {
	var addressInfo string = c.Redis
	//step1 parse username and password, and address and port from c.Redis
	// username and password is optional, default empty
	// address  is required
	// port is optional, default 6379
	// db is optional, default 0
	// format: username:password@address:port/db
	if atIndex := strings.LastIndex(c.Redis, "@"); atIndex >= 0 {
		//parse username and password
		addressInfo = c.Redis[atIndex+1:]
		c.RedisUsername = ""
		c.RedisPassword = ""
		if strings.Contains(c.Redis[:atIndex], ":") {
			c.RedisUsername = strings.Split(c.Redis[:atIndex], ":")[0]
			c.RedisPassword = strings.Split(c.Redis[:atIndex], ":")[1]
		}
	}
	//read db
	if strings.Contains(addressInfo, "/") {
		//parse address and port
		addressInfo = strings.Split(addressInfo, "/")[0]
		if RedisDB := strings.Split(addressInfo, "/")[1]; RedisDB == "" {
			c.RedisDB = 0
		} else if c.RedisDB, err = strconv.ParseInt(RedisDB, 10, 64); err != nil {
			log.Error().Err(err).Msg("Redis db is not a number")
			return nil, err
		}
	}
	// read port and address
	if c.RedisAddress = addressInfo; c.RedisAddress == "" {
		log.Fatal().Msg("Redis address is empty")
	}
	if !strings.Contains(c.RedisAddress, ":") {
		c.RedisAddress += ":6379"
	}

	//apply configuration
	redisOption := &redis.Options{
		Addr:         Cfg.RedisAddress,
		Username:     Cfg.RedisUsername,
		Password:     Cfg.RedisPassword, // no password set
		DB:           int(Cfg.RedisDB),  // use default DB
		PoolSize:     200,
		DialTimeout:  time.Second * 10,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	}
	rds := redis.NewClient(redisOption)
	//step2 connect to redis server
	// RedisPassword string `env:"RedisPassword"`
	// RedisDb       int    `env:"RedisDb,required=true"`
	return rds, nil
}
func (c *Configuration) HTTPEnabled() bool {
	return len(c.HTTP) > 0
}

var Cfg Configuration = Configuration{}

var Rds *redis.Client = nil

func init() {
	var err error
	log.Info().Msg("Step1.0: App Start! load config from OS env")

	if _, err := env.UnmarshalFromEnviron(&Cfg); err != nil {
		log.Fatal().Err(err).Msg("Load config from env failed")
	}
	zerolog.SetGlobalLevel(zerolog.Level(Cfg.LogLevel))

	if Cfg.JwtFieldsKept != "" {
		Cfg.JwtFieldsKept = strings.ToLower(Cfg.JwtFieldsKept)
	}
	log.Info().Any("Step1.1 Current Envs:", Cfg).Msg("Load config from env success")

	log.Info().Msg("Step1.2 Start checking Redis ")
	if Rds, err = Cfg.GetReidsClient(); err != nil {
		log.Fatal().Err(err).Any("Step1.2.1 Redis Enviroment not valid format of [username:password@address:port/db]", Cfg.Redis).Send()
	}
	//test connection
	if _, err := Rds.Ping(context.Background()).Result(); err != nil {
		log.Fatal().Err(err).Any("Step1.3 Redis server not rechable", Cfg.Redis).Send()
		return //if redis server is not valid, exit
	}
	log.Info().Str("Step1.3 Redis connection Success", Cfg.RedisAddress).Send()
	timeCmd := Rds.Time(context.Background())
	log.Info().Any("Step1.4 Redis server time: ", timeCmd.Val().String()).Send()
	//ping the address of redisAddress, if failed, print to log
	go pingServer(strings.Split(Cfg.RedisAddress, ":")[0])

	if Cfg.HTTPEnabled() {
		ind := strings.Index(Cfg.HTTP, "/")
		HTTPPort := Cfg.HTTP
		if ind >= 0 {
			Cfg.HTTPPath = "/" + strings.Split(Cfg.HTTP, "/")[1]
			HTTPPort = strings.Split(Cfg.HTTP, "/")[0]
		}
		Cfg.HTTPPath = "/" + Cfg.HTTPPath
		if ind := strings.Index(HTTPPort, ":"); ind >= 0 {
			HTTPPort = HTTPPort[ind+1:]
			Cfg.HTTPPort, _ = strconv.ParseInt(HTTPPort, 10, 64)
		} else {
			Cfg.HTTPPort = 80
		}
	}
	log.Info().Any("Http Enabled", Cfg.HTTPEnabled()).Any("Http Port", Cfg.HTTPPort).Any("Http Path", Cfg.HTTPPath).Msg("Step1.6 Redis HTTP Load Completed! ")
	log.Info().Msg("Step1.E: App loaded configuration completed!")

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
