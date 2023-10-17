package permission

import (
	"context"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/rs/zerolog/log"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/data"
)

var ConfigurationLoaded bool = false
var PermitMaps []cmap.ConcurrentMap[string, *Permission] = []cmap.ConcurrentMap[string, *Permission]{cmap.New[*Permission](), cmap.New[*Permission](), cmap.New[*Permission](), cmap.New[*Permission]()}
var PermitKeys []string = []string{"saavuuPermissionPost", "saavuuPermissionPut", "saavuuPermissionGet", "saavuuPermissionDel"}

type PermitType int64

const (
	Post PermitType = 0
	Put  PermitType = 1
	Get  PermitType = 2
	Del  PermitType = 3
)

func LoadPermissionFromRedis() {
	var (
		err  error
		_map map[string]*Permission
	)
	//wait while config.Rds is nil
	for config.Rds == nil {
		time.Sleep(time.Millisecond * 10)
	}

	for i, key := range PermitKeys {
		var paramRds = data.Ctx[string, *Permission]{Rds: config.Rds, Ctx: context.Background(), Key: key}
		if _map, err = paramRds.HGetAll(); err != nil {
			log.Warn().Str("key", key).Any("num", len(_map)).Err(err).Msg("Load permission Failed")
			continue
		}
		if mapChanged := permitMapUpdate(_map, PermitMaps[i]); mapChanged {
			log.Info().Str("key", key).Any("num", len(_map)).Msg("Load permission success")
		}
	}
	if !ConfigurationLoaded {
		ConfigurationLoaded = true
		log.Info().Msg("Load Configuration Permission From Redis success!")
	}
	time.Sleep(time.Second * 10)
	go LoadPermissionFromRedis()
}
func init() {
	go LoadPermissionFromRedis()
}
