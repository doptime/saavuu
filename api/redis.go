package api

import (
	"context"
	"time"

	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/rds"
)

var RdsOp = Ctx{Rds: config.ParamRds, Ctx: context.Background(), ServiceName: "api:redis"}

func (ac *Ctx) HGetAll(key string, mapOut interface{}) (err error) {
	return rds.HGetAll(ac.Ctx, ac.Rds, key, mapOut)
}

func (ac *Ctx) HSet(key string, field string, param interface{}) (err error) {
	return rds.HSet(ac.Ctx, ac.Rds, key, field, param)
}

func (ac *Ctx) HKeys(key string, fields interface{}) (err error) {
	return rds.HKeys(ac.Ctx, ac.Rds, key, fields)
}

func (ac *Ctx) Time() (tm time.Time, err error) {
	return rds.Time(ac.Ctx, ac.Rds)
}
