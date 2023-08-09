package data

import (
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

func (db *Ctx[k, v]) ZAdd(members ...redis.Z) (err error) {
	//MarshalRedisZ
	for i := range members {
		if members[i].Member != nil {
			members[i].Member, _ = msgpack.Marshal(members[i].Member)
		}
	}
	status := db.Rds.ZAdd(db.Ctx, db.Key, members...)
	return status.Err()
}
func (db *Ctx[k, v]) ZRem(members ...interface{}) (err error) {
	//msgpack marshal members to slice of bytes
	var bytes = make([][]byte, len(members))
	for i, member := range members {
		if bytes[i], err = msgpack.Marshal(member); err != nil {
			return err
		}
	}
	var redisPipe = db.Rds.Pipeline()
	for _, memberBytes := range bytes {
		redisPipe.ZRem(db.Ctx, db.Key, memberBytes)
	}
	_, err = redisPipe.Exec(db.Ctx)

	return err
}
func (db *Ctx[k, v]) ZRange(start, stop int64) (members []v, err error) {
	var cmd *redis.StringSliceCmd

	if cmd = db.Rds.ZRange(db.Ctx, db.Key, start, stop); cmd.Err() != nil && cmd.Err() != redis.Nil {
		return nil, cmd.Err()
	}
	return db.UnmarshalToSlice(cmd.Val())
}
func (db *Ctx[k, v]) ZRangeWithScores(start, stop int64) (members []v, scores []float64, err error) {
	cmd := db.Rds.ZRangeWithScores(db.Ctx, db.Key, start, stop)
	return db.UnmarshalRedisZ(cmd.Val())
}
func (db *Ctx[k, v]) ZRevRangeWithScores(start, stop int64) (members []v, scores []float64, err error) {
	cmd := db.Rds.ZRevRangeWithScores(db.Ctx, db.Key, start, stop)
	return db.UnmarshalRedisZ(cmd.Val())
}
func (db *Ctx[k, v]) ZRank(member interface{}) (rank int64, err error) {
	var (
		memberBytes []byte
	)
	//marshal member using msgpack
	if memberBytes, err = msgpack.Marshal(member); err != nil {
		return 0, err
	}
	cmd := db.Rds.ZRank(db.Ctx, db.Key, string(memberBytes))
	return cmd.Val(), cmd.Err()
}
func (db *Ctx[k, v]) ZRevRank(member interface{}) (rank int64, err error) {
	var (
		memberBytes []byte
	)
	//marshal member using msgpack
	if memberBytes, err = msgpack.Marshal(member); err != nil {
		return 0, err
	}
	cmd := db.Rds.ZRevRank(db.Ctx, db.Key, string(memberBytes))
	return cmd.Val(), cmd.Err()
}
func (db *Ctx[k, v]) ZScore(member interface{}) (score float64, err error) {
	var (
		memberBytes []byte
		cmd         *redis.FloatCmd
	)
	//marshal member using msgpack
	if memberBytes, err = msgpack.Marshal(member); err != nil {
		return 0, err
	}
	if cmd = db.Rds.ZScore(db.Ctx, db.Key, string(memberBytes)); cmd.Err() != nil && cmd.Err() != redis.Nil {
		return 0, err
	} else if cmd.Err() == redis.Nil {
		return 0, nil
	}
	return cmd.Result()
}
func (db *Ctx[k, v]) ZCard() (length int64, err error) {
	cmd := db.Rds.ZCard(db.Ctx, db.Key)
	return cmd.Result()
}
func (db *Ctx[k, v]) ZCount(min, max string) (length int64, err error) {
	cmd := db.Rds.ZCount(db.Ctx, db.Key, min, max)
	return cmd.Result()
}
func (db *Ctx[k, v]) ZRangeByScore(opt *redis.ZRangeBy) (out []v, err error) {
	cmd := db.Rds.ZRangeByScore(db.Ctx, db.Key, opt)

	return db.UnmarshalToSlice(cmd.Val())
}
func (db *Ctx[k, v]) ZRangeByScoreWithScores(opt *redis.ZRangeBy) (out []v, scores []float64, err error) {
	cmd := db.Rds.ZRangeByScoreWithScores(db.Ctx, db.Key, opt)
	return db.UnmarshalRedisZ(cmd.Val())
}
func (db *Ctx[k, v]) ZRevRangeByScore(opt *redis.ZRangeBy) (out []v, err error) {
	cmd := db.Rds.ZRevRangeByScore(db.Ctx, db.Key, opt)
	return db.UnmarshalToSlice(cmd.Val())
}
func (db *Ctx[k, v]) ZRevRange(start, stop int64) (out []v, err error) {
	var cmd *redis.StringSliceCmd

	if cmd = db.Rds.ZRevRange(db.Ctx, db.Key, start, stop); cmd.Err() != nil && cmd.Err() != redis.Nil {
		return nil, cmd.Err()
	}
	return db.UnmarshalToSlice(cmd.Val())
}
func (db *Ctx[k, v]) ZRevRangeByScoreWithScores(opt *redis.ZRangeBy) (out []v, scores []float64, err error) {
	cmd := db.Rds.ZRevRangeByScoreWithScores(db.Ctx, db.Key, opt)
	return db.UnmarshalRedisZ(cmd.Val())
}
func (db *Ctx[k, v]) ZRemRangeByRank(start, stop int64) (err error) {
	status := db.Rds.ZRemRangeByRank(db.Ctx, db.Key, start, stop)
	return status.Err()
}
func (db *Ctx[k, v]) ZRemRangeByScore(min, max string) (err error) {
	status := db.Rds.ZRemRangeByScore(db.Ctx, db.Key, min, max)
	return status.Err()
}
func (db *Ctx[k, v]) ZIncrBy(increment float64, member interface{}) (err error) {
	var (
		memberBytes []byte
	)
	//marshal member using msgpack
	if memberBytes, err = msgpack.Marshal(member); err != nil {
		return err
	}
	status := db.Rds.ZIncrBy(db.Ctx, db.Key, increment, string(memberBytes))
	return status.Err()
}
func (db *Ctx[k, v]) ZPopMax(count int64) (out []v, scores []float64, err error) {
	cmd := db.Rds.ZPopMax(db.Ctx, db.Key, count)
	return db.UnmarshalRedisZ(cmd.Val())
}
func (db *Ctx[k, v]) ZPopMin(count int64) (out []v, scores []float64, err error) {
	cmd := db.Rds.ZPopMin(db.Ctx, db.Key, count)
	return db.UnmarshalRedisZ(cmd.Val())
}
func (db *Ctx[k, v]) ZLexCount(min, max string) (length int64) {
	cmd := db.Rds.ZLexCount(db.Ctx, db.Key, min, max)
	return cmd.Val()
}
func (db *Ctx[k, v]) ZScan(cursor uint64, match string, count int64) ([]string, uint64, error) {
	cmd := db.Rds.ZScan(db.Ctx, db.Key, cursor, match, count)
	return cmd.Result()
}
