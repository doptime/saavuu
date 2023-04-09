package data

import (
	"reflect"
	"time"

	"github.com/yangkequn/saavuu/rds"
)

func (db *Ctx[v]) Get() (value v, err error) {

	vType := reflect.TypeOf((*v)(nil)).Elem()
	if vType.Kind() == reflect.Ptr {
		vValue := reflect.New(vType.Elem()).Interface().(v)
		err = rds.Get(db.Ctx, db.Rds, db.Key, vValue)
		return vValue, err
	}
	vValueWithPointer := reflect.New(vType).Interface().(*v)
	err = rds.Get(db.Ctx, db.Rds, db.Key, vValueWithPointer)
	return *vValueWithPointer, err
}
func (db *Ctx[v]) Set(param v, expiration time.Duration) (err error) {
	return rds.Set(db.Ctx, db.Rds, db.Key, param, expiration)
}
func (db *Ctx[v]) Del() (err error) {
	return rds.Del(db.Ctx, db.Rds, db.Key)
}
