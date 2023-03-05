package data

import (
	"github.com/yangkequn/saavuu/rds"
)

// append to Set
func (db *Ctx) SAdd(param interface{}) (err error) {
	return rds.SAdd(db.Ctx, db.Rds, db.Key, param)
}
func (db *Ctx) SRem(param interface{}) (err error) {
	return rds.SRem(db.Ctx, db.Rds, db.Key, param)
}
func (db *Ctx) SIsMember(param interface{}) (isMember bool, err error) {
	return rds.SIsMember(db.Ctx, db.Rds, db.Key, param)
}
func (db *Ctx) SMembers() (members []string, err error) {
	return rds.SMembers(db.Ctx, db.Rds, db.Key)
}
