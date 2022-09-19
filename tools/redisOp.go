package tools

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v5"
)

func Get(c context.Context, rds *redis.Client, key string, param interface{}) (err error) {
	cmd := rds.Get(c, key)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, param)
}
func Set(c context.Context, rds *redis.Client, key string, param interface{}, expiration time.Duration) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := rds.Set(c, key, bytes, expiration)
	return status.Err()
}
func HGet(c context.Context, rds *redis.Client, key string, field string, param interface{}) (err error) {
	cmd := rds.HGet(c, key, field)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, param)
}

func HSet(c context.Context, rds *redis.Client, key string, field string, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := rds.HSet(c, key, field, bytes)
	return status.Err()
}
func RPush(c context.Context, rds *redis.Client, key string, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := rds.RPush(c, key, bytes)
	return status.Err()
}
func LSet(c context.Context, rds *redis.Client, key string, index int64, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := rds.LSet(c, key, index, bytes)
	return status.Err()
}
func LGet(c context.Context, rds *redis.Client, key string, index int64, param interface{}) (err error) {
	cmd := rds.LIndex(c, key, index)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, param)
}
