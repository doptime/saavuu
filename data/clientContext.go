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

func NewCtx(ctx context.Context, Key string) *Ctx {
	if ctx == nil {
		ctx = context.Background()
	}
	return &Ctx{Ctx: ctx, Rds: config.DataRds, Key: Key}
}

func (db *Ctx) Time() (tm time.Time, err error) {
	return rds.Time(db.Ctx, db.Rds)
}
