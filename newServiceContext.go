package saavuu

import (
	"context"

	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/redisContext"
)

func NewRedisContext() *redisContext.RedisContext {
	var ctx context.Context = context.Background()
	return &redisContext.RedisContext{Ctx: ctx, ParamRds: config.ParamRds, DataRds: config.DataRds}
}
