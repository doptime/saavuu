package permission

import (
	"context"
	"strings"
	"time"

	"github.com/cespare/xxhash"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/logger"
)

var PermittedBatchOp map[uint64]interface{} = map[uint64]interface{}{}

func LoadDataItemBatchPermissionsFromRedis() (err error) {
	var (
		Permissions_TMP map[string]string = map[string]string{}
		KeyNum          int64             = 0
	)
	logger.Lshortfile.Println("start loading DataItemBatchPermissions from redis")
	// read DataItemBatchPermissions usiing ParamRds
	// DataItemBatchPermissions is a hash
	// split each value of DataItemBatchPermissions into string[] and store in PermittedBatchOp
	Permissions_TMP, err = config.ParamRds.HGetAll(context.Background(), "DataItemBatchPermissions").Result()
	if err != nil || KeyNum > 0 {
		logger.Lshortfile.Println("warning: can not load DataItemBatchPermissions from redis in param server. Consider Add hash item  DataItemBatchPermissions in redis,with key redis key before ':' and value as permitted batch operations seperated by ','")
		return err
	}
	for k, v := range Permissions_TMP {
		KeyNum++
		var SplittedOptions []string = strings.Split(v, ",")
		for _, oprationPermitted := range SplittedOptions {
			//conver k+"."+oprationPermitted to lower case
			keyLower := k + "." + strings.ToLower(oprationPermitted)
			//set keyhash to xxhash.Sum64String(keyLower)
			PermittedBatchOp[xxhash.Sum64String(keyLower)] = nil
		}
	}
	logger.Lshortfile.Println("loading DataItemBatchPermissions success. num keys:", KeyNum, " PermittedBatchOp size:", len(PermittedBatchOp))
	return nil
}
func RefreshDataItemBatchPermissions() {
	for {
		LoadDataItemBatchPermissionsFromRedis()
		time.Sleep(5 * time.Minute)
	}
}
func IsPermittedBatchOperation(dataKey string, operation string) bool {
	ind := strings.Index(dataKey, ":")
	if ind > 0 {
		dataKey = dataKey[:ind]
	}
	KeyHash := xxhash.Sum64String(dataKey + "." + operation)
	_, ok := PermittedBatchOp[KeyHash]
	return ok
}
