package saavuu

import (
	"time"

	"github.com/vmihailenco/msgpack"
)

func (scvCtx *ServiceContext) RedisGet(key string, param interface{}) (err error) {
	cmd := Config.rds.Get(scvCtx.ctx, key)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, param)
}
func (scvCtx *ServiceContext) RedisSet(key string, param interface{}, expiration time.Duration) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := Config.rds.Set(scvCtx.ctx, key, bytes, expiration)
	return status.Err()
}
func (scvCtx *ServiceContext) RedisHGet(key string, field string, param interface{}) (err error) {
	cmd := Config.rds.HGet(scvCtx.ctx, key, field)
	data, err := cmd.Bytes()
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, param)
}

func (scvCtx *ServiceContext) RedisHSet(key string, field string, param interface{}) (err error) {
	bytes, err := msgpack.Marshal(param)
	if err != nil {
		return err
	}
	status := Config.rds.HSet(scvCtx.ctx, key, field, bytes)
	return status.Err()
}
