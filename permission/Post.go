package permission

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var PermittedPostOp map[string]Permission = make(map[string]Permission)

func LoaIsPermittedPostField(operation string, Field *string, token *jwt.Token) (operationNew string, err error) {
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
			return operationNew, fmt.Errorf("jwt missing subTag " + subTag)
		}
		*Field = strings.Join(FieldParts, "")
	}
	return operationNew, nil
}
func IsPostPermitted(dataKey string, operation string) (ok bool) {
	return IsPermitted(PermittedPostOp, &permitKeyPost, dataKey, operation)
}
