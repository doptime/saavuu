package saavuu

import (
	"context"

	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/redisContext"
)

func NewParamContext(ctx context.Context) *redisContext.ParamCtx {
	if ctx == nil {
		ctx = context.Background()
	}
	return &redisContext.ParamCtx{Ctx: ctx, Rds: config.ParamRds}
}
func NewDataContext(ctx context.Context) *redisContext.DataCtx {
	if ctx == nil {
		ctx = context.Background()
	}
	return &redisContext.DataCtx{Ctx: ctx, Rds: config.DataRds}
}
