package data

import (
	"time"

	"github.com/redis/go-redis/v9"
)

func (db *Ctx[k, v]) Get() (value v, err error) {
	var (
		cmd  *redis.StringCmd
		data []byte
	)

	if cmd = db.Rds.Get(db.Ctx, db.Key); cmd.Err() != nil {
		return value, cmd.Err()
	}
	if data, err = cmd.Bytes(); err != nil {
		return value, err
	}
	return db.toValue(data)
}
func (db *Ctx[k, v]) Set(param v, expiration time.Duration) (err error) {
	if val, err := db.toValueStr(param); err != nil {
		return err
	} else {
		status := db.Rds.Set(db.Ctx, db.Key, val, expiration)
		return status.Err()
	}
}
func (db *Ctx[k, v]) Del() (err error) {
	status := db.Rds.Del(db.Ctx, db.Key)
	return status.Err()
}
