package permission

import (
	"context"
	"fmt"
	"time"

	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/logger"
)

type BatchPermission struct {
	Key      string
	CreateAt int64
	Actions  []string
}

var PermittedBatchOp map[string]BatchPermission = make(map[string]BatchPermission)
var lastLoadRedisBatchOpPermissionInfo string = ""

func LoadRedisBatchOpPermissionFromRedis() (err error) {
	var (
		Permissions_TMP map[string]string = map[string]string{}
		KeyNum          int64             = 0
	)
	// read RedisBatchOpPermission usiing ParamRds
	// RedisBatchOpPermission is a hash
	// split each value of RedisBatchOpPermission into string[] and store in PermittedBatchOp
	Permissions_TMP, err = config.ParamRds.HGetAll(context.Background(), "RedisBatchOpPermission").Result()
	if err != nil {
		logger.Lshortfile.Println("loading RedisBatchOpPermission  error: " + err.Error() + ". Consider Add hash item  RedisBatchOpPermission in redis,with key redis key before ':' and value as permitted batch operations seperated by ','")
		return err
	}
	for k, v := range Permissions_TMP {
		//use msgpack to unmarshal v to PermitStatus
		var batchPermission = BatchPermission{CreateAt: time.Now().Unix()}
		if msgpack.Unmarshal([]byte(v), &batchPermission) != nil {
			logger.Lshortfile.Println("loading RedisBatchOpPermission  error: " + err.Error() + ". Consider Add hash item  RedisBatchOpPermission in redis,with key redis key before ':' and value as permitted batch operations seperated by ','")
			continue
		}
		KeyNum++
		PermittedBatchOp[k] = batchPermission

	}
	//print info like this: info := fmt.Sprint("loading RedisBatchOpPermission success. num keys:%d PermittedBatchOp size:%d", KeyNum, len(PermittedBatchOp))
	info := fmt.Sprint("loading RedisBatchOpPermission success. num keys:", KeyNum)
	if info != lastLoadRedisBatchOpPermissionInfo {
		logger.Lshortfile.Println(info)
		lastLoadRedisBatchOpPermissionInfo = info
	}
	return nil
}
func RefreshRedisBatchOpPermission() {
	for {
		LoadRedisBatchOpPermissionFromRedis()
		time.Sleep(time.Second * 10)
	}
}
func IsPermittedBatchOperation(dataKey string, operation string) bool {
	batchPermission, ok := PermittedBatchOp[dataKey]
	//if datakey not in BatchPermission, then create BatchPermission, and add it to BatchPermission in redis
	if !ok {
		batchPermission = BatchPermission{Key: dataKey, CreateAt: time.Now().Unix(), Actions: []string{operation}}
		PermittedBatchOp[dataKey] = batchPermission
		//save to redis
		if b, err := msgpack.Marshal(batchPermission); err == nil {
			config.ParamRds.HSet(context.Background(), "RedisBatchOpPermission", dataKey, string(b))
		}
		return true
	}
	//caes ok

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
		if b, err := msgpack.Marshal(batchPermission); err == nil {
			config.ParamRds.HSet(context.Background(), "RedisBatchOpPermission", dataKey, string(b))
		}
		return true
	}

	return false
}
