package rCtx

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v5"
)

type DataCtx struct {
	Ctx context.Context
	Rds *redis.Client
}

func (dc *DataCtx) Get(key string, param interface{}) (err error) {
	cmd := dc.Rds.Get(dc.Ctx, key)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, param)
}
func (dc *DataCtx) Set(key string, param interface{}, expiration time.Duration) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := dc.Rds.Set(dc.Ctx, key, bytes, expiration)
	return status.Err()
}

func (dc *DataCtx) RPush(key string, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := dc.Rds.RPush(dc.Ctx, key, bytes)
	return status.Err()
}
func (dc *DataCtx) LSet(key string, index int64, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := dc.Rds.LSet(dc.Ctx, key, index, bytes)
	return status.Err()
}
func (dc *DataCtx) LGet(key string, index int64, param interface{}) (err error) {
	cmd := dc.Rds.LIndex(dc.Ctx, key, index)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, param)
}
func (dc *DataCtx) LLen(key string) (length int64) {
	cmd := dc.Rds.LLen(dc.Ctx, key)
	return cmd.Val()
}

// append to Set
func (dc *DataCtx) SAdd(key string, members ...interface{}) (err error) {
	status := dc.Rds.SAdd(dc.Ctx, key, members)
	return status.Err()
}
func (dc *DataCtx) SRem(key string, members ...interface{}) (err error) {
	status := dc.Rds.SRem(dc.Ctx, key, members)
	return status.Err()
}
func (dc *DataCtx) SIsMember(key string, param interface{}) (ok bool) {
	cmd := dc.Rds.SIsMember(dc.Ctx, key, param)
	return cmd.Val()
}
func (dc *DataCtx) SMembers(key string, param interface{}) (members []string, err error) {
	cmd := dc.Rds.SMembers(dc.Ctx, key)
	members, err = cmd.Result()
	if err != nil {
		return nil, err
	}
	return members, nil
}
