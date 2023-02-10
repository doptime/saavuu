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

func (db *Ctx) Get(key string, param interface{}) (err error) {
	return rds.Get(db.Ctx, db.Rds, key, param)
}
func (db *Ctx) Set(key string, param interface{}, expiration time.Duration) (err error) {
	return rds.Set(db.Ctx, db.Rds, key, param, expiration)
}

func (db *Ctx) RPush(key string, param interface{}) (err error) {
	return rds.RPush(db.Ctx, db.Rds, key, param)
}
func (db *Ctx) LSet(key string, index int64, param interface{}) (err error) {
	return rds.LSet(db.Ctx, db.Rds, key, index, param)
}
func (db *Ctx) LGet(key string, index int64, param interface{}) (err error) {
	return rds.LGet(db.Ctx, db.Rds, key, index, param)
}
func (db *Ctx) LLen(key string) (length int64, err error) {
	return rds.LLen(db.Ctx, db.Rds, key)
}

// append to Set
func (db *Ctx) SAdd(key string, param interface{}) (err error) {
	return rds.SAdd(db.Ctx, db.Rds, key, param)
}
func (db *Ctx) SRem(key string, param interface{}) (err error) {
	return rds.SRem(db.Ctx, db.Rds, key, param)
}
func (db *Ctx) SIsMember(key string, param interface{}) (isMember bool, err error) {
	return rds.SIsMember(db.Ctx, db.Rds, key, param)
}
func (db *Ctx) SMembers(key string) (members []string, err error) {
	return rds.SMembers(db.Ctx, db.Rds, key)
}

func (db *Ctx) Time() (tm time.Time, err error) {
	return rds.Time(db.Ctx, db.Rds)
}
