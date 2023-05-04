package permission

import (
	"context"
	"strings"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/data"
)

func IsPermitted(PermissionMap cmap.ConcurrentMap[string, Permission], PermissionKey *string, dataKey string, operation string) (ok bool) {
	dataKey = strings.Split(dataKey, ":")[0]
	// only care non-digit part of dataKey
	//split dataKey with number digit char, and get the first part
	//for example, if dataKey is "user1x3", then dataKey will be "user"
	for i, v := range dataKey {
		if v >= '0' && v <= '9' {
			dataKey = dataKey[:i]
			break
		}
	}
	if len(dataKey) == 0 {
		return false
	}

	permission, ok := PermissionMap.Get(dataKey)
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
	if config.Cfg.AutoPermission {
		permission.WhiteList = append(permission.WhiteList, operation)
	} else {
		permission.BlackList = append(permission.BlackList, operation)
	}
	PermissionMap.Set(dataKey, permission)
	//save to redis
	var paramRds = data.Ctx[Permission]{Rds: config.Rds, Ctx: context.Background(), Key: *PermissionKey}
	paramRds.HSet(dataKey, permission)
	return config.Cfg.AutoPermission
}
