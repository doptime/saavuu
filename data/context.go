package data

import (
	"context"
	"reflect"
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

func NewStruct[v any]() *Ctx[v] {
	var Key string
	_type := reflect.TypeOf((*v)(nil))
	//take name of type v as key
	for _type.Kind() == reflect.Ptr || _type.Kind() == reflect.Slice {
		_type = _type.Elem()
	}
	Key = _type.Name()
	//panic if Key is empty
	if Key == "" {
		panic("Key is empty, please give a key for this data")
	}
	return &Ctx[v]{Ctx: context.Background(), Rds: config.DataRds, Key: Key}
}

func New[v any](Key string) *Ctx[v] {
	//panic if Key is empty
	if Key == "" {
		panic("Key is empty, please give a key for this data")
	}
	return &Ctx[v]{Ctx: context.Background(), Rds: config.DataRds, Key: Key}
}
func (ctx *Ctx[v]) WithContext(c context.Context) *Ctx[v] {
	return &Ctx[v]{Ctx: c, Rds: ctx.Rds, Key: ctx.Key}
}

func (db *Ctx[v]) Time() (tm time.Time, err error) {
	return rds.Time(db.Ctx, db.Rds)
}

var NonKey = NewStruct[interface{}]()
