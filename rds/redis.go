package rds

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func Time(ctx context.Context, rc *redis.Client) (time time.Time, err error) {
	cmd := rc.Time(ctx)
	return cmd.Result()
}

// sacn key by pattern
func Scan(ctx context.Context, rc *redis.Client, cursor uint64, match string, count int64) (keys []string, nextCursor uint64, err error) {
	cmd := rc.Scan(ctx, cursor, match, count)
	return cmd.Result()
}
