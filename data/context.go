package data

import (
	"context"
	"reflect"
	"time"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/yangkequn/saavuu/config"
)

type Ctx[k comparable, v any] struct {
	Ctx             context.Context
	Rds             *redis.Client
	Key             string
	BloomFilterKeys *bloom.BloomFilter
}

func NewStruct[k comparable, v any]() *Ctx[k, v] {
	var (
		rds *redis.Client = config.RdsDefaultClient()
		Key string
	)
	_type := reflect.TypeOf((*v)(nil))
	//take name of type v as key
	for _type.Kind() == reflect.Ptr || _type.Kind() == reflect.Slice {
		_type = _type.Elem()
	}
	Key = _type.Name()

	log.Debug().Str("data NewStruct create start!", Key).Send()

	//panic if Key is empty
	if Key == "" {
		panic("Key is empty, please give a key for this data")
	}
	ctx := &Ctx[k, v]{Ctx: context.Background(), Rds: rds, Key: Key}

	log.Debug().Str("data NewStruct create end!", Key).Send()

	return ctx
}

func New[k comparable, v any](Key string) *Ctx[k, v] {
	var rds *redis.Client = config.RdsDefaultClient()
	log.Debug().Str("data New create start!", Key).Send()
	//panic if Key is empty
	if Key == "" {
		panic("Key is empty, please give a key for this data")
	}
	ctx := &Ctx[k, v]{Ctx: context.Background(), Rds: rds, Key: Key}
	log.Debug().Str("data New create end!", Key).Send()
	return ctx
}
func (ctx *Ctx[k, v]) WithContext(c context.Context) *Ctx[k, v] {
	return &Ctx[k, v]{Ctx: c, Rds: ctx.Rds, Key: ctx.Key}
}
func (ctx *Ctx[k, v]) WithRedis(RedisServerName string) *Ctx[k, v] {
	var (
		rds *redis.Client
		err error
	)
	if rds, err = config.RdsClientByName(RedisServerName); err != nil {
		log.Panic().AnErr("WithRedis error", err).Send()
		panic(err)
	}
	return &Ctx[k, v]{Ctx: ctx.Ctx, Rds: rds, Key: ctx.Key}
}

func (db *Ctx[k, v]) Time() (tm time.Time, err error) {
	cmd := db.Rds.Time(db.Ctx)
	return cmd.Result()
}
