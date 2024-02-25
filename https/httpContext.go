package https

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vmihailenco/msgpack/v5"
)

type HttpContext struct {
	Req                 *http.Request
	Rsb                 http.ResponseWriter
	jwtToken            *jwt.Token
	Ctx                 context.Context
	RedisDataSourceName string
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
		param        string
	)
	svcContext := &HttpContext{Req: r, Rsb: w, Ctx: ctx}
	//i.g. https://url.com/rSvc/HGET=UserAvatar=fa4Y3oyQk2swURaJ?Queries=*&RspType=image/jpeg
	if CmdKeyFields = strings.Split(r.URL.Path, "/"); len(CmdKeyFields) < 1 {
		return nil, ErrIncompleteRequest
	}
	if CmdKeyFields = strings.Split(CmdKeyFields[len(CmdKeyFields)-1], "-!"); len(CmdKeyFields) < 2 {
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
		//export enum RspType { json = "-!JSON", jpeg = "-!JPG", ogg = "-!OGG", mpeg = "-!MPEG", mp4 = "-!MP4", none = "", text = "-!TEXT", stream = "-!STREAM" }

		if param = CmdKeyFields[i]; param == "" {
			continue
		} else if ind := strings.Index(param, "="); ind > 0 {
			param = param[:ind+1]
		}
		switch param {
		case "JSON":
			svcContext.ResponseContentType = "application/json"
		case "JPG":
			svcContext.ResponseContentType = "image/jpeg"
		case "OGG":
			svcContext.ResponseContentType = "audio/ogg"
		case "MPEG":
			svcContext.ResponseContentType = "video/mpeg"
		case "MP4":
			svcContext.ResponseContentType = "video/mp4"
		case "TEXT":
			svcContext.ResponseContentType = "text/plain"
		case "STREAM":
			svcContext.ResponseContentType = "application/octet-stream"
		case "DS=": //redis db name RDB=redisDataSourceName
			if svcContext.RedisDataSourceName, err = url.QueryUnescape(CmdKeyFields[i][3:]); err != nil {
				return nil, err
			}
		}
	}
	return svcContext, nil
}

func (svc *HttpContext) MsgpackBodyBytes() (data []byte) {
	var (
		err error
	)
	if svc.Req.ContentLength == 0 {
		return nil
	}
	if !strings.HasPrefix(svc.Req.Header.Get("Content-Type"), "application/octet-stream") {
		return nil
	}
	if data, err = io.ReadAll(svc.Req.Body); err != nil {
		return nil
	}
	return data
}
func (svc *HttpContext) JsonBodyBytes() (data []byte) {
	var (
		err error
	)
	if svc.Req.ContentLength == 0 {
		return nil
	}
	if !strings.HasPrefix(svc.Req.Header.Get("Content-Type"), "application/json") {
		return nil
	}
	if data, err = io.ReadAll(svc.Req.Body); err != nil {
		return nil
	}
	return data
}

// Ensure the body is msgpack format
func (svc *HttpContext) MsgpackBody() (bytes []byte, err error) {
	var (
		data interface{}
	)
	if bytes = svc.MsgpackBodyBytes(); len(bytes) == 0 {
		return nil, fmt.Errorf("empty msgpack body")
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
