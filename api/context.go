package api

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/yangkequn/saavuu/config"
)

type Ctx struct {
	Ctx         context.Context
	Rds         *redis.Client
	ServiceName string
}

func New(ServiceName string) *Ctx {
	//ensure ServiceKey start with "api:"
	if len(ServiceName) < 4 || ServiceName[:4] != "api:" {
		ServiceName = "api:" + ServiceName
	}

	return &Ctx{Ctx: context.Background(), Rds: config.ParamRds, ServiceName: ServiceName}
}
func (ctx *Ctx) WithContext(c context.Context) *Ctx {
	return &Ctx{Ctx: c, Rds: ctx.Rds, ServiceName: ctx.ServiceName}
}
