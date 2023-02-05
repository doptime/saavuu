package saavuu

import (
	"context"

	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/rCtx"
)

func NewDataContext(ctx context.Context) *rCtx.DataCtx {
	if ctx == nil {
		ctx = context.Background()
	}
	return &rCtx.DataCtx{Ctx: ctx, Rds: config.DataRds}
}
