package data

import "github.com/redis/go-redis/v9"

// append to Set
func (db *Ctx[k, v]) SAdd(param v) (err error) {
	valStr, err := db.toValueStr(param)
	if err != nil {
		return err
	}
	status := db.Rds.SAdd(db.Ctx, db.Key, valStr)
	return status.Err()
}
func (db *Ctx[k, v]) SRem(param v) (err error) {
	valStr, err := db.toValueStr(param)
	if err != nil {
		return err
	}
	status := db.Rds.SRem(db.Ctx, db.Key, valStr)
	return status.Err()
}
func (db *Ctx[k, v]) SIsMember(param v) (isMember bool, err error) {
	valStr, err := db.toValueStr(param)
	if err != nil {
		return false, err
	}
	status := db.Rds.SIsMember(db.Ctx, db.Key, valStr)
	return status.Result()
}
func (db *Ctx[k, v]) SMembers() (members []v, err error) {
	var cmd *redis.StringSliceCmd
	if cmd = db.Rds.SMembers(db.Ctx, db.Key); cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return db.toValues(cmd.Val()...)
}
