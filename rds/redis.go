package rds

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/logger"
)

func HSet(ctx context.Context, rc *redis.Client, key string, field string, value interface{}) (err error) {
	bytes, err := msgpack.Marshal(value)
	if err != nil {
		return err
	}
	status := rc.HSet(ctx, key, field, bytes)
	return status.Err()
}

func Get(ctx context.Context, rc *redis.Client, key string, param interface{}) (err error) {
	cmd := rc.Get(ctx, key)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, param)
}
func Set(ctx context.Context, rc *redis.Client, key string, param interface{}, expiration time.Duration) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := rc.Set(ctx, key, bytes, expiration)
	return status.Err()
}

func HGetAll(ctx context.Context, rc *redis.Client, key string, mapOut interface{}) (err error) {
	mapElem := reflect.TypeOf(mapOut)
	if (mapElem.Kind() != reflect.Map) || (mapElem.Key().Kind() != reflect.String) {
		logger.Lshortfile.Println("mapOut must be a map[string] struct/interface{}")
		return errors.New("mapOut must be a map[string] struct/interface{}")
	}
	cmd := rc.HGetAll(ctx, key)
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

func HKeys(ctx context.Context, rc *redis.Client, key string) (fields []string, err error) {
	cmd := rc.HKeys(ctx, key)
	return cmd.Val(), cmd.Err()
}

func RPush(ctx context.Context, rc *redis.Client, key string, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := rc.RPush(ctx, key, bytes)
	return status.Err()
}
func LSet(ctx context.Context, rc *redis.Client, key string, index int64, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := rc.LSet(ctx, key, index, bytes)
	return status.Err()
}
func LGet(ctx context.Context, rc *redis.Client, key string, index int64, param interface{}) (err error) {
	cmd := rc.LIndex(ctx, key, index)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, param)
}
func LLen(ctx context.Context, rc *redis.Client, key string) (length int64, err error) {
	cmd := rc.LLen(ctx, key)
	return cmd.Result()
}
func SAdd(ctx context.Context, rc *redis.Client, key string, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := rc.SAdd(ctx, key, bytes)
	return status.Err()
}
func SRem(ctx context.Context, rc *redis.Client, key string, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := rc.SRem(ctx, key, bytes)
	return status.Err()
}
func SIsMember(ctx context.Context, rc *redis.Client, key string, param interface{}) (isMember bool, err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return false, err
	}
	cmd := rc.SIsMember(ctx, key, bytes)
	return cmd.Result()
}
func SMembers(ctx context.Context, rc *redis.Client, key string) (members []string, err error) {
	cmd := rc.SMembers(ctx, key)
	return cmd.Result()
}
func Time(ctx context.Context, rc *redis.Client) (time time.Time, err error) {
	cmd := rc.Time(ctx)
	return cmd.Result()
}
