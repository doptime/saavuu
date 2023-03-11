package permission

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func IsPermittedField(operation string, Field *string, token *jwt.Token) (operationNew string, err error) {
	var (
		mpclaims jwt.MapClaims
		ok       bool
		obj      interface{}
		subTag   string
		f64      float64
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
		if subTag = FieldParts[len(FieldParts)-1]; len(subTag) == 0 {
			return operationNew, fmt.Errorf("jwt missing subTag " + subTag)
		}
		if obj, ok = mpclaims[subTag]; !ok {
			return operationNew, fmt.Errorf("jwt missing subTag " + subTag)
		}
		// if 64 is int, convert to int
		if f64, ok = obj.(float64); ok && f64 == float64(int64(f64)) {
			obj = int64(f64)
		}
		FieldParts[len(FieldParts)-1] = fmt.Sprintf("%v", obj)
		*Field = strings.Join(FieldParts, "")
	}
	return operationNew, nil
}
