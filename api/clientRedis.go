package api

import (
	"github.com/yangkequn/saavuu/rds"
)

func (ac *Ctx) HGetAll(key string, mapOut interface{}) (err error) {
	return rds.HGetAll(ac.Ctx, ac.Rds, key, mapOut)
}

func (ac *Ctx) HSet(key string, field string, param interface{}) (err error) {
	return rds.HSet(ac.Ctx, ac.Rds, key, field, param)
}

func (ac *Ctx) HKeys(key string) (fields []string, err error) {
	return rds.HKeys(ac.Ctx, ac.Rds, key)
}
