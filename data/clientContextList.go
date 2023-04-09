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

	vType := reflect.TypeOf((*v)(nil)).Elem()
	if vType.Kind() == reflect.Ptr {
		vValue := reflect.New(vType.Elem()).Interface().(v)
		err = rds.LGet(db.Ctx, db.Rds, db.Key, index, vValue)
		return vValue, err
	}
	vValueWithPointer := reflect.New(vType).Interface().(*v)
	err = rds.LGet(db.Ctx, db.Rds, db.Key, index, vValueWithPointer)
	return *vValueWithPointer, err
}
func (db *Ctx[v]) LLen() (length int64, err error) {
	return rds.LLen(db.Ctx, db.Rds, db.Key)
}
