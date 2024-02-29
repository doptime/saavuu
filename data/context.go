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

func New[k comparable, v any](ops ...*DataOption) *Ctx[k, v] {
	var (
		rds    *redis.Client
		option *DataOption = &DataOption{}
		ok     bool
	)
	if len(ops) > 0 {
		option = ops[0]
	}
	//panic if Key is empty
	if !specification.GetValidDataKeyName((*v)(nil), &option.Key_) {
		log.Panic().Str("Key is empty in Data.New", option.Key_).Send()
	}
	if rds, ok = config.Rds[option.DataSource_]; !ok {
		log.Info().Str("DataSource not defined in enviroment", option.DataSource_).Send()
		return nil
	}
	ctx := &Ctx[k, v]{Ctx: context.Background(), Rds: rds, Key: option.Key_}
	log.Debug().Str("data New create end!", option.Key_).Send()
	return ctx
}
func (db *Ctx[k, v]) Time() (tm time.Time, err error) {
	cmd := db.Rds.Time(db.Ctx)
	return cmd.Result()
}
