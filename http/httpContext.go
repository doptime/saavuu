package http

import (
	"context"
	"net/http"

	. "saavuu/config"

	"github.com/golang-jwt/jwt/v4"
)

type HttpContext struct {
	Req                 *http.Request
	Rsb                 http.ResponseWriter
	Jwt                 *jwt.Token
	Ctx                 context.Context
	Key                 string
	Field               string
	QueryFields         string
	ExpectedReponseType string
}

func NewHttpContext(r *http.Request, w http.ResponseWriter) *HttpContext {
	svcContext := &HttpContext{Req: r, Rsb: w, Ctx: r.Context()}
	svcContext.Jwt, _ = JwtFromHttpRequest(r)
	svcContext.Req.ParseMultipartForm(Cfg.MaxBufferSize)
	svcContext.Key = svcContext.Req.FormValue("Key")
	svcContext.Field = svcContext.Req.FormValue("Field")
	svcContext.QueryFields = svcContext.Req.FormValue("Queries")
	svcContext.QueryFields = svcContext.Req.FormValue("Expect")
	return svcContext
}

func (svc *HttpContext) JwtField(field string) (s string, ok bool) {
	if svc.Jwt == nil {
		return "", false
	}
	mpclaims, ok := svc.Jwt.Claims.(jwt.MapClaims)
	if !ok {
		return "", ok
	}
	id, ok := mpclaims[field].(string)
	return id, ok
}
func (svc *HttpContext) Data() (d string, ok bool) {
	return svc.Req.FormValue("Data"), true
}
