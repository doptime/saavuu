package data

import (
	"fmt"
	"reflect"

	"github.com/jinzhu/copier"
	"github.com/vmihailenco/msgpack/v5"
)

func (db *Ctx[v]) upgradeSchemaFromRawString(raw string, upgrade func(in v) (out v)) (out v, err error) {
	var (
		vType       = reflect.TypeOf((*v)(nil)).Elem()
		vIsPtr bool = vType.Kind() == reflect.Ptr
		vValue interface{}
		obj    interface{}
	)
	if vIsPtr {
		vValue = reflect.New(vType.Elem()).Interface().(v)
	} else {
		vValue = reflect.New(vType).Interface().(*v)
	}
	//  step1: Advance conversion using copier, i.e. copy from string to float32
	msgpack.Unmarshal([]byte(raw), &obj)
	copier.Copy(vValue, &obj)
	// step2: read raw format from redis
	msgpack.Unmarshal([]byte(raw), vValue)
	// step3: upgrade scheme using user defined function
	// this is used, to allow user break point to work

	if vIsPtr {
		return upgrade(vValue.(v)), nil
	} else {
		return upgrade(*vValue.(*v)), nil
	}
}

func (db *Ctx[v]) UpgradeSchema(upgrader func(in v) (out v)) (err error) {
	var (
		keyType string
	)
	//get redis key type. if is hash, then upgrade all hash fields. if is list, then upgrade all list items. if is set, then upgrade all set members. if is zset, then upgrade all zset members.  if is string, then upgrade this string.
	//get redis key type
	if keyType, err = db.Rds.Type(db.Ctx, db.Key).Result(); err != nil {
		return err
	}
	switch keyType {
	case "hash":
		var mapIn map[string]string
		if mapIn, err = db.Rds.HGetAll(db.Ctx, db.Key).Result(); err != nil {
			return err
		}
		var mapOut map[string]interface{} = make(map[string]interface{})
		for k, v := range mapIn {
			if mapOut[k], err = db.upgradeSchemaFromRawString(v, upgrader); err != nil {
				return err
			}
		}
		return db.HSetAll(mapOut)
	case "list":
		//return not implemented error
		return fmt.Errorf("not implemented")
	case "set":
		return fmt.Errorf("not implemented")
	case "zset":
		return fmt.Errorf("not implemented")
	case "string":
		return fmt.Errorf("not implemented")
	}
	return nil
}
