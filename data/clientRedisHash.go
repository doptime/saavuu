package data

import (
	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/rds"
)

func (db *Ctx) HGet(key string, field interface{}, value interface{}) (err error) {
	return rds.HGet(db.Ctx, db.Rds, key, field, value)
}

func (db *Ctx) HSet(key string, field interface{}, value interface{}) (err error) {
	return rds.HSet(db.Ctx, db.Rds, key, field, value)
}

func (db *Ctx) HExists(key string, field string) (ok bool, err error) {
	cmd := db.Rds.HExists(db.Ctx, key, field)
	return cmd.Val(), cmd.Err()
}
func (db *Ctx) HGetAll(key string, mapOut interface{}) (err error) {
	return rds.HGetAll(db.Ctx, db.Rds, key, mapOut)
}
func (db *Ctx) HSetAll(key string, _map interface{}) (err error) {
	return rds.HSetAll(db.Ctx, db.Rds, key, _map)
}
func (db *Ctx) HMGET(key string, fields interface{}, mapOut interface{}) (err error) {
	return rds.HMGET(db.Ctx, db.Rds, key, fields, mapOut)
}

func (db *Ctx) HLen(key string) (length int64, err error) {
	cmd := db.Rds.HLen(db.Ctx, key)
	return cmd.Val(), cmd.Err()
}
func (db *Ctx) HDel(key string, field string) (err error) {
	status := db.Rds.HDel(db.Ctx, key, field)
	return status.Err()
}
func (db *Ctx) HKeys(key string, fields interface{}) (err error) {
	return rds.HKeys(db.Ctx, db.Rds, key, fields)
}
func (db *Ctx) HVals(key string, values *[]interface{}) (err error) {
	return rds.HValsPackFields(db.Ctx, db.Rds, key, values)
}
func (db *Ctx) HIncrBy(key string, field string, increment int64) (err error) {
	status := db.Rds.HIncrBy(db.Ctx, key, field, increment)
	return status.Err()
}
func (db *Ctx) HIncrByFloat(key string, field string, increment float64) (err error) {
	status := db.Rds.HIncrByFloat(db.Ctx, key, field, increment)
	return status.Err()
}
func (db *Ctx) HSetNX(key string, field string, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := db.Rds.HSetNX(db.Ctx, key, field, bytes)
	return status.Err()
}

// golang version of python scan_iter
func (db *Ctx) Scan(match string, cursor uint64, count int64) (keys []string, err error) {
	keys, _, err = db.Rds.Scan(db.Ctx, cursor, match, count).Result()
	return keys, err
}
