package data

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/rds"
)

type Ctx struct {
	Ctx context.Context
	Rds *redis.Client
	Key string
}

func New(Key string) *Ctx {
	return &Ctx{Ctx: context.Background(), Rds: config.DataRds, Key: Key}
}
func (ctx *Ctx) WithContext(c context.Context) *Ctx {
	return &Ctx{Ctx: c, Rds: ctx.Rds, Key: ctx.Key}
}

func (db *Ctx) Time() (tm time.Time, err error) {
	return rds.Time(db.Ctx, db.Rds)
}
