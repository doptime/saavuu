package api

import (
	"github.com/yangkequn/saavuu/rds"
)

func (pc *ApiCtx) HGetAll(key string, mapOut interface{}) (err error) {
	return rds.HGetAll(pc.Ctx, pc.Rds, key, mapOut)
}

func (pc *ApiCtx) HSet(key string, field string, param interface{}) (err error) {
	return rds.HSet(pc.Ctx, pc.Rds, key, field, param)
}

func (pc *ApiCtx) HKeys(key string) (fields []string, err error) {
	return rds.HKeys(pc.Ctx, pc.Rds, key)
}
