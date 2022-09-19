package saavuu

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v5"
)

type RedisContext struct {
	RdisClient *redis.Client
	ctx        context.Context
}

func (sc RedisContext) Get(key string, param interface{}) (err error) {
	cmd := sc.RdisClient.Get(sc.ctx, key)
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
	status := sc.RdisClient.Set(sc.ctx, key, bytes, expiration)
	return status.Err()
}
func (sc RedisContext) HGet(key string, field string, param interface{}) (err error) {
	cmd := sc.RdisClient.HGet(sc.ctx, key, field)
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
	status := sc.RdisClient.HSet(sc.ctx, key, field, bytes)
	return status.Err()
}
func (sc RedisContext) RPush(key string, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := sc.RdisClient.RPush(sc.ctx, key, bytes)
	return status.Err()
}
func (sc RedisContext) LSet(key string, index int64, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := sc.RdisClient.LSet(sc.ctx, key, index, bytes)
	return status.Err()
}
func (sc RedisContext) LGet(key string, index int64, param interface{}) (err error) {
	cmd := sc.RdisClient.LIndex(sc.ctx, key, index)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, param)
}
