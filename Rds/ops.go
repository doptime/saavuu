package rds

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/logger"
)

func HSet(ctx context.Context, rds *redis.Client, key string, field string, value interface{}) (err error) {
	bytes, err := msgpack.Marshal(value)
	if err != nil {
		return err
	}
	status := rds.HSet(ctx, key, field, bytes)
	return status.Err()
}

func Get(ctx context.Context, rds *redis.Client, key string, param interface{}) (err error) {
	cmd := rds.Get(ctx, key)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, param)
}
func Set(ctx context.Context, rds *redis.Client, key string, param interface{}, expiration time.Duration) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := rds.Set(ctx, key, bytes, expiration)
	return status.Err()
}

func HGetAll(ctx context.Context, rds *redis.Client, key string, mapOut interface{}) (err error) {
	mapElem := reflect.TypeOf(mapOut)
	if (mapElem.Kind() != reflect.Map) || (mapElem.Key().Kind() != reflect.String) {
		logger.Lshortfile.Println("mapOut must be a map[string] struct/interface{}")
		return errors.New("mapOut must be a map[string] struct/interface{}")
	}
	cmd := rds.HGetAll(ctx, key)
	data, err := cmd.Result()
	if err != nil {
		return err
	}
	//append all data to mapOut
	structSupposed := mapElem.Elem()
	for k, v := range data {
		//make a copy of stru , to obj
		obj := reflect.New(structSupposed).Interface()
		if err = msgpack.Unmarshal([]byte(v), &obj); err == nil {
			reflect.ValueOf(mapOut).SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(obj).Elem())
		}
	}
	return err
}

func HKeys(ctx context.Context, rds *redis.Client, key string) (fields []string, err error) {
	cmd := rds.HKeys(ctx, key)
	return cmd.Val(), cmd.Err()
}
