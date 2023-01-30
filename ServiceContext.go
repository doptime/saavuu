package saavuu

import (
	"context"

	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/rCtx"
)

func NewParamContext(ctx context.Context) *rCtx.ParamCtx {
	if ctx == nil {
		ctx = context.Background()
	}
	return &rCtx.ParamCtx{Ctx: ctx, Rds: config.ParamRds}
}
func NewDataContext(ctx context.Context) *rCtx.DataCtx {
	if ctx == nil {
		ctx = context.Background()
	}
	return &rCtx.DataCtx{Ctx: ctx, Rds: config.DataRds}
}
