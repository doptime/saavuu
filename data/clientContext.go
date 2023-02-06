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
}

func NewContext(ctx context.Context) *Ctx {
	if ctx == nil {
		ctx = context.Background()
	}
	return &Ctx{Ctx: ctx, Rds: config.DataRds}
}

func (dc *Ctx) Get(key string, param interface{}) (err error) {
	return rds.Get(dc.Ctx, dc.Rds, key, param)
}
func (dc *Ctx) Set(key string, param interface{}, expiration time.Duration) (err error) {
	return rds.Set(dc.Ctx, dc.Rds, key, param, expiration)
}

func (dc *Ctx) RPush(key string, param interface{}) (err error) {
	return rds.RPush(dc.Ctx, dc.Rds, key, param)
}
func (dc *Ctx) LSet(key string, index int64, param interface{}) (err error) {
	return rds.LSet(dc.Ctx, dc.Rds, key, index, param)
}
func (dc *Ctx) LGet(key string, index int64, param interface{}) (err error) {
	return rds.LGet(dc.Ctx, dc.Rds, key, index, param)
}
func (dc *Ctx) LLen(key string) (length int64, err error) {
	return rds.LLen(dc.Ctx, dc.Rds, key)
}

// append to Set
func (dc *Ctx) SAdd(key string, param interface{}) (err error) {
	return rds.SAdd(dc.Ctx, dc.Rds, key, param)
}
func (dc *Ctx) SRem(key string, param interface{}) (err error) {
	return rds.SRem(dc.Ctx, dc.Rds, key, param)
}
func (dc *Ctx) SIsMember(key string, param interface{}) (isMember bool, err error) {
	return rds.SIsMember(dc.Ctx, dc.Rds, key, param)
}
func (dc *Ctx) SMembers(key string) (members []string, err error) {
	return rds.SMembers(dc.Ctx, dc.Rds, key)
}

func (dc *Ctx) Time() (tm time.Time, err error) {
	return rds.Time(dc.Ctx, dc.Rds)
}
