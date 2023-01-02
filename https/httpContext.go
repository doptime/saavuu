package https

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type HttpContext struct {
	Req      *http.Request
	Rsb      http.ResponseWriter
	jwtToken *jwt.Token
	Ctx      context.Context
	// case get
	Cmd   string
	Key   string
	Field string
	// case post
	Service string

	QueryFields         string
	ResponseContentType string
}

func NewHttpContext(ctx context.Context, r *http.Request, w http.ResponseWriter) *HttpContext {
	svcContext := &HttpContext{Req: r, Rsb: w, Ctx: ctx}
	//for get
	svcContext.Cmd = svcContext.Req.FormValue("Cmd")
	svcContext.Key = svcContext.Req.FormValue("Key")
	svcContext.Field = svcContext.Req.FormValue("Field")

	//for post
	svcContext.Service = svcContext.Req.FormValue("Service")

	//for response
	svcContext.QueryFields = svcContext.Req.FormValue("Queries")
	svcContext.ResponseContentType = svcContext.Req.FormValue("RspType")
	return svcContext
}
func (svc *HttpContext) SetContentType() {
	if len(svc.ResponseContentType) > 0 && len(svc.QueryFields) > 0 {
		svc.Rsb.Header().Set("Content-Type", svc.ResponseContentType)
	}
}

func (svc *HttpContext) BodyMessage() (param map[string]interface{}, err error) {
	var data []byte = nil
	if data, err = svc.BodyBytes(); err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("empty body")
	}
	return map[string]interface{}{"MsgPack": data}, nil
}

func (svc *HttpContext) BodyBytes() (data []byte, err error) {
	var (
		ctype string = svc.Req.Header.Get("Content-Type")
	)
	if ctype != "application/octet-stream" {
		return nil, errors.New("unsupported content type")
	}
	if svc.Req.ContentLength == 0 {
		return nil, errors.New("empty body")
	}
	if data, err = ioutil.ReadAll(svc.Req.Body); err != nil {
		return nil, err
	}
	return data, nil
}
