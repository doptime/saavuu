package https

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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
	if CmdKeyFields = strings.Split(CmdKeyFields[len(CmdKeyFields)-1], "--"); len(CmdKeyFields) < 2 {
		return nil, ErrIncompleteRequest
	}
	// cmd and key and field, i.g. /HGET/UserAvatar?F=fa4Y3oyQk2swURaJ
	svcContext.Cmd, svcContext.Key = CmdKeyFields[0], CmdKeyFields[1]
	//url decoded already
	svcContext.Field = r.FormValue("F")

	//default response content type: application/json
	svcContext.ResponseContentType = "application/json"
	for i, l := 2, len(CmdKeyFields); i < l; i++ {
		//export enum RspType { json = "&RspType=application/json", jpeg = "&RspType=image/jpeg", ogg = "&RspType=audio/ogg", mpeg = "&RspType=video/mpeg", mp4 = "&RspType=video/mp4", none = "", text = "&RspType=text/plain", stream = "&RspType=application/octet-stream" }
		//export enum RspType { json = "--!JSON", jpeg = "--!JPG", ogg = "--!OGG", mpeg = "--!MPEG", mp4 = "--!MP4", none = "", text = "--!TEXT", stream = "--!STREAM" }
		var param string = CmdKeyFields[i]
		if param == "" || param[0] != '!' {
			continue
		}
		switch param {
		case "!JSON":
			svcContext.ResponseContentType = "application/json"
		case "!JPG":
			svcContext.ResponseContentType = "image/jpeg"
		case "!OGG":
			svcContext.ResponseContentType = "audio/ogg"
		case "!MPEG":
			svcContext.ResponseContentType = "video/mpeg"
		case "!MP4":
			svcContext.ResponseContentType = "video/mp4"
		case "!TEXT":
			svcContext.ResponseContentType = "text/plain"
		case "!STREAM":
			svcContext.ResponseContentType = "application/octet-stream"
		}
	}
	return svcContext, nil
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
