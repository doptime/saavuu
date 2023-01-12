package permission

import (
	"context"
	"fmt"
	"time"

	"github.com/yangkequn/saavuu"
	"github.com/yangkequn/saavuu/logger"
)

type PutPermission struct {
	Key       string
	CreateAt  int64
	WhiteList []string
	BlackList []string
}

var PermittedPutOp map[string]Permission = make(map[string]Permission)
var lastLoadPutPermissionInfo string = ""

func LoadPutPermissionFromRedis() {
	// read RedisPutPermission usiing ParamRds
	// RedisPutPermission is a hash
	// split each value of RedisPutPermission into string[] and store in PermittedPutOp

	paramCtx := saavuu.NewParamContext(context.Background())
	if err := paramCtx.HGetAll("RedisPutPermission", PermittedPutOp); err != nil {
		logger.Lshortfile.Println("loading RedisPutPermission  error: " + err.Error() + ". Consider Add hash item  RedisPutPermission in redis,with key redis key before ':' and value as permitted batch operations seperated by ','")
		time.Sleep(time.Second * 10)
		go LoadPutPermissionFromRedis()
		return
	}
	//print info like this: info := fmt.Sprint("loading RedisPutPermission success. num keys:%d PermittedPutOp size:%d", KeyNum, len(PermittedPutOp))
	info := fmt.Sprint("loading RedisPutPermission success. num keys:", len(PermittedPutOp))
	if info != lastLoadPutPermissionInfo {
		logger.Lshortfile.Println(info)
		lastLoadPutPermissionInfo = info
	}
	time.Sleep(time.Second * 10)
	go LoadPutPermissionFromRedis()
}
func IsPermittedPutOperation(dataKey string, operation string) bool {
	batchPermission, ok := PermittedPutOp[dataKey]
	//if datakey not in BatchPermission, then create BatchPermission, and add it to BatchPermission in redis
	if !ok {
		//auto set default permission, the operation is dis-allowed
		batchPermission = Permission{Key: dataKey, CreateAt: time.Now().Unix(), WhiteList: []string{}, BlackList: []string{operation}}
		PermittedPutOp[dataKey] = batchPermission
		//save to redis
		paramCtx := saavuu.NewParamContext(context.Background())
		paramCtx.HSet("RedisPutPermission", dataKey, batchPermission)
		return false
	}
	//caes ok

	//return true if allowed
	for _, v := range batchPermission.WhiteList {
		if v == operation {
			return true
		}
	}
	// if operation not in black list, then add it to black list
	for _, v := range batchPermission.BlackList {
		if v == operation {
			return false
		}
	}
	//add operation to black list for convenient modification
	batchPermission.BlackList = append(batchPermission.BlackList, operation)
	//save to redis
	paramCtx := saavuu.NewParamContext(context.Background())
	paramCtx.HSet("RedisPutPermission", dataKey, batchPermission)
	return false
}
