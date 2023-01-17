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
var lastLoadRedisBatchOpPermissionInfo string = ""

func LoadGetBatchPermissionFromRedis() {
	// read RedisBatchOpPermission usiing ParamRds
	// RedisBatchOpPermission is a hash
	// split each value of RedisBatchOpPermission into string[] and store in PermittedBatchOp

	var mapTmp map[string]Permission = make(map[string]Permission)
	paramCtx := saavuu.NewParamContext(context.Background())
	if err := paramCtx.HGetAll("RedisBatchOpPermission", mapTmp); err != nil {
		logger.Lshortfile.Println("loading RedisBatchOpPermission  error: " + err.Error() + ". Consider Add hash item  RedisBatchOpPermission in redis,with key redis key before ':' and value as permitted batch operations seperated by ','")
		time.Sleep(time.Second * 10)
		go LoadGetBatchPermissionFromRedis()
		return
	}

	//print info like this: info := fmt.Sprint("loading RedisBatchOpPermission success. num keys:%d PermittedBatchOp size:%d", KeyNum, len(PermittedBatchOp))
	info := fmt.Sprint("loading RedisBatchOpPermission success. num keys:", len(mapTmp))
	if info != lastLoadRedisBatchOpPermissionInfo {
		logger.Lshortfile.Println(info)
		lastLoadRedisBatchOpPermissionInfo = info
	}
	PermittedBatchOp = mapTmp
	time.Sleep(time.Second * 10)
	go LoadGetBatchPermissionFromRedis()
}
func IsPermittedBatchOperation(dataKey string, operation string) bool {
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
	paramCtx := saavuu.NewParamContext(context.Background())
	paramCtx.HSet("RedisBatchOpPermission", dataKey, permission)
	return config.Cfg.DevelopMode
}
