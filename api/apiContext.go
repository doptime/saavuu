package api

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/yangkequn/saavuu/config"
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
