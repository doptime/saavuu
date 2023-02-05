package api

import (
	"github.com/yangkequn/saavuu/rds"
)

func (pc *Ctx) HGetAll(key string, mapOut interface{}) (err error) {
	return rds.HGetAll(pc.Ctx, pc.Rds, key, mapOut)
}

func (pc *Ctx) HSet(key string, field string, param interface{}) (err error) {
	return rds.HSet(pc.Ctx, pc.Rds, key, field, param)
}

func (pc *Ctx) HKeys(key string) (fields []string, err error) {
	return rds.HKeys(pc.Ctx, pc.Rds, key)
}
