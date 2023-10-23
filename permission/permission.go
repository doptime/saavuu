package permission

import (
	"strings"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/yangkequn/saavuu/config"
)

type Permission struct {
	Key       string
	CreateAt  int64
	WhiteList []string
	BlackList []string
}

func IsPermitted(permitType PermitType, dataKey string, operation string) (ok bool) {
	permitIndex := int(permitType)
	var PermissionMap cmap.ConcurrentMap[string, *Permission] = PermitMaps[permitIndex]
	//for example, if dataKey is "user:1x3", then dataKey will be "user"
	if dataKey = strings.Split(dataKey, ":")[0]; len(dataKey) == 0 {
		return false
	}

	permission, ok := PermissionMap.Get(dataKey)
	//if datakey not in BatchPermission, then create BatchPermission, and add it to BatchPermission in redis
	if !ok {
		permission = &Permission{Key: dataKey, CreateAt: time.Now().Unix(), WhiteList: []string{}, BlackList: []string{}}
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
	if config.Cfg.Api.AutoPermission {
		permission.WhiteList = append(permission.WhiteList, operation)
		//save to redis
		var dataCtx = dataCtx(PermitType(permitType))
		dataCtx.HSet(dataKey, permission)
	} else {
		permission.BlackList = append(permission.BlackList, operation)
		//no changed to redis
	}
	PermissionMap.Set(dataKey, permission)
	return config.Cfg.Api.AutoPermission
}

func permitMapUpdate(newMap map[string]*Permission, oldMap cmap.ConcurrentMap[string, *Permission]) (modified bool) {
	modified = false
	for k, newV := range newMap {
		if oldV, ok := oldMap.Get(k); ok {
			if oldV.CreateAt < newV.CreateAt {
				oldMap.Set(k, newV)
				modified = true
			}
			//check if white list changed
			if len(oldV.WhiteList) != len(newV.WhiteList) {
				oldMap.Set(k, newV)
				modified = true
			} else {
				for i, vi := range oldV.WhiteList {
					if vi != newV.WhiteList[i] {
						oldMap.Set(k, newV)
						modified = true
						break
					}
				}
				for i, vi := range oldV.BlackList {
					if vi != newV.BlackList[i] {
						oldMap.Set(k, newV)
						modified = true
						break
					}
				}
			}
		} else {
			oldMap.Set(k, newV)
			modified = true
		}
	}
	//check if any key in oldMap not in newMap
	for k := range oldMap.Items() {
		if _, ok := newMap[k]; !ok {
			oldMap.Remove(k)
			modified = true
		}
	}
	return modified
}
