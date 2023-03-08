package permission

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/yangkequn/saavuu/api"
	"github.com/yangkequn/saavuu/config"
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
	var _map map[string]Permission = make(map[string]Permission)
	// read RedisDelPermission usiing ParamRds
	// RedisDelPermission is a hash
	// split each value of RedisDelPermission into string[] and store in PermittedDelOp
	if err := api.RdsOp.HGetAll("RedisDelPermission", _map); err != nil {
		logger.Lshortfile.Println("loading RedisDelPermission  error: " + err.Error())
	} else {
		if info := fmt.Sprint("loading RedisDelPermission success. num keys:", len(_map)); info != lastLoadDelPermissionInfo {
			logger.Lshortfile.Println(info)
			lastLoadDelPermissionInfo = info
		}
		PermittedDelOp = _map
	}
	time.Sleep(time.Second * 10)
	go LoadPutPermissionFromRedis()
}
func IsPermittedDelField(operation string, Field *string, token *jwt.Token) (operationNew string, err error) {
	var (
		mpclaims jwt.MapClaims
		ok       bool
	)
	// Field contains @*, replace @* with jwt value
	// 只要设置的时候，有@id,@pub，可以确保写不越权，因为 是"@" + operation
	if len(*Field) > 0 {
		operationNew = "@" + operation
		FieldParts := strings.Split(*Field, "@")
		if token == nil || token.Claims == nil {
			return operationNew, fmt.Errorf("JWT token is nil")
		}
		if mpclaims, ok = token.Claims.(jwt.MapClaims); !ok {
			return operationNew, fmt.Errorf("JWT token is invalid")
		}
		subTag := FieldParts[len(FieldParts)-1]
		if FieldParts[len(FieldParts)-1], ok = mpclaims[subTag].(string); !ok {
			return operationNew, fmt.Errorf("jwt missiong subTag " + subTag)
		}
		*Field = strings.Join(FieldParts, "")
	}
	return operationNew, nil
}
func IsDelPermitted(dataKey string, operation string) bool {
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

	permission, ok := PermittedDelOp[dataKey]
	//if datakey not in BatchPermission, then create BatchPermission, and add it to BatchPermission in redis
	if !ok {
		//auto set default permission, the operation is dis-allowed
		permission = Permission{Key: dataKey, CreateAt: time.Now().Unix(), WhiteList: []string{}, BlackList: []string{}}
	}
	//caes ok

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
	if config.Cfg.DevelopMode {
		permission.WhiteList = append(permission.WhiteList, operation)
	} else {
		permission.BlackList = append(permission.BlackList, operation)
	}
	PermittedDelOp[dataKey] = permission
	//save to redis
	api.RdsOp.HSet("RedisDelPermission", dataKey, permission)
	return config.Cfg.DevelopMode
}
