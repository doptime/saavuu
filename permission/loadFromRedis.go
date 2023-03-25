package permission

import (
	"context"
	"fmt"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/data"
	"github.com/yangkequn/saavuu/logger"
)

var lastLoadPermissionInfo map[string]string = make(map[string]string)
var permitKeyPut string = "RedisPutPermission"
var permitKeyPost string = "RedisPostPermission"
var permitKeyGet string = "RedisGetPermission"
var permitKeyDel string = "RedisDeletePermission"

func LoadPPermissionFromRedis() {
	// read RedisPutPermission usiing ParamRds
	// RedisPutPermission is a hash
	// split each value of RedisPutPermission into string[] and store in PermittedPutOp

	//a slice name key holding  "RedisPutPermission","RedisPostPermission","RedisGetPermission","RedisDeletePermission"
	var keys []string = []string{permitKeyPut, permitKeyPost, permitKeyGet, permitKeyDel}
	//a slice name desMap holding  &PermittedPutOp,&PermittedPostOp,&PermittedGetOp,&PermittedDelOp
	var desMap []*cmap.ConcurrentMap[string, Permission] = []*cmap.ConcurrentMap[string, Permission]{&PermittedPutOp, &PermittedPostOp, &PermittedGetOp, &PermittedDelOp}
	for i, key := range keys {

		var _map map[string]Permission = make(map[string]Permission)
		var paramRds = data.Ctx[Permission]{Rds: config.ParamRds, Ctx: context.Background(), Key: key}
		if err := paramRds.HGetAll(_map); err != nil {
			logger.Lshortfile.Println("loading " + key + "  error: " + err.Error())
		} else {
			var mapDes cmap.ConcurrentMap[string, Permission] = cmap.New[Permission]()
			mapDes.MSet(_map)
			lastInfo, ok := lastLoadPermissionInfo[key]
			if info := fmt.Sprint("loading "+key+" success. num keys:", mapDes.Count()); !ok || info != lastInfo {
				logger.Lshortfile.Println(info)
				lastLoadPermissionInfo[key] = info
			}
			*desMap[i] = mapDes
		}
	}
	time.Sleep(time.Second * 10)
	go LoadPPermissionFromRedis()
}
