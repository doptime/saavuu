package data

import (
	"encoding/json"
	"reflect"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func (db *Ctx[k, v]) HGet(field k) (value v, err error) {

	var (
		cmd      *redis.StringCmd
		valBytes []byte
		fieldStr string
	)
	if fieldStr, err = db.toKeyStr(field); err != nil {
		return value, err
	}

	if cmd = db.Rds.HGet(db.Ctx, db.Key, fieldStr); cmd.Err() != nil {
		return value, cmd.Err()
	}
	if valBytes, err = cmd.Bytes(); err != nil {
		return value, err
	}
	return db.toValue(valBytes)
}

// HSet accepts values in following formats:
//
//   - HSet("myhash", "key1", "value1", "key2", "value2")
//
//   - HSet("myhash", map[string]interface{}{"key1": "value1", "key2": "value2"})
func (db *Ctx[k, v]) HSet(values ...interface{}) (err error) {
	var (
		KeyValuesStrs []string
	)
	if KeyValuesStrs, err = db.toKeyValueStrs(values...); err != nil {
		return err
	}
	status := db.Rds.HSet(db.Ctx, db.Key, KeyValuesStrs)
	return status.Err()
}

func (db *Ctx[k, v]) HExists(field k) (ok bool, err error) {

	var (
		cmd      *redis.BoolCmd
		fieldStr string
	)
	if fieldStr, err = db.toKeyStr(field); err != nil {
		return false, err
	}
	cmd = db.Rds.HExists(db.Ctx, db.Key, fieldStr)
	return cmd.Result()

}
func (db *Ctx[k, v]) HGetAll() (mapOut map[k]v, err error) {
	var (
		cmd *redis.MapStringStringCmd
		key k
		val v
	)
	mapOut = make(map[k]v)
	if cmd = db.Rds.HGetAll(db.Ctx, db.Key); cmd.Err() != nil {
		return mapOut, cmd.Err()
	}
	//append all data to mapOut
	for k, v := range cmd.Val() {
		if key, err = db.toKey([]byte(k)); err != nil {
			log.Info().AnErr("HGetAll: key unmarshal error:", err)
			continue
		}
		if val, err = db.toValue([]byte(v)); err != nil {
			log.Info().AnErr("HGetAll: value unmarshal error:", err)
			continue
		}
		mapOut[key] = val
	}
	return mapOut, err
}

func (db *Ctx[k, v]) HMGET(fields ...k) (values []v, err error) {
	var (
		cmd          *redis.SliceCmd
		fieldsString []string
		rawValues    []string
	)
	if fieldsString, err = db.toKeyStrs(fields...); err != nil {
		return nil, err
	}
	if cmd = db.Rds.HMGet(db.Ctx, db.Key, fieldsString...); cmd.Err() != nil {
		return nil, cmd.Err()
	}
	rawValues = make([]string, len(cmd.Val()))
	for i, val := range cmd.Val() {
		if val == nil {
			continue
		}
		rawValues[i] = val.(string)
	}
	return db.toValues(rawValues...)
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
	}
	cmd = db.Rds.HDel(db.Ctx, db.Key, fieldStrs...)
	return cmd.Err()
}
func (db *Ctx[k, v]) HKeys() (fields []k, err error) {
	var (
		cmd *redis.StringSliceCmd
	)
	if cmd = db.Rds.HKeys(db.Ctx, db.Key); cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return db.toKeys(cmd.Val())
}
func (db *Ctx[k, v]) HVals() (values []v, err error) {
	var cmd *redis.StringSliceCmd
	if cmd = db.Rds.HVals(db.Ctx, db.Key); cmd.Err() != nil {
		return nil, cmd.Err()
	}
	return db.toValues(cmd.Val()...)
}
func (db *Ctx[k, v]) HIncrBy(field k, increment int64) (err error) {
	var (
		cmd      *redis.IntCmd
		fieldStr string
	)
	if fieldStr, err = db.toKeyStr(field); err != nil {
		return err
	}
	cmd = db.Rds.HIncrBy(db.Ctx, db.Key, fieldStr, increment)
	return cmd.Err()
}

func (db *Ctx[k, v]) HIncrByFloat(field k, increment float64) (err error) {
	var (
		cmd      *redis.FloatCmd
		fieldStr string
	)
	if fieldStr, err = db.toKeyStr(field); err != nil {
		return err
	}
	cmd = db.Rds.HIncrByFloat(db.Ctx, db.Key, fieldStr, increment)
	return cmd.Err()

}
func (db *Ctx[k, v]) HSetNX(field k, value v) (err error) {
	var (
		cmd      *redis.BoolCmd
		fieldStr string
		valStr   string
	)
	if fieldStr, err = db.toKeyStr(field); err != nil {
		return err
	}
	if valStr, err = db.toValueStr(value); err != nil {
		return err
	}
	cmd = db.Rds.HSetNX(db.Ctx, db.Key, fieldStr, valStr)
	return cmd.Err()
}
