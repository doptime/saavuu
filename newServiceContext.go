package saavuu

import (
	"context"

	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/redisContext"
)

func NewParamContext(ctx context.Context) *redisContext.ParamContext {
	if ctx == nil {
		ctx = context.Background()
	}
	return &redisContext.ParamContext{Ctx: ctx, Rds: config.ParamRds}
}
func NewDataContext(ctx context.Context) *redisContext.DataContext {
	if ctx == nil {
		ctx = context.Background()
	}
	return &redisContext.DataContext{Ctx: ctx, Rds: config.DataRds}
}
