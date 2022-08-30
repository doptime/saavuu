package redis

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
