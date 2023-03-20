package https

import (
	"errors"
	"fmt"
	"strings"

	"github.com/yangkequn/saavuu/config"

	"github.com/golang-jwt/jwt/v5"
)

func (svc *HttpContext) ParseJwtToken() (err error) {
	var (
		jwtStr string
	)
	if svc.jwtToken != nil {
		return nil
	}
	if jwtStr = svc.Req.Header.Get("Authorization"); len(jwtStr) == 0 {
		return errors.New("no JWT token")
	}
	//decode jwt string to map[string] interface{} with jwtSrcrets as jwt secret
	keyFunction := func(token *jwt.Token) (value interface{}, err error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(config.Cfg.JwtSecret), nil
	}
	if svc.jwtToken, err = jwt.ParseWithClaims(jwtStr, jwt.MapClaims{}, keyFunction, jwt.WithJSONNumber()); err != nil {
		return fmt.Errorf("invalid JWT token: %v", err)
	}
	return nil
}
func (svc *HttpContext) MergeJwtField(paramIn map[string]interface{}) {
	//remove nay field that starts with "JWT_" in paramIn
	//prevent forged jwt field
	for k := range paramIn {
		if strings.HasPrefix(k, "JWT_") {
			delete(paramIn, k)
		}
	}

	if err := svc.ParseJwtToken(); err != nil {
		return
	}
	//save every field in svc.Jwt.Claims to in
	mpclaims, ok := svc.jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return
	}
	for k, v := range mpclaims {
		if strings.Contains(config.Cfg.JwtFieldsKept, strings.ToLower(k)) {
			paramIn["JWT_"+k] = v
		}
	}
}

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
