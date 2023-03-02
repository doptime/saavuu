package data

import (
	"errors"
	"reflect"

	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/logger"
	"github.com/yangkequn/saavuu/rds"
)

func (db *Ctx) HGet(key string, field interface{}, value interface{}) (err error) {
	return rds.HGet(db.Ctx, db.Rds, key, field, value)
}

func (db *Ctx) HSet(key string, field interface{}, value interface{}) (err error) {
	return rds.HSet(db.Ctx, db.Rds, key, field, value)
}

func (db *Ctx) HExists(key string, field string) (ok bool, err error) {
	cmd := db.Rds.HExists(db.Ctx, key, field)
	return cmd.Val(), cmd.Err()
}
func (db *Ctx) HGetAll(key string, mapOut interface{}) (err error) {
	return rds.HGetAll(db.Ctx, db.Rds, key, mapOut)
}
func (db *Ctx) HGetMap(key string, mapOut interface{}) (err error) {
	return rds.HGetMapPackFields(db.Ctx, db.Rds, key, mapOut)
}
func (db *Ctx) HSetMap(key string, _map interface{}) (err error) {
	return rds.HSetMapPackFields(db.Ctx, db.Rds, key, _map)
}
func (db *Ctx) HMGET(key string, _map interface{}, fields ...string) (err error) {
	mapElem := reflect.TypeOf(_map)
	if (mapElem.Kind() != reflect.Map) || (mapElem.Key().Kind() != reflect.String) {
		logger.Lshortfile.Println("mapOut must be a map[string] struct/interface{}")
		return errors.New("mapOut must be a map[string] struct/interface{}")
	}
	structSupposed := mapElem.Elem()
	cmd := db.Rds.HMGet(db.Ctx, key, fields...)
	if cmd.Err() == nil {
		//unmarshal each value of cmd.Val() to interface{}, using msgpack
		for i, v := range cmd.Val() {
			if v == nil {
				//set _map with nil
				reflect.ValueOf(_map).SetMapIndex(reflect.ValueOf(fields[i]), reflect.Zero(structSupposed))
				continue
			}
			obj := reflect.New(structSupposed).Interface()
			if err = msgpack.Unmarshal([]byte(v.(string)), &obj); err == nil {
				reflect.ValueOf(_map).SetMapIndex(reflect.ValueOf(fields[i]), reflect.ValueOf(obj).Elem())
			}
		}
	}
	return cmd.Err()
}
func (db *Ctx) HMGETPackFields(key string, fields []interface{}, values *[]interface{}) (err error) {
	return rds.HMGETPackFields(db.Ctx, db.Rds, key, fields, values)
}

func (db *Ctx) HGetAllDefault(key string) (param map[string]interface{}, err error) {
	cmd := db.Rds.HGetAll(db.Ctx, key)
	data, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	param = make(map[string]interface{})
	//make a copoy of valueStruct
	// unmarshal value of data to the copy
	// store unmarshaled result to param
	for k, v := range data {
		var obj interface{}
		if err = msgpack.Unmarshal([]byte(v), &obj); err == nil {
			param[k] = obj
		}
	}
	return param, nil
}
func (db *Ctx) HLen(key string) (length int64, err error) {
	cmd := db.Rds.HLen(db.Ctx, key)
	return cmd.Val(), cmd.Err()
}
func (db *Ctx) HDel(key string, field string) (err error) {
	status := db.Rds.HDel(db.Ctx, key, field)
	return status.Err()
}
func (db *Ctx) HKeys(key string, fields interface{}) (err error) {
	return rds.HKeys(db.Ctx, db.Rds, key, fields)
}
func (db *Ctx) HVals(key string, values *[]interface{}) (err error) {
	return rds.HValsPackFields(db.Ctx, db.Rds, key, values)
}
func (db *Ctx) HIncrBy(key string, field string, increment int64) (err error) {
	status := db.Rds.HIncrBy(db.Ctx, key, field, increment)
	return status.Err()
}
func (db *Ctx) HIncrByFloat(key string, field string, increment float64) (err error) {
	status := db.Rds.HIncrByFloat(db.Ctx, key, field, increment)
	return status.Err()
}
func (db *Ctx) HSetNX(key string, field string, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := db.Rds.HSetNX(db.Ctx, key, field, bytes)
	return status.Err()
}

// golang version of python scan_iter
func (db *Ctx) Scan(match string, cursor uint64, count int64) (keys []string, err error) {
	keys, _, err = db.Rds.Scan(db.Ctx, cursor, match, count).Result()
	return keys, err
}
