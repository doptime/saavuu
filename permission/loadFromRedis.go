package permission

import (
	"context"
	"fmt"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/rs/zerolog/log"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/data"
)

var lastLoadPermissionInfo map[string]string = make(map[string]string)
var permitKeyPut string = "saavuuPermissionPut"
var permitKeyPost string = "saavuuPermissionPost"
var permitKeyGet string = "saavuuPermissionGet"
var permitKeyDel string = "saavuuPermissionDel"

func LoadPPermissionFromRedis() {
	var err error
	//wait while config.Rds is nil
	for config.Rds == nil {
		time.Sleep(time.Millisecond * 100)
	}
	// read RedisPutPermission usiing Rds
	// RedisPutPermission is a hash
	// split each value of RedisPutPermission into string[] and store in PermittedPutOp

	//a slice name key holding  "RedisPutPermission","RedisPostPermission","RedisGetPermission","RedisDeletePermission"
	var keys []string = []string{permitKeyPut, permitKeyPost, permitKeyGet, permitKeyDel}
	//a slice name desMap holding  &PermittedPutOp,&PermittedPostOp,&PermittedGetOp,&PermittedDelOp
	var desMap []cmap.ConcurrentMap[string, *Permission] = []cmap.ConcurrentMap[string, *Permission]{PermittedPutOp, PermittedPostOp, PermittedGetOp, PermittedDelOp}
	for i, key := range keys {
		var _map map[string]*Permission
		var paramRds = data.Ctx[string, *Permission]{Rds: config.Rds, Ctx: context.Background(), Key: key}
		if _map, err = paramRds.HGetAll(); err != nil {
			log.Info().Msg("loading " + key + "  error: " + err.Error())
		} else {
			var mapDes cmap.ConcurrentMap[string, *Permission] = cmap.New[*Permission]()
			mapDes.MSet(_map)
			lastInfo, ok := lastLoadPermissionInfo[key]
			if info := fmt.Sprint("loading "+key+" success. num keys:", mapDes.Count()); !ok || info != lastInfo {
				log.Info().Msg(info)
				lastLoadPermissionInfo[key] = info
			}
			desMap[i] = mapDes
		}
	}
	go ContinuousReloadPermission()
}
func ContinuousReloadPermission() {
	time.Sleep(time.Second * 10)
	LoadPPermissionFromRedis()
}
