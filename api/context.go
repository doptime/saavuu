package api

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/yangkequn/saavuu/config"
	"golang.org/x/exp/slices"
)

type Ctx[i any, o any] struct {
	Ctx context.Context
	Rds *redis.Client

	Debug       bool
	ServiceName string
	Func        func(InServiceName i) (ret o, err error)
}

var disAllowedServiceNames = []string{"string", "int32", "int64", "float32", "float64", "int", "uint", "float", "bool", "byte", "rune", "complex64", "complex128"}

// create Api context.
// This New function is for the case the API is defined outside of this package.
// If the API is defined in this package, use Api() instead.
func New[i any, o any](ServiceName string) *Ctx[i, o] {
	//remove "api:" prefix
	if len(ServiceName) >= 4 && ServiceName[:4] == "api:" {
		ServiceName = ServiceName[4:]
	}

	//if SerivceName Starts with "In", remove it
	if len(ServiceName) >= 2 && (ServiceName[0:2] == "In" || ServiceName[0:2] == "in") {
		ServiceName = ServiceName[2:]
	}

	//first byte of ServiceName should be lower case
	if ServiceName[0] >= 'A' && ServiceName[0] <= 'Z' {
		ServiceName = string(ServiceName[0]+32) + ServiceName[1:]
	}

	if len(ServiceName) == 0 {
		log.Panic().Msg("Empty ServiceName is empty")
	}
	//panic if servicename is string int32 int64 float32 float64, int, uint, float, bool, byte, rune, complex64, complex128
	if slices.Contains(disAllowedServiceNames, ServiceName) {
		log.Panic().Msg(ServiceName + ":ServiceName misnamed. Check your code")
	}
	//ensure ServiceKey start with "api:"
	ServiceName = "api:" + ServiceName

	return &Ctx[i, o]{Ctx: context.Background(), Rds: config.Rds, Debug: false, ServiceName: ServiceName}
}

// allow setting breakpoint for input decoding
func (ctx *Ctx[i, o]) UseDebug() *Ctx[i, o] {
	ctx.Debug = true
	return ctx
}

// force use new context
func (ctx *Ctx[i, o]) UseContext(c context.Context) *Ctx[i, o] {
	ctx.Ctx = c
	return ctx
}

// force use RPC mode
func (ctx *Ctx[i, o]) UseRPC() *Ctx[i, o] {
	ctx.Func = nil
	return ctx
}
