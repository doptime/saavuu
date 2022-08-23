package saavuu

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

func JwtFromHttpRequest(r *http.Request) (t *jwt.Token, err error) {
	jwtStr := r.Header.Get("Authorization")
	if jwtStr == "" {
		return nil, fmt.Errorf("no Authorization header")
	}
	//decode jwt string to map[string] interface{} with jwtSrcrets as jwt secret
	//map[string] interface{} is the type of jwt.Claims
	keyFunction := func(t *jwt.Token) (value interface{}, err error) {
		return []byte(Config.JwtToken), nil
	}
	return jwt.ParseWithClaims(jwtStr, &jwt.MapClaims{}, keyFunction)
}
