package redisContext

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v5"
)

type DataCtx struct {
	Ctx context.Context
	Rds *redis.Client
}

func (dc *DataCtx) Get(key string, param interface{}) (err error) {
	cmd := dc.Rds.Get(dc.Ctx, key)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, param)
}
func (dc *DataCtx) Set(key string, param interface{}, expiration time.Duration) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := dc.Rds.Set(dc.Ctx, key, bytes, expiration)
	return status.Err()
}
func (dc *DataCtx) HGet(key string, field string, param interface{}) (err error) {
	cmd := dc.Rds.HGet(dc.Ctx, key, field)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, param)
}
func (dc *DataCtx) HSet(key string, field string, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := dc.Rds.HSet(dc.Ctx, key, field, bytes)
	return status.Err()
}
func (dc *DataCtx) HExists(key string, field string) (ok bool) {
	cmd := dc.Rds.HExists(dc.Ctx, key, field)
	return cmd.Val()
}

func (dc *DataCtx) RPush(key string, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := dc.Rds.RPush(dc.Ctx, key, bytes)
	return status.Err()
}
func (dc *DataCtx) LSet(key string, index int64, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := dc.Rds.LSet(dc.Ctx, key, index, bytes)
	return status.Err()
}
func (dc *DataCtx) LGet(key string, index int64, param interface{}) (err error) {
	cmd := dc.Rds.LIndex(dc.Ctx, key, index)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, param)
}
func (dc *DataCtx) LLen(key string) (length int64) {
	cmd := dc.Rds.LLen(dc.Ctx, key)
	return cmd.Val()
}

// append to Set
func (dc *DataCtx) SAdd(key string, members ...interface{}) (err error) {
	status := dc.Rds.SAdd(dc.Ctx, key, members)
	return status.Err()
}
func (dc *DataCtx) SRem(key string, members ...interface{}) (err error) {
	status := dc.Rds.SRem(dc.Ctx, key, members)
	return status.Err()
}
func (dc *DataCtx) SIsMember(key string, param interface{}) (ok bool) {
	cmd := dc.Rds.SIsMember(dc.Ctx, key, param)
	return cmd.Val()
}
func (dc *DataCtx) SMembers(key string, param interface{}) (members []string, err error) {
	cmd := dc.Rds.SMembers(dc.Ctx, key)
	members, err = cmd.Result()
	if err != nil {
		return nil, err
	}
	return members, nil
}
func (dc *DataCtx) HGetAll(key string, decodeFun func(string) (obj interface{}, erro error)) (param map[string]interface{}, err error) {
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
		if Decoded, err := decodeFun(v); err == nil {
			param[k] = Decoded
		} else {
			return nil, err
		}
	}
	return param, nil
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
