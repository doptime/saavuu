package data

import (
	"reflect"

	"github.com/yangkequn/saavuu/rds"
)

func (db *Ctx[v]) HGet(field interface{}) (value v, err error) {
	vType := reflect.TypeOf((*v)(nil)).Elem()
	if vType.Kind() == reflect.Ptr {
		vValue := reflect.New(vType.Elem()).Interface().(v)
		err = rds.HGet(db.Ctx, db.Rds, db.Key, field, vValue)
		return vValue, err
	}
	vValueWithPointer := reflect.New(vType).Interface().(*v)
	err = rds.HGet(db.Ctx, db.Rds, db.Key, field, vValueWithPointer)
	return *vValueWithPointer, err
}

func (db *Ctx[v]) HSet(field interface{}, value v) (err error) {
	return rds.HSet(db.Ctx, db.Rds, db.Key, field, value)
}

func (db *Ctx[v]) HExists(field v) (ok bool, err error) {
	return rds.HExists(db.Ctx, db.Rds, db.Key, field)
}
func (db *Ctx[v]) HGetAll(mapOut interface{}) (err error) {
	return rds.HGetAll(db.Ctx, db.Rds, db.Key, mapOut)
}
func (db *Ctx[v]) HSetAll(_map interface{}) (err error) {
	return rds.HSetAll(db.Ctx, db.Rds, db.Key, _map)
}
func (db *Ctx[v]) HMGET(fields interface{}, mapOut interface{}) (err error) {
	return rds.HMGET(db.Ctx, db.Rds, db.Key, fields, mapOut)
}

func (db *Ctx[v]) HLen() (length int64, err error) {
	cmd := db.Rds.HLen(db.Ctx, db.Key)
	return cmd.Val(), cmd.Err()
}
func (db *Ctx[v]) HDel(field interface{}) (err error) {
	return rds.HDel(db.Ctx, db.Rds, db.Key, field)
}
func (db *Ctx[v]) HKeys(fields interface{}) (err error) {
	return rds.HKeys(db.Ctx, db.Rds, db.Key, fields)
}
func (db *Ctx[v]) HVals(values interface{}) (err error) {
	return rds.HVals(db.Ctx, db.Rds, db.Key, values)
}
func (db *Ctx[v]) HIncrBy(field interface{}, increment int64) (err error) {
	return rds.HIncrBy(db.Ctx, db.Rds, db.Key, field, increment)
}
func (db *Ctx[v]) HIncrByFloat(field string, increment float64) (err error) {
	return rds.HIncrByFloat(db.Ctx, db.Rds, db.Key, field, increment)
}
func (db *Ctx[v]) HSetNX(field interface{}, param interface{}) (err error) {
	return rds.HSetNX(db.Ctx, db.Rds, db.Key, field, param)
}

// golang version of python scan_iter
func (db *Ctx[v]) Scan(match string, cursor uint64, count int64) (keys []string, err error) {
	keys, _, err = db.Rds.Scan(db.Ctx, cursor, match, count).Result()
	return keys, err
}
