package data

import (
	"github.com/redis/go-redis/v9"
)

func (dc *Ctx) ZAdd(key string, members ...redis.Z) (err error) {
	status := dc.Rds.ZAdd(dc.Ctx, key, members...)
	return status.Err()
}
func (dc *Ctx) ZRem(key string, members ...interface{}) (err error) {
	status := dc.Rds.ZRem(dc.Ctx, key, members)
	return status.Err()
}
func (dc *Ctx) ZRange(key string, start, stop int64) (members []string, err error) {
	cmd := dc.Rds.ZRange(dc.Ctx, key, start, stop)
	return cmd.Result()
}
func (dc *Ctx) ZRangeWithScores(key string, start, stop int64) (members []redis.Z, err error) {
	cmd := dc.Rds.ZRangeWithScores(dc.Ctx, key, start, stop)
	return cmd.Result()
}
func (dc *Ctx) ZRevRangeWithScores(key string, start, stop int64) (members []redis.Z, err error) {
	cmd := dc.Rds.ZRevRangeWithScores(dc.Ctx, key, start, stop)
	return cmd.Result()
}
func (dc *Ctx) ZRank(key string, member string) (rank int64, err error) {
	cmd := dc.Rds.ZRank(dc.Ctx, key, member)
	return cmd.Val(), cmd.Err()
}
func (dc *Ctx) ZRevRank(key string, member string) (rank int64) {
	cmd := dc.Rds.ZRevRank(dc.Ctx, key, member)
	return cmd.Val()
}
func (dc *Ctx) ZScore(key string, member string) (score float64) {
	cmd := dc.Rds.ZScore(dc.Ctx, key, member)
	return cmd.Val()
}
func (dc *Ctx) ZCard(key string) (length int64) {
	cmd := dc.Rds.ZCard(dc.Ctx, key)
	return cmd.Val()
}
func (dc *Ctx) ZCount(key string, min, max string) (length int64) {
	cmd := dc.Rds.ZCount(dc.Ctx, key, min, max)
	return cmd.Val()
}
func (dc *Ctx) ZRangeByScoreWithScores(key string, opt *redis.ZRangeBy) (members []redis.Z, err error) {
	cmd := dc.Rds.ZRangeByScoreWithScores(dc.Ctx, key, opt)
	return cmd.Result()
}
func (dc *Ctx) ZRevRangeByScore(key string, opt *redis.ZRangeBy) (members []string, err error) {
	cmd := dc.Rds.ZRevRangeByScore(dc.Ctx, key, opt)
	return cmd.Result()
}
func (dc *Ctx) ZRevRangeByScoreWithScores(key string, opt *redis.ZRangeBy) (members []redis.Z, err error) {
	cmd := dc.Rds.ZRevRangeByScoreWithScores(dc.Ctx, key, opt)
	return cmd.Result()
}
func (dc *Ctx) ZRemRangeByRank(key string, start, stop int64) (err error) {
	status := dc.Rds.ZRemRangeByRank(dc.Ctx, key, start, stop)
	return status.Err()
}
func (dc *Ctx) ZRemRangeByScore(key string, min, max string) (err error) {
	status := dc.Rds.ZRemRangeByScore(dc.Ctx, key, min, max)
	return status.Err()
}
func (dc *Ctx) ZIncrBy(key string, increment float64, member string) (err error) {
	status := dc.Rds.ZIncrBy(dc.Ctx, key, increment, member)
	return status.Err()
}
func (dc *Ctx) ZUnionStore(destination string, store *redis.ZStore) (err error) {
	status := dc.Rds.ZUnionStore(dc.Ctx, destination, store)
	return status.Err()
}
func (dc *Ctx) ZInterStore(destination string, store *redis.ZStore) (err error) {
	status := dc.Rds.ZInterStore(dc.Ctx, destination, store)
	return status.Err()
}
func (dc *Ctx) ZPopMax(key string, count int64) (members []redis.Z, err error) {
	cmd := dc.Rds.ZPopMax(dc.Ctx, key, count)
	return cmd.Result()
}
func (dc *Ctx) ZPopMin(key string, count int64) (members []redis.Z, err error) {
	cmd := dc.Rds.ZPopMin(dc.Ctx, key, count)
	return cmd.Result()
}
func (dc *Ctx) ZLexCount(key string, min, max string) (length int64) {
	cmd := dc.Rds.ZLexCount(dc.Ctx, key, min, max)
	return cmd.Val()
}
func (dc *Ctx) ZScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	cmd := dc.Rds.ZScan(dc.Ctx, key, cursor, match, count)
	return cmd.Result()
}
