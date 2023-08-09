package data

import (
	"github.com/yangkequn/saavuu/rds"
)

// append to Set
func (db *Ctx[k, v]) SAdd(param v) (err error) {
	return rds.SAdd(db.Ctx, db.Rds, db.Key, param)
}
func (db *Ctx[k, v]) SRem(param v) (err error) {
	return rds.SRem(db.Ctx, db.Rds, db.Key, param)
}
func (db *Ctx[k, v]) SIsMember(param v) (isMember bool, err error) {
	return rds.SIsMember(db.Ctx, db.Rds, db.Key, param)
}
func (db *Ctx[k, v]) SMembers() (members []string, err error) {
	return rds.SMembers(db.Ctx, db.Rds, db.Key)
}
