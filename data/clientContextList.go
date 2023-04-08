package data

import (
	"reflect"

	rds "github.com/yangkequn/saavuu/rds"
)

func (db *Ctx[v]) RPush(param v) (err error) {
	return rds.RPush(db.Ctx, db.Rds, db.Key, param)
}
func (db *Ctx[v]) LSet(index int64, param v) (err error) {
	return rds.LSet(db.Ctx, db.Rds, db.Key, index, param)
}
func (db *Ctx[v]) LGet(index int64) (param v, err error) {
	//create param with type v
	param = reflect.New(reflect.TypeOf((*v)(nil)).Elem()).Interface().(v)
	return param, rds.LGet(db.Ctx, db.Rds, db.Key, index, &param)
}
func (db *Ctx[v]) LLen() (length int64, err error) {
	return rds.LLen(db.Ctx, db.Rds, db.Key)
}
