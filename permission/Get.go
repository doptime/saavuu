package permission

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/yangkequn/saavuu"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/logger"
)

var PermittedBatchOp map[string]Permission = make(map[string]Permission)
var lastLoadRedisGetPermissionInfo string = ""

func LoadGetPermissionFromRedis() {
	var _map map[string]Permission = make(map[string]Permission)
	// read RedisGetPermission usiing ParamRds
	// RedisGetPermission is a hash
	// split each value of RedisGetPermission into string[] and store in PermittedBatchOp
	paramCtx := saavuu.NewApiContext(context.Background())
	if err := paramCtx.HGetAll("RedisGetPermission", _map); err != nil {
		logger.Lshortfile.Println("loading RedisGetPermission  error: " + err.Error())
	} else {
		if info := fmt.Sprint("loading RedisGetPermission success. num keys:", len(_map)); info != lastLoadRedisGetPermissionInfo {
			logger.Lshortfile.Println(info)
			lastLoadRedisGetPermissionInfo = info
		}
		PermittedBatchOp = _map
	}
	time.Sleep(time.Second * 10)
	go LoadGetPermissionFromRedis()
}

// Only Batch Get Operation is checked. HGET etc are not checked
func IsGetPermitted(dataKey string, operation string) bool {
	dataKey = strings.Split(dataKey, ":")[0]
	permission, ok := PermittedBatchOp[dataKey]
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
	//return false if not allowed
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

	PermittedBatchOp[dataKey] = permission
	//save to redis
	paramCtx := saavuu.NewApiContext(context.Background())
	paramCtx.HSet("RedisGetPermission", dataKey, permission)
	return config.Cfg.DevelopMode
}
