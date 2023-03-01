package rds

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

func SAdd(ctx context.Context, rc *redis.Client, key string, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := rc.SAdd(ctx, key, bytes)
	return status.Err()
}
func SRem(ctx context.Context, rc *redis.Client, key string, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := rc.SRem(ctx, key, bytes)
	return status.Err()
}
func SIsMember(ctx context.Context, rc *redis.Client, key string, param interface{}) (isMember bool, err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return false, err
	}
	cmd := rc.SIsMember(ctx, key, bytes)
	return cmd.Result()
}
func SMembers(ctx context.Context, rc *redis.Client, key string) (members []string, err error) {
	cmd := rc.SMembers(ctx, key)
	return cmd.Result()
}
