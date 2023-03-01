package rds

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

func RPush(ctx context.Context, rc *redis.Client, key string, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := rc.RPush(ctx, key, bytes)
	return status.Err()
}
func LSet(ctx context.Context, rc *redis.Client, key string, index int64, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := rc.LSet(ctx, key, index, bytes)
	return status.Err()
}
func LGet(ctx context.Context, rc *redis.Client, key string, index int64, param interface{}) (err error) {
	cmd := rc.LIndex(ctx, key, index)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, param)
}
func LLen(ctx context.Context, rc *redis.Client, key string) (length int64, err error) {
	cmd := rc.LLen(ctx, key)
	return cmd.Result()
}
