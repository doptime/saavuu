package data

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/rds"
)

type Ctx[v any] struct {
	Ctx context.Context
	Rds *redis.Client
	Key string
}

func New[v any](Key string) *Ctx[v] {
	return &Ctx[v]{Ctx: context.Background(), Rds: config.DataRds, Key: Key}
}
func (ctx *Ctx[v]) WithContext(c context.Context) *Ctx[v] {
	return &Ctx[v]{Ctx: c, Rds: ctx.Rds, Key: ctx.Key}
}

func (db *Ctx[v]) Time() (tm time.Time, err error) {
	return rds.Time(db.Ctx, db.Rds)
}

var NonKey = New[interface{}]("")
