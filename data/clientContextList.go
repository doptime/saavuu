package data

import rds "github.com/yangkequn/saavuu/rds"

func (db *Ctx) RPush(param interface{}) (err error) {
	return rds.RPush(db.Ctx, db.Rds, db.Key, param)
}
func (db *Ctx) LSet(index int64, param interface{}) (err error) {
	return rds.LSet(db.Ctx, db.Rds, db.Key, index, param)
}
func (db *Ctx) LGet(index int64, param interface{}) (err error) {
	return rds.LGet(db.Ctx, db.Rds, db.Key, index, param)
}
func (db *Ctx) LLen() (length int64, err error) {
	return rds.LLen(db.Ctx, db.Rds, db.Key)
}
