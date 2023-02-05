package data

import (
	"errors"
	"reflect"

	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/logger"
	"github.com/yangkequn/saavuu/rds"
)

func (dc *DataCtx) HGet(key string, field string, param interface{}) (err error) {
	//use reflect to check if param is a pointer
	if reflect.TypeOf(param).Kind() != reflect.Ptr {
		logger.Lshortfile.Println("param must be a pointer")
		return errors.New("param must be a pointer")
	}

	cmd := dc.Rds.HGet(dc.Ctx, key, field)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, param)
}

func (dc *DataCtx) HSet(key string, field string, param interface{}) (err error) {
	return rds.HSet(dc.Ctx, dc.Rds, key, field, param)
}

func (dc *DataCtx) HExists(key string, field string) (ok bool, err error) {
	cmd := dc.Rds.HExists(dc.Ctx, key, field)
	return cmd.Val(), cmd.Err()
}
func (dc *DataCtx) HGetAll(key string, mapOut interface{}) (err error) {
	return rds.HGetAll(dc.Ctx, dc.Rds, key, mapOut)
}
func (dc *DataCtx) HSetAll(key string, _map interface{}) (err error) {
	mapElem := reflect.TypeOf(_map)
	if (mapElem.Kind() != reflect.Map) || (mapElem.Key().Kind() != reflect.String) {
		logger.Lshortfile.Println("mapOut must be a map[string] struct/interface{}")
		return errors.New("mapOut must be a map[string] struct/interface{}")
	}
	//HSet each element of _map to redis
	var result error
	pipe := dc.Rds.Pipeline()
	for _, k := range reflect.ValueOf(_map).MapKeys() {
		v := reflect.ValueOf(_map).MapIndex(k)
		if bytes, err := msgpack.Marshal(v.Interface()); err != nil {
			result = err
		} else {
			pipe.HSet(dc.Ctx, key, k.String(), bytes)
		}
	}
	pipe.Exec(dc.Ctx)
	return result
}
func (dc *DataCtx) HMGET(key string, _map interface{}, fields ...string) (err error) {
	mapElem := reflect.TypeOf(_map)
	if (mapElem.Kind() != reflect.Map) || (mapElem.Key().Kind() != reflect.String) {
		logger.Lshortfile.Println("mapOut must be a map[string] struct/interface{}")
		return errors.New("mapOut must be a map[string] struct/interface{}")
	}
	structSupposed := mapElem.Elem()
	cmd := dc.Rds.HMGet(dc.Ctx, key, fields...)
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

func (dc *DataCtx) HGetAllDefault(key string) (param map[string]interface{}, err error) {
	cmd := dc.Rds.HGetAll(dc.Ctx, key)
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
func (dc *DataCtx) HLen(key string) (length int64, err error) {
	cmd := dc.Rds.HLen(dc.Ctx, key)
	return cmd.Val(), cmd.Err()
}
func (dc *DataCtx) HDel(key string, field string) (err error) {
	status := dc.Rds.HDel(dc.Ctx, key, field)
	return status.Err()
}
func (dc *DataCtx) HKeys(key string) (fields []string, err error) {
	return rds.HKeys(dc.Ctx, dc.Rds, key)
}
func (dc *DataCtx) HVals(key string) (values []interface{}, err error) {
	cmd := dc.Rds.HVals(dc.Ctx, key)
	data := cmd.Val()
	//unmarshal each value of cmd.Val() to interface{}, using msgpack
	for _, v := range data {
		var obj interface{}
		if err = msgpack.Unmarshal([]byte(v), &obj); err == nil {
			values = append(values, obj)
		}
	}
	return values, nil
}
func (dc *DataCtx) HIncrBy(key string, field string, increment int64) (err error) {
	status := dc.Rds.HIncrBy(dc.Ctx, key, field, increment)
	return status.Err()
}
func (dc *DataCtx) HIncrByFloat(key string, field string, increment float64) (err error) {
	status := dc.Rds.HIncrByFloat(dc.Ctx, key, field, increment)
	return status.Err()
}
func (dc *DataCtx) HSetNX(key string, field string, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := dc.Rds.HSetNX(dc.Ctx, key, field, bytes)
	return status.Err()
}

// golang version of python scan_iter
func (dc *DataCtx) Scan(match string, cursor uint64, count int64) (keys []string, err error) {
	keys, _, err = dc.Rds.Scan(dc.Ctx, cursor, match, count).Result()
	return keys, err
}
