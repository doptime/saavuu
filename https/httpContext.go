package https

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/vmihailenco/msgpack/v5"
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

	ResponseFields      string
	ResponseContentType string
}

var ErrIncompleteRequest = errors.New("incomplete request")

func NewHttpContext(ctx context.Context, r *http.Request, w http.ResponseWriter) (httpCtx *HttpContext, err error) {
	var (
		CmdKeyFields []string
	)
	svcContext := &HttpContext{Req: r, Rsb: w, Ctx: ctx}
	//i.g. https://url.com/rSvc/HGET=UserAvatar=fa4Y3oyQk2swURaJ?Queries=*&RspType=image/jpeg
	if CmdKeyFields = strings.Split(r.URL.Path, "/"); len(CmdKeyFields) < 1 {
		return nil, ErrIncompleteRequest
	}
	//this last part of url is cmd and key and field, i.g. /HGET?K=UserAvatar&F=fa4Y3oyQk2swURaJ
	//read first param as cmd
	svcContext.Cmd = CmdKeyFields[len(CmdKeyFields)-1]
	svcContext.Key = r.FormValue("K")
	svcContext.Field = r.FormValue("F")
	//for response
	if svcContext.ResponseFields = svcContext.Req.FormValue("QF"); len(svcContext.ResponseFields) == 0 {
		svcContext.ResponseFields = "*"
	} else if svcContext.ResponseFields == "null" {
		svcContext.ResponseFields = ""
	}

	//default response content type: application/json
	if svcContext.ResponseContentType = svcContext.Req.FormValue("RspType"); svcContext.ResponseContentType == "" {
		svcContext.ResponseContentType = "application/json"
	}
	return svcContext, nil
}
func (svc *HttpContext) SetContentType() {
	if len(svc.ResponseContentType) > 0 && len(svc.ResponseFields) > 0 {
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
	if data, err = io.ReadAll(svc.Req.Body); err != nil {
		return nil, err
	}
	return data, nil
}

// Ensure the body is msgpack format
func (svc *HttpContext) MsgpackBody() (bytes []byte, err error) {
	var (
		data interface{}
	)
	if bytes, err = svc.BodyBytes(); err != nil {
		return nil, err
	}
	//should make sure the data is msgpack format
	if err = msgpack.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}
	if bytes, err = msgpack.Marshal(data); err != nil {
		return nil, err
	}
	//return remarshaled bytes, because golang msgpack is better fullfill than javascript msgpack
	return bytes, nil
}
