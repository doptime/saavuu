package data

import (
	"time"

	"github.com/yangkequn/saavuu/rds"
)

func (db *Ctx) Get(param interface{}) (err error) {
	return rds.Get(db.Ctx, db.Rds, db.Key, param)
}
func (db *Ctx) Set(param interface{}, expiration time.Duration) (err error) {
	return rds.Set(db.Ctx, db.Rds, db.Key, param, expiration)
}
