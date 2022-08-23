package saavuu

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

type ServiceContext struct {
	req         *http.Request
	rsb         http.ResponseWriter
	Jwt         *jwt.Token
	ctx         context.Context
	Key         string
	Field       string
	QueryFields string
}

func LoadHttpContext(r *http.Request, w http.ResponseWriter) *ServiceContext {
	svcContext := &ServiceContext{req: r, rsb: w, ctx: r.Context()}
	svcContext.Jwt, _ = JwtFromHttpRequest(r)
	svcContext.req.ParseMultipartForm(Config.MaxBufferSize)
	svcContext.Key = svcContext.req.FormValue("Key")
	svcContext.Field = svcContext.req.FormValue("Field")
	svcContext.QueryFields = svcContext.req.FormValue("Queries")
	return svcContext
}

func (svc *ServiceContext) JwtField(field string) (s string, ok bool) {
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
func (svc *ServiceContext) Data() (d string, ok bool) {
	return svc.req.FormValue("Data"), true
}
