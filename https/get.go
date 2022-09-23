package https

import (
	"errors"
	"strings"

	. "github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/redisContext"

	"github.com/vmihailenco/msgpack/v5"
)

func replaceAtUseJwt(scvCtx *HttpContext, jwts map[string]interface{}, s string) (r string, err error) {
	var (
		jwtValue    interface{}
		ok          bool
		jwtFieldStr string
	)
	for i := strings.Index(s, "@"); i >= 0; i = strings.Index(s, "@") {
		ie := i + 1
		for ; ie < len(s) && ((s[ie] >= '0' && s[ie] <= '9') || (s[ie] >= 'a' && s[ie] <= 'z') || (s[ie] >= 'A' && s[ie] <= 'Z')); ie++ {
		}
		if ie == i+1 {
			return "", errors.New("invalid field")
		}
		label := s[i+1 : ie]
		if jwtValue, ok = jwts["JWT_"+label]; !ok {
			return "", errors.New("invalid jwt field")
		}
		if jwtFieldStr, ok = jwtValue.(string); !ok {
			return "", errors.New("invalid jwt field")
		}
		s = s[:i] + jwtFieldStr + s[ie:]
	}
	return s, nil
}
func (svcCtx *HttpContext) GetHandler() (ret interface{}, err error) {
	var (
		jwts         map[string]interface{} = map[string]interface{}{}
		data         []byte
		resultBytes  []byte                 = []byte{}
		resultString string                 = ""
		result       map[string]interface{} = map[string]interface{}{}
	)

	svcCtx.MergeJwtField(jwts)

	if svcCtx.Key, err = replaceAtUseJwt(svcCtx, jwts, svcCtx.Key); err != nil {
		return nil, err
	} else if len(svcCtx.Key) == 0 {
		return nil, errors.New("no key")
	}

	if svcCtx.Field, err = replaceAtUseJwt(svcCtx, jwts, svcCtx.Field); err != nil {
		return nil, err
	}

	//check auth. only Key start with upper case are allowed to access
	if len(svcCtx.Key) <= 0 || !(svcCtx.Key[0] >= 'A' && svcCtx.Key[0] <= 'Z') {
		return nil, errors.New("no auth")
	}
	//case Is a member of a set
	if strings.Index(svcCtx.Key, "SMember:") == 0 {
		dc := redisContext.DataCtx{Ctx: svcCtx.Ctx, Rds: DataRds}
		if ok := dc.SIsMember(svcCtx.Key[8:], svcCtx.Field); ok {
			return "{member:true}", nil
		}
		return "{member:false}", nil

	}
	if strings.Index(svcCtx.Key, "HEXISTS:") == 0 {
		dc := redisContext.DataCtx{Ctx: svcCtx.Ctx, Rds: DataRds}
		if ok := dc.HExists(svcCtx.Key[8:], svcCtx.Field); ok {
			return "{member:true}", nil
		}
		return "{member:false}", nil
	}

	cmd := DataRds.HGet(svcCtx.Ctx, svcCtx.Key, svcCtx.Field)
	if data, err = cmd.Bytes(); err != nil {
		return "", nil
	}
	//fill content type, to support binary or json response
	if svcCtx.ResponseContentType != "application/json" {
		if msgpack.Unmarshal(data, &resultBytes) == nil {
			return resultBytes, err
		}
		if msgpack.Unmarshal(data, &resultString) == nil {
			return resultString, err
		}
	}
	if err = msgpack.Unmarshal(data, &result); err != nil {
		return nil, errors.New("unsupported data type")
	}
	return result, nil
}
