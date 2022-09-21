package https

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/vmihailenco/msgpack/v5"
)

type HttpContext struct {
	Req      *http.Request
	Rsb      http.ResponseWriter
	jwtToken *jwt.Token
	Ctx      context.Context
	// case get
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
	svcContext.Key = svcContext.Req.FormValue("Key")
	svcContext.Field = svcContext.Req.FormValue("Field")
	//for post
	svcContext.Service = svcContext.Req.FormValue("Service")

	//for response
	svcContext.QueryFields = svcContext.Req.FormValue("Queries")
	svcContext.ResponseContentType = svcContext.Req.FormValue("RspType")
	return svcContext
}

func (svc *HttpContext) BodyMessage() (msgPack map[string]interface{}, err error) {
	var (
		param map[string]interface{} = map[string]interface{}{}
		data  []byte                 = make([]byte, svc.Req.ContentLength)
		ctype string                 = svc.Req.Header.Get("Content-Type")
	)
	//read request body to data
	if ctype == "application/json" {

		if data, err = ioutil.ReadAll(svc.Req.Body); err != nil {
			return nil, err
		}
		//unpack with json and convert to msgpack
		if err = json.Unmarshal(data, &param); err != nil {
			return nil, err
		}
		return param, nil
	} else if ctype == "octet-stream" || ctype == "application/octet-stream" {
		if data, err = ioutil.ReadAll(svc.Req.Body); err != nil {
			return nil, err
		}
		//unpack with msgpack
		if err = msgpack.Unmarshal(data, &param); err != nil {
			return nil, err
		}
		return param, nil
	}
	return nil, errors.New("unsupported content type")
}
