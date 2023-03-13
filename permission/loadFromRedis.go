package permission

import (
	"context"
	"fmt"
	"time"

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
	var desMap []*map[string]Permission = []*map[string]Permission{&PermittedPutOp, &PermittedPostOp, &PermittedGetOp, &PermittedDelOp}
	for i, key := range keys {

		var _map map[string]Permission = make(map[string]Permission)

		var paramRds = data.Ctx{Rds: config.ParamRds, Ctx: context.Background(), Key: key}
		if err := paramRds.HGetAll(_map); err != nil {
			logger.Lshortfile.Println("loading " + key + "  error: " + err.Error())
		} else {
			lastInfo, ok := lastLoadPermissionInfo[key]
			if info := fmt.Sprint("loading "+key+" success. num keys:", len(_map)); !ok || info != lastInfo {
				logger.Lshortfile.Println(info)
				lastLoadPermissionInfo[key] = info
			}
			*desMap[i] = _map
		}
	}
	time.Sleep(time.Second * 10)
	go LoadPPermissionFromRedis()
}
