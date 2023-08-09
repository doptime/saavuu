package data

import (
	"encoding/json"
	"reflect"

	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/yangkequn/saavuu/rds"
)

func (db *Ctx[k, v]) HGet(field interface{}) (value v, err error) {
	vType := reflect.TypeOf((*v)(nil)).Elem()
	if vType.Kind() == reflect.Ptr {
		vValue := reflect.New(vType.Elem()).Interface().(v)
		err = rds.HGet(db.Ctx, db.Rds, db.Key, field, vValue)
		return vValue, err
	}
	vValueWithPointer := reflect.New(vType).Interface().(*v)
	err = rds.HGet(db.Ctx, db.Rds, db.Key, field, vValueWithPointer)
	return *vValueWithPointer, err
}

func (db *Ctx[k, v]) HSet(field interface{}, value v) (err error) {
	return rds.HSet(db.Ctx, db.Rds, db.Key, field, value)
}

func (db *Ctx[k, v]) HExists(field v) (ok bool, err error) {
	return rds.HExists(db.Ctx, db.Rds, db.Key, field)
}
func (db *Ctx[k, v]) HGetAll(mapOut interface{}) (err error) {
	return rds.HGetAll(db.Ctx, db.Rds, db.Key, mapOut)
}
func (db *Ctx[k, v]) HSetAll(_map interface{}) (err error) {
	return rds.HSetAll(db.Ctx, db.Rds, db.Key, _map)
}

func (db *Ctx[k, v]) HMSET(values []v) (err error) {
	var (
		valueBytes [][]byte
		bytes      []byte
	)
	for _, value := range values {
		if bytes, err = msgpack.Marshal(value); err != nil {
			return err
		}
		valueBytes = append(valueBytes, bytes)
	}
	cmd := db.Rds.HMSet(db.Ctx, db.Key, valueBytes)
	return cmd.Err()
}
func (db *Ctx[k, v]) HMGET(fields interface{}) (values []v, err error) {
	var (
		cmd          *redis.SliceCmd
		fieldsString []string
	)
	if fieldsString, err = rds.FieldsToSlice(fields); err != nil {
		return nil, err
	}

	if cmd = db.Rds.HMGet(db.Ctx, db.Key, fieldsString...); cmd.Err() != nil {
		return nil, cmd.Err()
	}
	values = make([]v, len(fieldsString))
	valueStruct := reflect.TypeOf((*v)(nil)).Elem()
	isElemPtr := valueStruct.Kind() == reflect.Ptr

	//save all data to mapOut
	for i, val := range cmd.Val() {
		if val == nil {
			continue
		}
		obj := reflect.New(valueStruct).Interface()
		if isElemPtr {
			if err = msgpack.Unmarshal([]byte(val.(string)), obj); err == nil {
				values[i] = *obj.(*v)
			}
			continue
		}
		if err = msgpack.Unmarshal([]byte(val.(string)), &obj); err == nil {
			values[i] = obj.(v)
		}
	}
	return values, nil
}

func (db *Ctx[k, v]) HLen() (length int64, err error) {
	cmd := db.Rds.HLen(db.Ctx, db.Key)
	return cmd.Val(), cmd.Err()
}
func (db *Ctx[k, v]) HDel(fields ...k) (err error) {
	var (
		cmd       *redis.IntCmd
		fieldStrs []string
		bytes     []byte
	)
	if len(fields) == 0 {
		return nil
	}
	//if k is  string, then use HDEL directly
	if reflect.TypeOf(fields[0]).Kind() == reflect.String {
		fieldStrs = interface{}(fields).([]string)
	} else {
		//if k is not string, then marshal k to string
		fieldStrs = make([]string, len(fields))
		for i, field := range fields {
			if bytes, err = json.Marshal(field); err != nil {
				return err
			}
			fieldStrs[i] = string(bytes)
		}
		return err
	}
	cmd = db.Rds.HDel(db.Ctx, db.Key, fieldStrs...)
	return cmd.Err()
}
func (db *Ctx[k, v]) HKeys(fields interface{}) (err error) {
	return rds.HKeys(db.Ctx, db.Rds, db.Key, fields)
}
func (db *Ctx[k, v]) HVals() (values []v, err error) {
	values = make([]v, 0)
	return values, rds.HVals(db.Ctx, db.Rds, db.Key, &values)
}
func (db *Ctx[k, v]) HIncrBy(field interface{}, increment int64) (err error) {
	return rds.HIncrBy(db.Ctx, db.Rds, db.Key, field, increment)
}
func (db *Ctx[k, v]) HIncrByFloat(field string, increment float64) (err error) {
	return rds.HIncrByFloat(db.Ctx, db.Rds, db.Key, field, increment)
}
func (db *Ctx[k, v]) HSetNX(field interface{}, param interface{}) (err error) {
	return rds.HSetNX(db.Ctx, db.Rds, db.Key, field, param)
}

// golang version of python scan_iter
func (db *Ctx[k, v]) Scan(match string, cursor uint64, count int64) (keys []string, err error) {
	keys, _, err = db.Rds.Scan(db.Ctx, cursor, match, count).Result()
	return keys, err
}
