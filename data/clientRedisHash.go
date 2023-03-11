package data

import (
	"github.com/yangkequn/saavuu/rds"
)

func (db *Ctx) HGet(field interface{}, value interface{}) (err error) {
	return rds.HGet(db.Ctx, db.Rds, db.Key, field, value)
}

func (db *Ctx) HSet(field interface{}, value interface{}) (err error) {
	return rds.HSet(db.Ctx, db.Rds, db.Key, field, value)
}

func (db *Ctx) HExists(field interface{}) (ok bool, err error) {
	return rds.HExists(db.Ctx, db.Rds, db.Key, field)
}
func (db *Ctx) HGetAll(mapOut interface{}) (err error) {
	return rds.HGetAll(db.Ctx, db.Rds, db.Key, mapOut)
}
func (db *Ctx) HSetAll(_map interface{}) (err error) {
	return rds.HSetAll(db.Ctx, db.Rds, db.Key, _map)
}
func (db *Ctx) HMGET(fields interface{}, mapOut interface{}) (err error) {
	return rds.HMGET(db.Ctx, db.Rds, db.Key, fields, mapOut)
}

func (db *Ctx) HLen() (length int64, err error) {
	cmd := db.Rds.HLen(db.Ctx, db.Key)
	return cmd.Val(), cmd.Err()
}
func (db *Ctx) HDel(field interface{}) (err error) {
	return rds.HDel(db.Ctx, db.Rds, db.Key, field)
}
func (db *Ctx) HKeys(fields interface{}) (err error) {
	return rds.HKeys(db.Ctx, db.Rds, db.Key, fields)
}
func (db *Ctx) HVals(values interface{}) (err error) {
	return rds.HVals(db.Ctx, db.Rds, db.Key, values)
}
func (db *Ctx) HIncrBy(field interface{}, increment int64) (err error) {
	return rds.HIncrBy(db.Ctx, db.Rds, db.Key, field, increment)
}
func (db *Ctx) HIncrByFloat(field string, increment float64) (err error) {
	return rds.HIncrByFloat(db.Ctx, db.Rds, db.Key, field, increment)
}
func (db *Ctx) HSetNX(field interface{}, param interface{}) (err error) {
	return rds.HSetNX(db.Ctx, db.Rds, db.Key, field, param)
}

// golang version of python scan_iter
func (db *Ctx) Scan(match string, cursor uint64, count int64) (keys []string, err error) {
	keys, _, err = db.Rds.Scan(db.Ctx, cursor, match, count).Result()
	return keys, err
}
