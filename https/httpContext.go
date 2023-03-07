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

	QueryFields         string
	ResponseContentType string
}

var ErrIncompleteRequest = errors.New("incomplete request")

func NewHttpContext(ctx context.Context, r *http.Request, w http.ResponseWriter) (httpCtx *HttpContext, err error) {
	var (
		CmdKeyField  string
		CmdKeyFields []string
	)
	svcContext := &HttpContext{Req: r, Rsb: w, Ctx: ctx}
	if CmdKeyFields = strings.Split(r.URL.Path, "/"); len(CmdKeyFields) < 1 {
		return nil, ErrIncompleteRequest
	}
	//for get
	if CmdKeyField = CmdKeyFields[len(CmdKeyFields)-1]; len(CmdKeyField) < 4 {
		return nil, ErrIncompleteRequest
	}
	if CmdKeyFields = strings.Split(strings.Split(CmdKeyField, "?")[0], "="); len(CmdKeyFields) < 1 {
		return nil, ErrIncompleteRequest
	}
	svcContext.Cmd = CmdKeyFields[0]
	//default key is cmd
	svcContext.Key = svcContext.Cmd
	if len(CmdKeyFields) > 1 {
		svcContext.Key = CmdKeyFields[1]
	}
	if len(CmdKeyFields) == 2 {
		svcContext.Field = CmdKeyFields[2]
	} else {
		svcContext.Field = strings.Join(CmdKeyFields[2:], "=")
	}

	//for response
	svcContext.QueryFields = svcContext.Req.FormValue("Queries")
	//default response content type: application/json
	if svcContext.ResponseContentType = svcContext.Req.FormValue("RspType"); svcContext.ResponseContentType == "" {
		svcContext.ResponseContentType = "application/json"
	}
	return svcContext, nil
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
