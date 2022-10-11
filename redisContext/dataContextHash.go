package redisContext

import "github.com/vmihailenco/msgpack/v5"

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
func (dc *DataCtx) HLen(key string) (length int64) {
	cmd := dc.Rds.HLen(dc.Ctx, key)
	return cmd.Val()
}
func (dc *DataCtx) HDel(key string, field string) (err error) {
	status := dc.Rds.HDel(dc.Ctx, key, field)
	return status.Err()
}
func (dc *DataCtx) HKeys(key string) (fields []string) {
	cmd := dc.Rds.HKeys(dc.Ctx, key)
	return cmd.Val()
}
func (dc *DataCtx) HVals(key string) (values []string) {
	cmd := dc.Rds.HVals(dc.Ctx, key)
	return cmd.Val()
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
