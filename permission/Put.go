package permission

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/yangkequn/saavuu/api"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/logger"
)

var PermittedPutOp map[string]Permission = make(map[string]Permission)
var lastLoadPutPermissionInfo string = ""

func LoadPutPermissionFromRedis() {
	var _map map[string]Permission = make(map[string]Permission)
	// read RedisPutPermission usiing ParamRds
	// RedisPutPermission is a hash
	// split each value of RedisPutPermission into string[] and store in PermittedPutOp

	paramCtx := api.NewContext(context.Background())
	if err := paramCtx.HGetAll("RedisPutPermission", _map); err != nil {
		logger.Lshortfile.Println("loading RedisPutPermission  error: " + err.Error())
	} else {
		if info := fmt.Sprint("loading RedisPutPermission success. num keys:", len(_map)); info != lastLoadPutPermissionInfo {
			logger.Lshortfile.Println(info)
			lastLoadPutPermissionInfo = info
		}
		PermittedPutOp = _map
	}
	time.Sleep(time.Second * 10)
	go LoadPutPermissionFromRedis()
}
func IsPutPermitted(dataKey string, operation string) bool {
	dataKey = strings.Split(dataKey, ":")[0]
	permission, ok := PermittedPutOp[dataKey]
	//if datakey not in BatchPermission, then create BatchPermission, and add it to BatchPermission in redis
	if !ok {
		permission = Permission{Key: dataKey, CreateAt: time.Now().Unix(), WhiteList: []string{}, BlackList: []string{}}
	}

	//return true if allowed
	for _, v := range permission.WhiteList {
		if v == operation {
			return true
		}
	}
	// if operation not in black list, then add it to black list
	for _, v := range permission.BlackList {
		if v == operation {
			return false
		}
	}

	// if using develop mode, then add operation to white list; else add operation to black list
	if config.Cfg.DevelopMode {
		permission.WhiteList = append(permission.WhiteList, operation)
	} else {
		permission.BlackList = append(permission.BlackList, operation)
	}
	PermittedPutOp[dataKey] = permission
	//save to redis
	paramCtx := api.NewContext(context.Background())
	paramCtx.HSet("RedisPutPermission", dataKey, permission)
	return config.Cfg.DevelopMode
}
