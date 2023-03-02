package rds

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/logger"
)

func HGet(ctx context.Context, rc *redis.Client, key string, field interface{}, value interface{}) (err error) {
	var (
		cmd              *redis.StringCmd
		fieldBytes, data []byte
	)
	if field == nil {
		return errors.New("field is nil")
	}
	//case string do not need to marshal
	if _field, ok := field.(string); ok {
		cmd = rc.HGet(ctx, key, _field)
	} else if fieldBytes, err = json.Marshal(field); err != nil {
		//case fail to marshal
		return err
	} else {
		//case marshal success
		cmd = rc.HGet(ctx, key, string(fieldBytes))
	}

	if data, err = cmd.Bytes(); err != nil {
		return err
	}
	return msgpack.Unmarshal(data, value)
}

func HSet(ctx context.Context, rc *redis.Client, key string, field interface{}, value interface{}) (err error) {
	var (
		fieldBytes []byte
		valueBytes []byte
		status     *redis.IntCmd
	)
	if field == nil {
		return errors.New("field is nil")
	}

	if valueBytes, err = msgpack.Marshal(value); err != nil {
		return err
	}

	if _, ok := field.(string); ok {
		status = rc.HSet(ctx, key, field, valueBytes)
	} else if fieldBytes, err = json.Marshal(field); err != nil {
		return err
	} else {
		status = rc.HSet(ctx, key, fieldBytes, valueBytes)

	}
	return status.Err()
}

func HGetAll(ctx context.Context, rc *redis.Client, key string, mapOut interface{}) (err error) {
	var (
		cmd *redis.MapStringStringCmd
	)
	mapElem := reflect.TypeOf(mapOut)
	//if mapOut is  a pointer to  map , such as: var mapOut *map[uint32]interface{}
	if mapElem.Kind() == reflect.Ptr {
		if mapElem.Elem().Kind() != reflect.Map {
			logger.Lshortfile.Println("mapOut must be a map[interface{}] interface{} or *map[interface{}] interface{}")
			return errors.New("mapOut must be a map[interface{}] interface{} or *map[interface{}] interface{}")
		}
		//if mapOut is a pointer to nil map, make a new one
		if reflect.ValueOf(mapOut).Elem().IsNil() {
			reflect.ValueOf(mapOut).Elem().Set(reflect.MakeMap(mapElem.Elem()))
		}
		//make sure mapElem is a map
		mapElem = mapElem.Elem()
		mapOut = reflect.ValueOf(mapOut).Elem().Interface()
	}
	//make sure mapElem is a map
	if mapElem.Kind() != reflect.Map {
		logger.Lshortfile.Println("mapOut must be a map[interface{}] interface{}")
		return errors.New("mapOut must be a map[interface{}] interface{}")
	}
	if cmd = rc.HGetAll(ctx, key); cmd.Err() != nil {
		return cmd.Err()
	}

	//append all data to mapOut
	KeyStructSupposed := mapElem.Key()
	isKeyString := KeyStructSupposed.Kind() == reflect.String
	valueStructSupposed := mapElem.Elem()
	for k, v := range cmd.Val() {
		//make a copy of stru , to valObj
		valObj := reflect.New(valueStructSupposed).Interface()
		if err = msgpack.Unmarshal([]byte(v), &valObj); err != nil {
			logger.Lshortfile.Println("HGetAll: value unmarshal error:", err)
			continue
		}
		//case key is  string, no need to unmarshal key
		if isKeyString {
			reflect.ValueOf(mapOut).SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(valObj).Elem())
			continue
		}
		// case key is not string, unmarshal key
		keyObj := reflect.New(KeyStructSupposed).Interface()
		if err = json.Unmarshal([]byte(k), &keyObj); err != nil {
			logger.Lshortfile.Println("HGetAll: key unmarshal error:", err)
			continue
		}
		reflect.ValueOf(mapOut).SetMapIndex(reflect.ValueOf(keyObj).Elem(), reflect.ValueOf(valObj).Elem())
	}
	return err
}
func HSetAll(ctx context.Context, rc *redis.Client, key string, mapIn interface{}) (err error) {
	mapElem := reflect.TypeOf(mapIn)
	if mapElem.Kind() != reflect.Map {
		logger.Lshortfile.Println("mapIn must be a map[interface{}] struct/interface{}")
		return errors.New("mapIn must be a map[interface{}] struct/interface{}")
	}
	mapOut := make(map[string]interface{})
	//append all key value of mapIn to mapOut, using msgpack
	for _, k := range reflect.ValueOf(mapIn).MapKeys() {
		//marshal key to bytes
		keyBytes, err := json.Marshal(k.Interface())
		if err != nil {
			logger.Lshortfile.Println("HSetMap: key marshal error:", err)
			continue
		}
		//marshal value to bytes
		valueBytes, err := msgpack.Marshal(reflect.ValueOf(mapIn).MapIndex(k).Interface())
		if err != nil {
			logger.Lshortfile.Println("HSetMap: value marshal error:", err)
			continue
		}
		mapOut[string(keyBytes)] = valueBytes
	}
	//hset to redis
	return rc.HSet(ctx, key, mapOut).Err()
}
func HMGET(ctx context.Context, rc *redis.Client, key string, fields interface{}, mapOut interface{}) (err error) {
	var (
		cmd *redis.SliceCmd
	)
	//make sure fields should be a slice
	fieldsType := reflect.TypeOf(fields)
	if fieldsType.Kind() != reflect.Slice {
		logger.Lshortfile.Println("fields must be a slice")
		return errors.New("fields must be a slice")
	}
	fieldsElem := reflect.ValueOf(fields)
	//mapOut should be a map
	mapElem := reflect.TypeOf(mapOut)
	//if mapOut is  a pointer to  map , such as: var mapOut *map[uint32]interface{}
	if mapElem.Kind() == reflect.Ptr {
		if mapElem.Elem().Kind() != reflect.Map {
			logger.Lshortfile.Println("mapOut must be a map[interface{}] interface{} or *map[interface{}] interface{}")
			return errors.New("mapOut must be a map[interface{}] interface{} or *map[interface{}] interface{}")
		}
		//if mapOut is a pointer to nil map, make a new one
		//i.g. var mapOut map[uint32]interface{}   =>  var mapOut map[uint32]interface{} = make(map[uint32]interface{})
		if reflect.ValueOf(mapOut).Elem().IsNil() {
			reflect.ValueOf(mapOut).Elem().Set(reflect.MakeMap(mapElem.Elem()))
		}
		//make sure mapElem is a map
		mapElem = mapElem.Elem()
		mapOut = reflect.ValueOf(mapOut).Elem().Interface()
	}
	if mapElem.Kind() != reflect.Map {
		logger.Lshortfile.Println("mapOut must be a map[interface{}] interface{}")
		return errors.New("mapOut must be a map[interface{}] interface{}")
	}
	//if mapOut is nil, make a new one
	if reflect.ValueOf(mapOut).IsNil() {
		reflect.ValueOf(mapOut).Set(reflect.MakeMap(mapElem))
	}
	//if fieldsElem is not []string, marshal each field to string
	var fieldsString []string
	var isStringField bool
	if fieldsString, isStringField = fields.([]string); !isStringField {
		//marshal each field to string
		fieldsString = make([]string, 0, fieldsElem.Len())

		for i := 0; i < fieldsElem.Len(); i++ {
			b, err := json.Marshal(reflect.ValueOf(fields).Index(i).Interface())
			if err != nil {
				logger.Lshortfile.Println("HMGET: field marshal error:", err)
				continue
			}
			fieldsString = append(fieldsString, string(b))
		}
	}
	if cmd = rc.HMGet(ctx, key, fieldsString...); cmd.Err() != nil {
		return cmd.Err()
	}

	//append all data to mapOut
	KeyStructSupposed := mapElem.Key()
	valueStructSupposed := mapElem.Elem()

	//save all data to mapOut
	for i, v := range cmd.Val() {
		key := reflect.New(KeyStructSupposed).Interface()
		if err = json.Unmarshal([]byte(fieldsString[i]), &key); err != nil {
			logger.Lshortfile.Println("HMGET: key unmarshal error:", err)
			continue
		}
		if v == nil {
			//set _map with nil
			reflect.ValueOf(mapOut).SetMapIndex(reflect.ValueOf(key).Elem(), reflect.Zero(valueStructSupposed))
			continue
		}
		obj := reflect.New(valueStructSupposed).Interface()
		if err = msgpack.Unmarshal([]byte(v.(string)), &obj); err == nil {
			reflect.ValueOf(mapOut).SetMapIndex(reflect.ValueOf(key).Elem(), reflect.ValueOf(obj).Elem())
		}
	}
	return cmd.Err()
}
func isPointerToSlice(obj interface{}) (ok bool) {
	objType := reflect.TypeOf(obj)
	if objType.Kind() != reflect.Ptr {
		return false
	}
	if objType.Elem().Kind() != reflect.Slice {
		return false
	}
	return true
}

func HKeys(ctx context.Context, rc *redis.Client, key string, fields interface{}) (err error) {
	if !isPointerToSlice(fields) {
		logger.Lshortfile.Println("fields must be a pointer to slice")
		return errors.New("fields must be a pointer to slice")
	}
	cmd := rc.HKeys(ctx, key)
	//if fields if *[]string, return directly
	//not needed to unmarshal fields
	if reflect.TypeOf(fields).Elem().Elem().Kind() == reflect.String {
		reflect.ValueOf(fields).Elem().Set(reflect.ValueOf(cmd.Val()))
		return cmd.Err()
	}
	// structFields := reflect.TypeOf(fields).Elem()
	// *fields = make([]interface{}, 0, len(cmd.Val()))
	structFields := reflect.TypeOf(fields).Elem().Elem()
	slice := reflect.MakeSlice(reflect.TypeOf(fields).Elem(), 0, len(cmd.Val()))
	reflect.ValueOf(fields).Elem().Set(slice)
	for _, v := range cmd.Val() {
		field := reflect.New(structFields).Interface()
		if err = json.Unmarshal([]byte(v), &field); err != nil {
			logger.Lshortfile.Println("HKeys1: field unmarshal error:", err)
			continue
		}
		//*fields = append(*fields, field)
		reflect.ValueOf(fields).Elem().Set(reflect.Append(reflect.ValueOf(fields).Elem(), reflect.ValueOf(field).Elem()))
	}
	return cmd.Err()
}

func HValsPackFields(ctx context.Context, rc *redis.Client, key string, values *[]interface{}) (err error) {
	cmd := rc.HVals(ctx, key)
	data := cmd.Val()
	*values = make([]interface{}, 0, len(data))
	valueStruct := reflect.TypeOf(values).Elem().Elem()
	//unmarshal each value of cmd.Val() to interface{}, using msgpack
	for _, v := range data {
		obj := reflect.New(valueStruct).Interface()
		if err = msgpack.Unmarshal([]byte(v), &obj); err != nil {
			logger.Lshortfile.Println("HVals: value unmarshal error:", err)
			continue
		}
		*values = append(*values, obj)
	}
	return cmd.Err()
}
