package permission

import (
	"context"
	"fmt"
	"time"

	"github.com/yangkequn/saavuu"
	"github.com/yangkequn/saavuu/logger"
)

type BatchPermission struct {
	Key      string
	CreateAt int64
	Actions  []string
}

var PermittedBatchOp map[string]BatchPermission = make(map[string]BatchPermission)
var lastLoadRedisBatchOpPermissionInfo string = ""

func LoadGetBatchPermissionFromRedis() {
	// read RedisBatchOpPermission usiing ParamRds
	// RedisBatchOpPermission is a hash
	// split each value of RedisBatchOpPermission into string[] and store in PermittedBatchOp

	paramCtx := saavuu.NewParamContext(context.Background())
	if err := paramCtx.HGetAll("RedisBatchOpPermission", PermittedBatchOp); err != nil {
		logger.Lshortfile.Println("loading RedisBatchOpPermission  error: " + err.Error() + ". Consider Add hash item  RedisBatchOpPermission in redis,with key redis key before ':' and value as permitted batch operations seperated by ','")
		time.Sleep(time.Second * 10)
		go LoadGetBatchPermissionFromRedis()
		return
	}

	//print info like this: info := fmt.Sprint("loading RedisBatchOpPermission success. num keys:%d PermittedBatchOp size:%d", KeyNum, len(PermittedBatchOp))
	info := fmt.Sprint("loading RedisBatchOpPermission success. num keys:", len(PermittedBatchOp))
	if info != lastLoadRedisBatchOpPermissionInfo {
		logger.Lshortfile.Println(info)
		lastLoadRedisBatchOpPermissionInfo = info
	}
	time.Sleep(time.Second * 10)
	go LoadGetBatchPermissionFromRedis()
}
func IsPermittedBatchOperation(dataKey string, operation string) bool {
	batchPermission, ok := PermittedBatchOp[dataKey]
	//if datakey not in BatchPermission, then create BatchPermission, and add it to BatchPermission in redis
	if !ok {
		batchPermission = BatchPermission{Key: dataKey, CreateAt: time.Now().Unix(), Actions: []string{operation}}
		PermittedBatchOp[dataKey] = batchPermission
		//save to redis
		paramCtx := saavuu.NewParamContext(context.Background())
		paramCtx.HSet("RedisBatchOpPermission", dataKey, batchPermission)
		return true
	}

	//return true if allowed
	for _, v := range batchPermission.Actions {
		if v == operation {
			return true
		}
	}

	//auto set default permission
	// if CeatedAt within 3 days, the operation is allowed
	// after 3 days, the permission is locked, and the operation is not allowed
	if time.Now().Unix()-batchPermission.CreateAt < 3*24*3600 {
		batchPermission.Actions = append(batchPermission.Actions, operation)
		PermittedBatchOp[dataKey] = batchPermission
		//save to redis
		paramCtx := saavuu.NewParamContext(context.Background())
		paramCtx.HSet("RedisBatchOpPermission", dataKey, batchPermission)
		return true
	}

	return false
}
