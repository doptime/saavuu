package api

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/rds"
)

type Ctx struct {
	Ctx context.Context
	Rds *redis.Client
}

func NewContext(ctx context.Context) *Ctx {
	if ctx == nil {
		ctx = context.Background()
	}
	return &Ctx{Ctx: ctx, Rds: config.ParamRds}
}

func (pc *Ctx) HGetAll(key string, mapOut interface{}) (err error) {
	return rds.HGetAll(pc.Ctx, pc.Rds, key, mapOut)
}

func (pc *Ctx) HSet(key string, field string, param interface{}) (err error) {
	return rds.HSet(pc.Ctx, pc.Rds, key, field, param)
}

func (pc *Ctx) HKeys(key string) (fields []string, err error) {
	return rds.HKeys(pc.Ctx, pc.Rds, key)
}
