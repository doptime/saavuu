package https

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/yangkequn/saavuu/config"
)

func ConvertMapToJwtString(param map[string]interface{}) (jwtString string, err error) {
	//convert map to jwt.claims
	claims := jwt.MapClaims{}
	for k, v := range param {
		claims[k] = v
	}
	//create jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//sign jwt token
	jwtString, err = token.SignedString([]byte(config.Cfg.JwtSecret))
	return jwtString, err
}
