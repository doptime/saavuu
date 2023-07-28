package rds

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/msgpack/v5"
)

var ErrInvalidField = errors.New("invalid field")

func HGet(ctx context.Context, rc *redis.Client, key string, field interface{}, value interface{}) (err error) {
	var (
		cmd              *redis.StringCmd
		fieldBytes, data []byte
	)
	if field == nil {
		return ErrInvalidField
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
		return ErrInvalidField
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

func HExists(ctx context.Context, rc *redis.Client, key string, field interface{}) (ok bool, err error) {
	var (
		cmd      *redis.BoolCmd
		fieldStr string
	)
	if field == nil {
		return false, ErrInvalidField
	}
	if fieldStr, ok = field.(string); ok {
		cmd = rc.HExists(ctx, key, fieldStr)
	} else if fieldBytes, err := json.Marshal(field); err != nil {
		return false, err
	} else {
		cmd = rc.HExists(ctx, key, string(fieldBytes))
	}
	return cmd.Result()
}
func HDel(ctx context.Context, rc *redis.Client, key string, field interface{}) (err error) {
	var (
		cmd      *redis.IntCmd
		fieldStr string
		ok       bool
	)
	if field == nil {
		return ErrInvalidField
	}
	if fieldStr, ok = field.(string); ok {
		cmd = rc.HDel(ctx, key, fieldStr)
	} else if fieldBytes, err := json.Marshal(field); err != nil {
		return err
	} else {
		cmd = rc.HDel(ctx, key, string(fieldBytes))
	}
	return cmd.Err()
}

func HGetAll(ctx context.Context, rc *redis.Client, key string, mapOut interface{}) (err error) {
	var (
		cmd *redis.MapStringStringCmd
	)
	mapElem := reflect.TypeOf(mapOut)
	//if mapOut is  a pointer to  map , such as: var mapOut *map[uint32]interface{}
	if mapElem.Kind() == reflect.Ptr {
		if mapElem.Elem().Kind() != reflect.Map {
			log.Info().Msg("mapOut must be a map[interface{}] interface{} or *map[interface{}] interface{}")
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
		log.Info().Msg("mapOut must be a map[interface{}] interface{}")
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
			log.Info().AnErr("HGetAll: value unmarshal error:", err)
			continue
		}
		//case key is  string, no need to unmarshal key
		if isKeyString {
			reflect.ValueOf(mapOut).SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(valObj).Elem())
			continue
		}
		// case key is not string, unmarshal key
		keyObj := reflect.New(KeyStructSupposed).Interface()
		if err = json.Unmarshal([]byte(k), &keyObj); err == nil {
			reflect.ValueOf(mapOut).SetMapIndex(reflect.ValueOf(keyObj).Elem(), reflect.ValueOf(valObj).Elem())
		} else if KeyStructSupposed.Kind() == reflect.Interface {
			//if KeyStructSupposed is interface{}, save key as string
			reflect.ValueOf(mapOut).SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(valObj).Elem())
		} else {
			log.Info().AnErr("HGetAll: key unmarshal error:", err)
		}
	}
	return err
}
func HSetAll(ctx context.Context, rc *redis.Client, key string, mapIn interface{}) (err error) {
	mapElem := reflect.TypeOf(mapIn)
	if mapElem.Kind() != reflect.Map {
		log.Info().Msg("mapIn must be a map[interface{}] struct/interface{}")
		return errors.New("mapIn must be a map[interface{}] struct/interface{}")
	}
	mapOut := make(map[string]interface{})
	//append all key value of mapIn to mapOut, using msgpack
	for _, k := range reflect.ValueOf(mapIn).MapKeys() {
		//marshal key to bytes
		keyBytes, err := json.Marshal(k.Interface())
		if err != nil {
			log.Info().AnErr("HSetMap: key marshal error:", err)
			continue
		}
		//marshal value to bytes
		valueBytes, err := msgpack.Marshal(reflect.ValueOf(mapIn).MapIndex(k).Interface())
		if err != nil {
			log.Info().AnErr("HSetMap: value marshal error:", err)
			continue
		}
		mapOut[string(keyBytes)] = valueBytes
	}
	//hset to redis
	return rc.HSet(ctx, key, mapOut).Err()
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
		log.Info().Msg("fields must be a pointer to slice")
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
			log.Info().AnErr("HKeys1: field unmarshal error:", err)
			continue
		}
		//*fields = append(*fields, field)
		reflect.ValueOf(fields).Elem().Set(reflect.Append(reflect.ValueOf(fields).Elem(), reflect.ValueOf(field).Elem()))
	}
	return cmd.Err()
}

func HVals(ctx context.Context, rc *redis.Client, key string, values interface{}) (err error) {
	if !isPointerToSlice(values) {
		log.Info().Msg("values must be a pointer to slice")
		return errors.New("values must be a pointer to slice")
	}
	cmd := rc.HVals(ctx, key)
	// structFields := reflect.TypeOf(values).Elem()
	// *values = make([]interface{}, 0, len(cmd.Val()))
	structFields := reflect.TypeOf(values).Elem().Elem()
	slice := reflect.MakeSlice(reflect.TypeOf(values).Elem(), 0, len(cmd.Val()))
	reflect.ValueOf(values).Elem().Set(slice)
	for _, v := range cmd.Val() {
		value := reflect.New(structFields).Interface()
		if err = msgpack.Unmarshal([]byte(v), &value); err != nil {
			log.Info().AnErr("HVals: value unmarshal error:", err)
			continue
		}
		//*values = append(*values, value)
		reflect.ValueOf(values).Elem().Set(reflect.Append(reflect.ValueOf(values).Elem(), reflect.ValueOf(value).Elem()))
	}
	return cmd.Err()
}
func HIncrBy(ctx context.Context, rc *redis.Client, key string, field interface{}, incr int64) (err error) {
	var (
		cmd      *redis.IntCmd
		fieldStr string
		ok       bool
	)
	if field == nil {
		return ErrInvalidField
	}
	if fieldStr, ok = field.(string); ok {
		cmd = rc.HIncrBy(ctx, key, fieldStr, incr)
	} else if fieldBytes, err := json.Marshal(field); err != nil {
		return err
	} else {
		cmd = rc.HIncrBy(ctx, key, string(fieldBytes), incr)
	}
	return cmd.Err()
}
func HIncrByFloat(ctx context.Context, rc *redis.Client, key string, field interface{}, incr float64) (err error) {
	var (
		cmd      *redis.FloatCmd
		fieldStr string
		ok       bool
	)
	if field == nil {
		return ErrInvalidField
	}
	if fieldStr, ok = field.(string); ok {
		cmd = rc.HIncrByFloat(ctx, key, fieldStr, incr)
	} else if fieldBytes, err := json.Marshal(field); err != nil {
		return err
	} else {
		cmd = rc.HIncrByFloat(ctx, key, string(fieldBytes), incr)
	}
	return cmd.Err()
}
func HSetNX(ctx context.Context, rc *redis.Client, key string, field interface{}, value interface{}) (err error) {
	var (
		cmd        *redis.BoolCmd
		fieldStr   string
		ok         bool
		valueBytes []byte
	)
	if field == nil {
		return ErrInvalidField
	}
	if valueBytes, err = msgpack.Marshal(value); err != nil {
		return err
	}
	if fieldStr, ok = field.(string); ok {
		cmd = rc.HSetNX(ctx, key, fieldStr, value)
	} else if fieldBytes, err := json.Marshal(field); err != nil {
		return err
	} else {
		cmd = rc.HSetNX(ctx, key, string(fieldBytes), valueBytes)
	}
	return cmd.Err()
}
