package rds

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

func Get(ctx context.Context, rc *redis.Client, key string, param interface{}) (err error) {
	cmd := rc.Get(ctx, key)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, param)
}
func Set(ctx context.Context, rc *redis.Client, key string, param interface{}, expiration time.Duration) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := rc.Set(ctx, key, bytes, expiration)
	return status.Err()
}
