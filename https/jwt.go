package https

import (
	"fmt"
	. "saavuu/config"

	"github.com/golang-jwt/jwt/v4"
)

func (svc *HttpContext) JwtField(field string) (f interface{}) {
	if svc.jwtToken == nil {
		jwtStr := svc.Req.Header.Get("Authorization")
		if len(jwtStr) == 0 {
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
		if svc.jwtToken == nil {
			return nil
		}

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
