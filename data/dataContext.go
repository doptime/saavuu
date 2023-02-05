package data

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
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
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := dc.Rds.RPush(dc.Ctx, key, bytes)
	return status.Err()
}
func (dc *Ctx) LSet(key string, index int64, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := dc.Rds.LSet(dc.Ctx, key, index, bytes)
	return status.Err()
}
func (dc *Ctx) LGet(key string, index int64, param interface{}) (err error) {
	cmd := dc.Rds.LIndex(dc.Ctx, key, index)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, param)
}
func (dc *Ctx) LLen(key string) (length int64) {
	cmd := dc.Rds.LLen(dc.Ctx, key)
	return cmd.Val()
}

// append to Set
func (dc *Ctx) SAdd(key string, members ...interface{}) (err error) {
	status := dc.Rds.SAdd(dc.Ctx, key, members)
	return status.Err()
}
func (dc *Ctx) SRem(key string, members ...interface{}) (err error) {
	status := dc.Rds.SRem(dc.Ctx, key, members)
	return status.Err()
}
func (dc *Ctx) SIsMember(key string, param interface{}) (ok bool, err error) {
	cmd := dc.Rds.SIsMember(dc.Ctx, key, param)
	return cmd.Result()
}
func (dc *Ctx) Time() (tm time.Time, err error) {
	cmd := dc.Rds.Time(dc.Ctx)
	return cmd.Result()
}
func (dc *Ctx) SMembers(key string, param interface{}) (members []string, err error) {
	cmd := dc.Rds.SMembers(dc.Ctx, key)
	members, err = cmd.Result()
	if err != nil {
		return nil, err
	}
	return members, nil
}
