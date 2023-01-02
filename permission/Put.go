package permission

import (
	"context"
	"fmt"
	"time"

	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/logger"
)

type PutPermission struct {
	Key       string
	CreateAt  int64
	WhiteList []string
	BlackList []string
}

var PermittedPutOp map[string]PutPermission = make(map[string]PutPermission)
var lastLoadPutPermissionInfo string = ""

func LoadPutPermissionFromRedis() (err error) {
	var (
		Permissions_TMP map[string]string = map[string]string{}
		KeyNum          int64             = 0
	)
	// read RedisPutPermission usiing ParamRds
	// RedisPutPermission is a hash
	// split each value of RedisPutPermission into string[] and store in PermittedPutOp
	Permissions_TMP, err = config.ParamRds.HGetAll(context.Background(), "RedisPutPermission").Result()
	if err != nil {
		logger.Lshortfile.Println("loading RedisPutPermission  error: " + err.Error() + ". Consider Add hash item  RedisPutPermission in redis,with key redis key before ':' and value as permitted batch operations seperated by ','")
		return err
	}
	for k, v := range Permissions_TMP {
		//use msgpack to unmarshal v to PermitStatus
		var permission = PutPermission{CreateAt: time.Now().Unix()}
		if msgpack.Unmarshal([]byte(v), &permission) != nil {
			logger.Lshortfile.Println("loading RedisPutPermission  error: " + err.Error() + ". Consider Add hash item  RedisPutPermission in redis,with key redis key before ':' and value as permitted batch operations seperated by ','")
			continue
		}
		KeyNum++
		PermittedPutOp[k] = permission

	}
	//print info like this: info := fmt.Sprint("loading RedisPutPermission success. num keys:%d PermittedPutOp size:%d", KeyNum, len(PermittedPutOp))
	info := fmt.Sprint("loading RedisPutPermission success. num keys:", KeyNum)
	if info != lastLoadPutPermissionInfo {
		logger.Lshortfile.Println(info)
		lastLoadPutPermissionInfo = info
	}
	return nil
}
func RefreshPutPermission() {
	for {
		LoadPutPermissionFromRedis()
		time.Sleep(time.Second * 10)
	}
}
func IsPermittedPutOperation(dataKey string, operation string) bool {
	batchPermission, ok := PermittedPutOp[dataKey]
	//if datakey not in BatchPermission, then create BatchPermission, and add it to BatchPermission in redis
	if !ok {
		//auto set default permission, the operation is dis-allowed
		batchPermission = PutPermission{Key: dataKey, CreateAt: time.Now().Unix(), WhiteList: []string{}, BlackList: []string{operation}}
		PermittedPutOp[dataKey] = batchPermission
		//save to redis
		if b, err := msgpack.Marshal(batchPermission); err == nil {
			config.ParamRds.HSet(context.Background(), "RedisPutPermission", dataKey, string(b))
		}
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
	if b, err := msgpack.Marshal(batchPermission); err == nil {
		config.ParamRds.HSet(context.Background(), "RedisPutPermission", dataKey, string(b))
	}
	return false
}
