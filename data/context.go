package data

import (
	"context"
	"time"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/specification"
)

type Ctx[k comparable, v any] struct {
	Ctx             context.Context
	Rds             *redis.Client
	Key             string
	BloomFilterKeys *bloom.BloomFilter
}

func New[k comparable, v any](ops ...With) *Ctx[k, v] {
	var (
		rds    *redis.Client
		option *Options = mergeOptions(ops...)
		ok     bool
	)
	//panic if Key is empty
	if !specification.GetValidDataKeyName((*v)(nil), &option.Key) {
		log.Panic().Str("Key is empty in Data.New", option.Key).Send()
	}
	if rds, ok = config.Rds[option.DataSourceName]; !ok {
		log.Info().Str("DataSourceName not defined in enviroment", option.DataSourceName).Send()
		return nil
	}
	ctx := &Ctx[k, v]{Ctx: context.Background(), Rds: rds, Key: option.Key}
	log.Debug().Str("data New create end!", option.Key).Send()
	return ctx
}
func (db *Ctx[k, v]) Time() (tm time.Time, err error) {
	cmd := db.Rds.Time(db.Ctx)
	return cmd.Result()
}
