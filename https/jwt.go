package https

import (
	"fmt"
	"strings"

	. "github.com/yangkequn/saavuu/config"

	"github.com/golang-jwt/jwt/v4"
)

func (svc *HttpContext) JwtField(field string) (f interface{}) {
	var token *jwt.Token = svc.JwtToken()
	if token == nil {
		return nil
	}
	// return field in svc.Jwt.Claims
	mpclaims, ok := svc.jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil
	}
	if f, ok = mpclaims[field]; !ok {
		return nil
	}
	return f
}

func (svc *HttpContext) JwtToken() (token *jwt.Token) {
	var jwtStr string
	if svc.jwtToken == nil {
		if jwtStr = svc.Req.Header.Get("Authorization"); len(jwtStr) == 0 {
			return nil
		}
		//decode jwt string to map[string] interface{} with jwtSrcrets as jwt secret
		keyFunction := func(token *jwt.Token) (value interface{}, err error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("Invalid signing method")
			}
			return []byte(Cfg.JwtSecret), nil
		}
		svc.jwtToken, _ = jwt.ParseWithClaims(jwtStr, jwt.MapClaims{}, keyFunction)
	}
	return svc.jwtToken
}
func (svc *HttpContext) MergeJwtField(paramIn map[string]interface{}) {
	//remove nay field that starts with "jwt_" in paramIn
	//prevent forged jwt field
	for k, _ := range paramIn {
		if strings.HasPrefix(k, "jwt_") {
			delete(paramIn, k)
		}
	}

	var token *jwt.Token = svc.JwtToken()
	if token == nil {
		return
	}
	//save every field in svc.Jwt.Claims to in
	mpclaims, ok := svc.jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return
	}
	for k, v := range mpclaims {
		if !strings.Contains(Cfg.JwtIgnoreFields, strings.ToLower(k)) {
			paramIn["jwt_"+k] = v
		}
	}
}
