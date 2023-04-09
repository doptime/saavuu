package api

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/yangkequn/saavuu/config"
)

type Ctx[v any] struct {
	Ctx         context.Context
	Rds         *redis.Client
	Debug       bool
	ServiceName string
}

func New[v any](ServiceName string) *Ctx[v] {
	//ensure ServiceKey start with "api:"
	if len(ServiceName) < 4 || ServiceName[:4] != "api:" {
		ServiceName = "api:" + ServiceName
	}

	return &Ctx[v]{Ctx: context.Background(), Rds: config.ParamRds, Debug: false, ServiceName: ServiceName}
}
func (ctx *Ctx[v]) WithDebug() *Ctx[v] {
	return &Ctx[v]{Ctx: ctx.Ctx, Rds: ctx.Rds, Debug: true, ServiceName: ctx.ServiceName}
}
func (ctx *Ctx[v]) WithContext(c context.Context) *Ctx[v] {
	return &Ctx[v]{Ctx: c, Rds: ctx.Rds, ServiceName: ctx.ServiceName}
}
