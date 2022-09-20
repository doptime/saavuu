package redisContext

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v5"
)

type RedisContext struct {
	Ctx       context.Context
	RdsClient *redis.Client
}

func (sc RedisContext) Get(key string, param interface{}) (err error) {
	cmd := sc.RdsClient.Get(sc.Ctx, key)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, param)
}
func (sc RedisContext) Set(key string, param interface{}, expiration time.Duration) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := sc.RdsClient.Set(sc.Ctx, key, bytes, expiration)
	return status.Err()
}
func (sc RedisContext) HGet(key string, field string, param interface{}) (err error) {
	cmd := sc.RdsClient.HGet(sc.Ctx, key, field)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, param)
}
func (sc RedisContext) HSet(key string, field string, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := sc.RdsClient.HSet(sc.Ctx, key, field, bytes)
	return status.Err()
}
func (sc RedisContext) RPush(key string, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := sc.RdsClient.RPush(sc.Ctx, key, bytes)
	return status.Err()
}
func (sc RedisContext) LSet(key string, index int64, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := sc.RdsClient.LSet(sc.Ctx, key, index, bytes)
	return status.Err()
}
func (sc RedisContext) LGet(key string, index int64, param interface{}) (err error) {
	cmd := sc.RdsClient.LIndex(sc.Ctx, key, index)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, param)
}
func (sc RedisContext) LLen(key string) (length int64) {
	cmd := sc.RdsClient.LLen(sc.Ctx, key)
	return cmd.Val()
}
