package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/vmihailenco/msgpack/v5"
)

type HttpContext struct {
	Req *http.Request
	Rsb http.ResponseWriter
	Jwt *jwt.Token
	Ctx context.Context
	// case get
	Key   string
	Field string
	// case post
	Service string

	QueryFields         string
	ExpectedReponseType string
}

func NewHttpContext(r *http.Request, w http.ResponseWriter) *HttpContext {
	svcContext := &HttpContext{Req: r, Rsb: w, Ctx: r.Context()}
	svcContext.Jwt, _ = JwtFromHttpRequest(r)
	svcContext.Key = svcContext.Req.FormValue("Key")
	svcContext.Field = svcContext.Req.FormValue("Field")
	svcContext.Service = svcContext.Req.FormValue("Service")

	svcContext.QueryFields = svcContext.Req.FormValue("Queries")
	svcContext.QueryFields = svcContext.Req.FormValue("Expect")
	return svcContext
}

func (svc *HttpContext) JwtField(field string) (f interface{}) {
	if svc.Jwt == nil {
		return nil
	}
	mpclaims, ok := svc.Jwt.Claims.(jwt.MapClaims)
	if !ok {
		return nil
	}
	if f, ok = mpclaims[field]; !ok {
		return nil
	}
	return f
}
func (svc *HttpContext) BodyMessage() (msgPack map[string]interface{}, err error) {
	var (
		param map[string]interface{} = map[string]interface{}{}
		data  []byte                 = make([]byte, svc.Req.ContentLength)
		ctype string                 = svc.Req.Header.Get("Content-Type")
		body  io.ReadCloser
	)
	//read request body to data
	if ctype == "application/json" {

		if body, err = svc.Req.GetBody(); err != nil {
			return nil, err
		}
		body.Read(data)
		body.Close()
		//unpack with json and convert to msgpack
		if err = json.Unmarshal(data, &param); err != nil {
			return nil, err
		}
		return param, nil
	} else if ctype == "octet-stream" {
		if body, err = svc.Req.GetBody(); err != nil {
			return nil, err
		}
		body.Read(data)
		body.Close()
		//unpack with msgpack
		if err = msgpack.Unmarshal(data, &param); err != nil {
			return nil, err
		}
		return param, nil
	}
	return nil, errors.New("unsupported content type")
}
