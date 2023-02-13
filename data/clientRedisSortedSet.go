package data

import (
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

func (db *Ctx) ZAdd(key string, members ...redis.Z) (err error) {
	db.MarshalRedisZ(members...)
	status := db.Rds.ZAdd(db.Ctx, key, members...)
	return status.Err()
}
func (db *Ctx) ZRem(key string, members ...interface{}) (err error) {
	//msgpack marshal members
	var memberBytes [][]byte
	if memberBytes, err = db.MarshalSlice(members...); err != nil {
		return err
	}
	status := db.Rds.ZRem(db.Ctx, key, memberBytes)
	return status.Err()
}
func (db *Ctx) ZRange(key string, start, stop int64, out interface{}) (err error) {
	var cmd *redis.StringSliceCmd

	if cmd = db.Rds.ZRange(db.Ctx, key, start, stop); cmd.Err() != nil && cmd.Err() != redis.Nil {
		return cmd.Err()
	}
	return db.UnmarshalStrings(cmd.Val(), out)
}
func (db *Ctx) ZRangeWithScores(key string, start, stop int64, outStruct interface{}) (results []redis.Z, err error) {
	cmd := db.Rds.ZRangeWithScores(db.Ctx, key, start, stop)
	return db.UnmarshalRedisZ(cmd.Val(), outStruct)
}
func (db *Ctx) ZRevRangeWithScores(key string, start, stop int64, outStruct interface{}) (results []redis.Z, err error) {
	cmd := db.Rds.ZRevRangeWithScores(db.Ctx, key, start, stop)
	return db.UnmarshalRedisZ(cmd.Val(), outStruct)
}
func (db *Ctx) ZRank(key string, member string) (rank int64, err error) {
	var (
		memberBytes []byte
	)
	//marshal member using msgpack
	if memberBytes, err = msgpack.Marshal(member); err != nil {
		return 0, err
	}
	cmd := db.Rds.ZRank(db.Ctx, key, string(memberBytes))
	return cmd.Val(), cmd.Err()
}
func (db *Ctx) ZRevRank(key string, member string) (rank int64) {
	var (
		memberBytes []byte
		err         error
	)
	//marshal member using msgpack
	if memberBytes, err = msgpack.Marshal(member); err != nil {
		return 0
	}
	cmd := db.Rds.ZRevRank(db.Ctx, key, string(memberBytes))
	return cmd.Val()
}
func (db *Ctx) ZScore(key string, member string) (score float64) {
	var (
		memberBytes []byte
		err         error
	)
	//marshal member using msgpack
	if memberBytes, err = msgpack.Marshal(member); err != nil {
		return 0
	}
	cmd := db.Rds.ZScore(db.Ctx, key, string(memberBytes))
	return cmd.Val()
}
func (db *Ctx) ZCard(key string) (length int64) {
	cmd := db.Rds.ZCard(db.Ctx, key)
	return cmd.Val()
}
func (db *Ctx) ZCount(key string, min, max string) (length int64) {
	cmd := db.Rds.ZCount(db.Ctx, key, min, max)
	return cmd.Val()
}
func (db *Ctx) ZRangeByScoreWithScores(key string, opt *redis.ZRangeBy, outStruct interface{}) (results []redis.Z, err error) {
	cmd := db.Rds.ZRangeByScoreWithScores(db.Ctx, key, opt)
	return db.UnmarshalRedisZ(cmd.Val(), outStruct)
}
func (db *Ctx) ZRevRangeByScore(key string, opt *redis.ZRangeBy) (members []string, err error) {
	cmd := db.Rds.ZRevRangeByScore(db.Ctx, key, opt)
	return cmd.Result()
}
func (db *Ctx) ZRevRangeByScoreWithScores(key string, opt *redis.ZRangeBy, outStruct interface{}) (results []redis.Z, err error) {
	cmd := db.Rds.ZRevRangeByScoreWithScores(db.Ctx, key, opt)
	return db.UnmarshalRedisZ(cmd.Val(), outStruct)
}
func (db *Ctx) ZRemRangeByRank(key string, start, stop int64) (err error) {
	status := db.Rds.ZRemRangeByRank(db.Ctx, key, start, stop)
	return status.Err()
}
func (db *Ctx) ZRemRangeByScore(key string, min, max string) (err error) {
	status := db.Rds.ZRemRangeByScore(db.Ctx, key, min, max)
	return status.Err()
}
func (db *Ctx) ZIncrBy(key string, increment float64, member string) (err error) {
	status := db.Rds.ZIncrBy(db.Ctx, key, increment, member)
	return status.Err()
}
func (db *Ctx) ZUnionStore(destination string, store *redis.ZStore) (err error) {
	status := db.Rds.ZUnionStore(db.Ctx, destination, store)
	return status.Err()
}
func (db *Ctx) ZInterStore(destination string, store *redis.ZStore) (err error) {
	status := db.Rds.ZInterStore(db.Ctx, destination, store)
	return status.Err()
}
func (db *Ctx) ZPopMax(key string, count int64, outStruct interface{}) (results []redis.Z, err error) {
	cmd := db.Rds.ZPopMax(db.Ctx, key, count)
	return db.UnmarshalRedisZ(cmd.Val(), outStruct)
}
func (db *Ctx) ZPopMin(key string, count int64, outStruct interface{}) (results []redis.Z, err error) {
	cmd := db.Rds.ZPopMin(db.Ctx, key, count)
	return db.UnmarshalRedisZ(cmd.Val(), outStruct)
}
func (db *Ctx) ZLexCount(key string, min, max string) (length int64) {
	cmd := db.Rds.ZLexCount(db.Ctx, key, min, max)
	return cmd.Val()
}
func (db *Ctx) ZScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	cmd := db.Rds.ZScan(db.Ctx, key, cursor, match, count)
	return cmd.Result()
}
