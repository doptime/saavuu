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
	cmap "github.com/orcaman/concurrent-map/v2"
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
	Name     string
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
	Redis []*ConfigRedis `env:"REDIS,required=true"`
	Jwt   ConfigJWT      `env:"JWT"`
	Http  ConfigHttp     `env:"HTTP"`
	Api   ConfigAPI      `env:"API"`
	//{"DebugLevel": 0,"InfoLevel": 1,"WarnLevel": 2,"ErrorLevel": 3,"FatalLevel": 4,"PanicLevel": 5,"NoLevel": 6,"Disabled": 7	  }
	LogLevel int8 `env:"LogLevel,default=1"`
}

// set default values
var Cfg Configuration = Configuration{
	Redis:    []*ConfigRedis{},
	Jwt:      ConfigJWT{Secret: "", Fields: ""},
	Http:     ConfigHttp{CORES: "*", Port: 80, Path: "/", Enable: false, MaxBufferSize: 10485760},
	Api:      ConfigAPI{RPCFirst: false, AutoPermission: false, ServiceBatchSize: 64},
	LogLevel: 1,
}

var Rds map[string]*redis.Client = map[string]*redis.Client{}

func RdsClientDefault() *redis.Client {
	var (
		ok  bool
		rds *redis.Client
	)
	if rds, ok = Rds[""]; !ok {
		log.Panic().Msg("default redis client not found")
	}
	return rds
}
func RdsClientByName(name string) (rds *redis.Client, err error) {
	var (
		ok bool
	)
	if rds, ok = Rds[name]; !ok {
		err = fmt.Errorf("redis client with name %s not found", name)
		return nil, err
	}

	return rds, nil
}

func LoadConfig() (err error) {
	if redisEnv := os.Getenv("Redis"); len(redisEnv) > 0 {
		//regex trim space in "} , {"
		for len(strings.Replace(redisEnv, " {", "{", -1)) != len(redisEnv) {
			redisEnv = strings.Replace(redisEnv, " {", "{", -1)
		}
		for len(strings.Replace(redisEnv, "} ", "}", -1)) != len(redisEnv) {
			redisEnv = strings.Replace(redisEnv, "} ", "}", -1)
		}
		redisEnv = strings.Replace(redisEnv, "},{", "}},{{", -1)
		rdsStrs := strings.Split(redisEnv, "},{")
		for _, rdsStr := range rdsStrs {
			rdsCfg := &ConfigRedis{}
			if err := json.Unmarshal([]byte(rdsStr), &rdsCfg); err != nil {
				correctFormat := "{Name,Username,Password,Host,Port,DB},{Name,Username,Password,Host,Port,DB}"
				log.Fatal().Err(err).Str("redisEnv", rdsStr).Msg("Step1.0 Load config from Redis env failed, correct format: " + correctFormat)
			}
			if rdsCfg.Name == "default" || rdsCfg.Name == "_" {
				rdsCfg.Name = ""
			}
			Cfg.Redis = append(Cfg.Redis, rdsCfg)
		}
	}

	// Load and parse JWT config
	if jwtEnv := os.Getenv("Jwt"); jwtEnv != "" {
		if err := json.Unmarshal([]byte(jwtEnv), &Cfg.Jwt); err != nil {
			log.Fatal().Err(err).Str("jwtEnv", jwtEnv).Msg("Step1.0 Load config from JWT env failed")
		}
	}

	// Load and parse HTTP config
	Cfg.Http.Enable, Cfg.Http.Path, Cfg.Http.CORES = true, "/", "*"
	if httpEnv := os.Getenv("Http"); len(httpEnv) > 0 {
		if err := json.Unmarshal([]byte(httpEnv), &Cfg.Http); err != nil {
			log.Fatal().Err(err).Str("httpEnv", httpEnv).Msg("Step1.0 Load config from HTTP env failed")
		}
	}

	// Load and parse API config
	if apiEnv := os.Getenv("Api"); apiEnv != "" {
		if err := json.Unmarshal([]byte(apiEnv), &Cfg.Api); err != nil {
			log.Fatal().Err(err).Str("apiEnv", apiEnv).Msg("Step1.0 Load config from API env failed")
		}
	}

	// Load LogLevel
	if logLevelEnv := os.Getenv("LogLevel"); len(logLevelEnv) > 0 {
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

	for _, rdsCfg := range Cfg.Redis {
		//apply configuration
		redisOption := &redis.Options{
			Addr:         rdsCfg.Host + ":" + rdsCfg.Port,
			Username:     rdsCfg.Username,
			Password:     rdsCfg.Password, // no password set
			DB:           int(rdsCfg.DB),  // use default DB
			PoolSize:     200,
			DialTimeout:  time.Second * 10,
			ReadTimeout:  -1,
			WriteTimeout: time.Second * 300,
		}
		rdsClient := redis.NewClient(redisOption)
		//test connection
		if _, err := rdsClient.Ping(context.Background()).Result(); err != nil {
			log.Fatal().Err(err).Any("Step1.3 Redis server not rechable", rdsCfg.Host).Send()
			return //if redis server is not valid, exit
		}
		//save to the list
		log.Info().Str("Step1.3 Redis Load ", "Success").Any("RedisUsername", rdsCfg.Username).Any("RedisPassword", rdsCfg.Password).Any("RedisHost", rdsCfg.Host).Any("RedisPort", rdsCfg.Port).Send()
		Rds[rdsCfg.Name] = rdsClient
		timeCmd := rdsClient.Time(context.Background())
		log.Info().Any("Step1.4 Redis server time: ", timeCmd.Val().String()).Send()
		//ping the address of redisAddress, if failed, print to log
		pingServer(rdsCfg.Host)

	}

	log.Info().Msg("Step1.E: App loaded done")

}

var pingTaskServers = cmap.ConcurrentMap[string, bool]{}

func pingServer(domain string) {
	var (
		pinger *ping.Pinger
		err    error
	)
	for _, ok := pingTaskServers.Get(domain); ok; {
		return
	}
	pingTaskServers.Set(domain, true)

	if pinger, err = ping.NewPinger(domain); err != nil {
		log.Info().AnErr("Step1.5 ERROR NewPinger", err).Send()
	}
	pinger.Count = 4
	pinger.Timeout = time.Second * 10
	pinger.OnRecv = func(pkt *ping.Packet) {}

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
	go func() {
		if err := pinger.Run(); err != nil {
			log.Info().AnErr("Step1.5 ERROR Ping", err).Send()
		}
	}()
}
