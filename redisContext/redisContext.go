package redisContext

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v5"
)

type RedisContext struct {
	Ctx      context.Context
	DataRds  *redis.Client
	ParamRds *redis.Client
}

//RedisCall: 1.use RPush to push data to redis. 2.use BLPop to pop data from selected channel
//return: error
func (sc RedisContext) RdsApiBasic(ServiceKey string, paramIn map[string]interface{}) (result []byte, err error) {
	var (
		b       []byte
		BackTo  string = fmt.Sprintf("%x", rand.Int63())
		results []string
	)
	paramIn["BackTo"] = BackTo

	if b, err = msgpack.Marshal(paramIn); err != nil {
		return nil, err
	}
	ppl := sc.ParamRds.Pipeline()
	ppl.RPush(sc.Ctx, ServiceKey, b)
	//长期不执行的任务，抛弃
	ppl.Expire(sc.Ctx, ServiceKey, time.Second*60)
	if _, err := ppl.Exec(sc.Ctx); err != nil {
		return nil, err
	}
	//BLPop 返回结果 [key1,value1,key2,value2]
	if results, err = sc.ParamRds.BLPop(sc.Ctx, time.Second*20, BackTo).Result(); err != nil {
		return nil, err
	}
	return []byte(results[1]), nil
}
func (sc RedisContext) RdsApi(ServiceKey string, structIn interface{}, out interface{}) (err error) {
	var (
		paramIn = map[string]interface{}{}
		ok      bool
	)
	//convert struct to map, to allow add field "BackTo"
	if paramIn, ok = structIn.(map[string]interface{}); !ok {
		b, err := msgpack.Marshal(structIn)
		if err != nil {
			return err
		}
		if err = msgpack.Unmarshal(b, &paramIn); err != nil {
			return err
		}
	}

	result, err := sc.RdsApiBasic(ServiceKey, paramIn)
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(result, out)
}

func (sc RedisContext) Get(key string, param interface{}) (err error) {
	cmd := sc.DataRds.Get(sc.Ctx, key)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, param)
}
func (sc RedisContext) Set(key string, param interface{}, expiration time.Duration) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := sc.DataRds.Set(sc.Ctx, key, bytes, expiration)
	return status.Err()
}
func (sc RedisContext) HGet(key string, field string, param interface{}) (err error) {
	cmd := sc.DataRds.HGet(sc.Ctx, key, field)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, param)
}
func (sc RedisContext) HSet(key string, field string, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := sc.DataRds.HSet(sc.Ctx, key, field, bytes)
	return status.Err()
}
func (sc RedisContext) RPush(key string, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := sc.DataRds.RPush(sc.Ctx, key, bytes)
	return status.Err()
}
func (sc RedisContext) LSet(key string, index int64, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := sc.DataRds.LSet(sc.Ctx, key, index, bytes)
	return status.Err()
}
func (sc RedisContext) LGet(key string, index int64, param interface{}) (err error) {
	cmd := sc.DataRds.LIndex(sc.Ctx, key, index)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, param)
}
func (sc RedisContext) LLen(key string) (length int64) {
	cmd := sc.DataRds.LLen(sc.Ctx, key)
	return cmd.Val()
}
