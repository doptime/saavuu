package api

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/yangkequn/saavuu/config"
)

type Ctx[i any, o any] struct {
	Ctx         context.Context
	Rds         *redis.Client
	Debug       bool
	ServiceName string
	DoLocal     func(InServiceName i) (ret o, err error)
}

// create Api context.
// This New function is for the case the API is defined outside of this package.
// If the API is defined in this package, use Api() instead.
func New[i any, o any](ServiceName string) *Ctx[i, o] {
	//ensure ServiceKey start with "api:"
	if len(ServiceName) < 4 || ServiceName[:4] != "api:" {
		ServiceName = "api:" + ServiceName
	}

	return &Ctx[i, o]{Ctx: context.Background(), Rds: config.ParamRds, Debug: false, ServiceName: ServiceName}
}
func (ctx *Ctx[i, o]) WithDebug() *Ctx[i, o] {
	return &Ctx[i, o]{Ctx: ctx.Ctx, Rds: ctx.Rds, Debug: true, ServiceName: ctx.ServiceName}
}
func (ctx *Ctx[i, o]) WithContext(c context.Context) *Ctx[i, o] {
	return &Ctx[i, o]{Ctx: c, Rds: ctx.Rds, ServiceName: ctx.ServiceName}
}
