package permission

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/yangkequn/saavuu"
	"github.com/yangkequn/saavuu/logger"
)

type Permission struct {
	Key       string
	CreateAt  int64
	WhiteList []string
	BlackList []string
}

var PermittedDelOp map[string]Permission = make(map[string]Permission)
var lastLoadDelPermissionInfo string = ""

func LoadDelPermissionFromRedis() {
	// read RedisDelPermission usiing ParamRds
	// RedisDelPermission is a hash
	// split each value of RedisDelPermission into string[] and store in PermittedDelOp
	paramCtx := saavuu.NewParamContext(context.Background())
	var mapTmp map[string]Permission = make(map[string]Permission)
	if err := paramCtx.HGetAll("RedisDelPermission", mapTmp); err != nil {
		logger.Lshortfile.Println("loading RedisDelPermission  error: " + err.Error() + ". Consider Add hash item  RedisDelPermission in redis,with key redis key before ':' and value as permitted batch operations seperated by ','")

		time.Sleep(time.Second * 10)
		go LoadPutPermissionFromRedis()
		return
	}
	//print info like this: info := fmt.Sprint("loading RedisDelPermission success. num keys:%d PermittedDelOp size:%d", KeyNum, len(PermittedDelOp))
	info := fmt.Sprint("loading RedisDelPermission success. num keys:", len(mapTmp))
	if info != lastLoadDelPermissionInfo {
		logger.Lshortfile.Println(info)
		lastLoadDelPermissionInfo = info
	}
	PermittedDelOp = mapTmp
	time.Sleep(time.Second * 10)
	go LoadPutPermissionFromRedis()
}
func IsPermittedDelOperation(dataKey string, operation string) bool {
	// remove :... from dataKey
	dataKey = strings.Split(dataKey, ":")[0]

	batchPermission, ok := PermittedDelOp[dataKey]
	//if datakey not in BatchPermission, then create BatchPermission, and add it to BatchPermission in redis
	if !ok {
		//auto set default permission, the operation is dis-allowed
		batchPermission = Permission{Key: dataKey, CreateAt: time.Now().Unix(), WhiteList: []string{}, BlackList: []string{operation}}
		PermittedDelOp[dataKey] = batchPermission
		//save to redis
		paramCtx := saavuu.NewParamContext(context.Background())
		paramCtx.HSet("RedisDelPermission", dataKey, batchPermission)
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
	paramCtx.HSet("RedisDelPermission", dataKey, batchPermission)
	return false
}
