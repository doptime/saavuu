package saavuu

import (
	"context"

	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/rCtx"
)

func NewApiContext(ctx context.Context) *rCtx.ApiCtx {
	if ctx == nil {
		ctx = context.Background()
	}
	return &rCtx.ApiCtx{Ctx: ctx, Rds: config.ParamRds}
}
