package data

import (
	"time"

	"github.com/redis/go-redis/v9"
)

func (db *Ctx[k, v]) Get(key k) (value v, err error) {
	var (
		keyStr string
		cmd    *redis.StringCmd
		data   []byte
	)
	if keyStr, err = db.toKeyStr(key); err != nil {
		return value, err
	}

	if cmd = db.Rds.Get(db.Ctx, db.Key+":"+keyStr); cmd.Err() != nil {
		return value, cmd.Err()
	}
	if data, err = cmd.Bytes(); err != nil {
		return value, err
	}
	return db.toValue(data)
}
func (db *Ctx[k, v]) Keys() (out []k, err error) {
	var (
		cmd  *redis.StringSliceCmd
		keys []string
	)
	cmd = db.Rds.Keys(db.Ctx, db.Key+":*")
	if keys, err = cmd.Result(); err != nil {
		return nil, err
	}
	return db.toKeys(keys)
}

func (db *Ctx[k, v]) Set(key k, param v, expiration time.Duration) (err error) {
	var (
		keyStr string
		valStr string
	)
	if keyStr, err = db.toKeyStr(key); err != nil {
		return err
	}
	if valStr, err = db.toValueStr(param); err != nil {
		return err
	} else {
		status := db.Rds.Set(db.Ctx, db.Key+":"+keyStr, valStr, expiration)
		return status.Err()
	}
}
func (db *Ctx[k, v]) Del(key k) (err error) {
	var (
		keyStr string
	)
	if keyStr, err = db.toKeyStr(key); err != nil {
		return err
	}
	status := db.Rds.Del(db.Ctx, keyStr)
	return status.Err()
}
