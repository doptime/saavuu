package rCtx

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/logger"
)

func (dc *DataCtx) HGet(key string, field string, param interface{}) (err error) {
	//use reflect to check if param is a pointer
	if reflect.TypeOf(param).Kind() != reflect.Ptr {
		logger.Lshortfile.Fatal("param must be a pointer")
	}

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

func (dc *DataCtx) HGetAll(key string, stru interface{}) (param map[string]interface{}, err error) {
	cmd := dc.Rds.HGetAll(dc.Ctx, key)
	data, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	param = make(map[string]interface{})
	for k, v := range data {
		//make a copy of stru , to obj
		obj := reflect.New(reflect.TypeOf(stru).Elem()).Interface()
		if err = msgpack.Unmarshal([]byte(v), &obj); err == nil {
			param[k] = obj
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

// golang version of python scan_iter
func (dc *DataCtx) Scan(match string, cursor uint64, count int64) (keys []string, err error) {
	keys, _, err = dc.Rds.Scan(dc.Ctx, cursor, match, count).Result()
	return keys, err
}

// update redis value schema, the value should be a pointer to  struct of msgpack
// match is the key pattern, if match end with *, scan all keys start with match
// example: rename field "V" to "VAL"
//
//	demoStruct struct {
//		//use client's time
//		Val           int64 `msgpack:"V,alias:VAL"`
//	}
func (dc *DataCtx) UpdateSchema(match string, dataStruct interface{}) (err error) {
	var (
		val  []byte
		keys []string = []string{match}
		data map[string]string
	)
	//error check, error if reflect of dataStruct is not a pointer
	if reflect.TypeOf(dataStruct).Kind() != reflect.Ptr {
		return fmt.Errorf("dataStruct must be a pointer")
	}
	//if keyStart end with *,iter scan all key start with keyStart
	if strings.HasSuffix(match, "*") {
		if keys, err = dc.Scan(match, 0, 1024*1024*1024); err != nil {
			return err
		}
	}
	//check type of redis value
	for _, key := range keys {
		if dc.Rds.Type(dc.Ctx, key).Val() == "hash" {
			cmd := dc.Rds.HGetAll(dc.Ctx, key)
			if data, err = cmd.Result(); err != nil {
				return err
			}
			pipe := dc.Rds.Pipeline()
			for field, v := range data {
				if msgpack.Unmarshal([]byte(v), dataStruct) != nil {
					return err
				}
				bytes, err := msgpack.Marshal(dataStruct)
				if err != nil {
					return err
				}
				if status := pipe.HSet(dc.Ctx, key, field, bytes); status.Err() != nil {
					return status.Err()
				}
			}
			if _, err := pipe.Exec(dc.Ctx); err != nil {
				return err
			}
		} else if dc.Rds.Type(dc.Ctx, key).Val() == "string" {
			if val, err = dc.Rds.Get(dc.Ctx, key).Bytes(); err != nil {
				return err
			}
			if msgpack.Unmarshal(val, dataStruct) != nil {
				return err
			}
			bytes, err := msgpack.Marshal(dataStruct)
			if err != nil {
				return err
			}
			if err = dc.Rds.Set(dc.Ctx, key, bytes, 0).Err(); err != nil {
				return err
			}
		} else if dc.Rds.Type(dc.Ctx, key).Val() == "list" {
			//not impleted yet
			return fmt.Errorf("not impleted yet")

		} else if dc.Rds.Type(dc.Ctx, key).Val() == "set" {
			//not impleted yet
			return fmt.Errorf("not impleted yet")
		} else if dc.Rds.Type(dc.Ctx, key).Val() == "zset" {
			//not impleted yet
			return fmt.Errorf("not impleted yet")
		} else {
			return fmt.Errorf("unknown type")
		}
	}
	return nil
}

// update redis value schema, the value should be a pointer to  struct of msgpack
// match is the key pattern, if match end with *, scan all keys start with match
// demo :
// dc := saavuu.NewDataContext(context.Background())
// var meditEpisode *MeditationEpisode = &MeditationEpisode{}
// StructureOldToNew := func(old interface{}) interface{} {

//		type Bar struct {
//			Duration uint32 `msgpack:"D"`
//		}
//		var oldMeditEpisode *MeditationEpisode = old.(*MeditationEpisode)
//		_newStruct := Bar{
//			Duration:     uint32(oldMeditEpisode.Duration) * 1000,
//		}
//		return _newStruct
//	}
//
// dc.UpdateSchemaViaFunc("TrajMedit:ekebmgfi24g6:*", meditEpisode, StructureOldToNew)
func (dc *DataCtx) UpdateSchemaViaFunc(match string, dataStruct interface{}, StructOldToNew func(interface{}) interface{}) (err error) {
	var (
		val  []byte
		keys []string = []string{match}
		data map[string]string
	)
	//error check, error if reflect of dataStruct is not a pointer
	if reflect.TypeOf(dataStruct).Kind() != reflect.Ptr {
		return fmt.Errorf("dataStruct must be a pointer")
	}
	//if keyStart end with *,iter scan all key start with keyStart
	if strings.HasSuffix(match, "*") {
		if keys, err = dc.Scan(match, 0, 1024*1024*1024); err != nil {
			return err
		}
	}
	//check type of redis value
	for _, key := range keys {
		if dc.Rds.Type(dc.Ctx, key).Val() == "hash" {
			cmd := dc.Rds.HGetAll(dc.Ctx, key)
			if data, err = cmd.Result(); err != nil {
				return err
			}
			pipe := dc.Rds.Pipeline()
			for field, v := range data {
				if msgpack.Unmarshal([]byte(v), dataStruct) != nil {
					return err
				}
				newStruct := StructOldToNew(dataStruct)
				bytes, err := msgpack.Marshal(newStruct)
				if err != nil {
					return err
				}
				if status := pipe.HSet(dc.Ctx, key, field, bytes); status.Err() != nil {
					return status.Err()
				}
			}
			if _, err := pipe.Exec(dc.Ctx); err != nil {
				return err
			}
		} else if dc.Rds.Type(dc.Ctx, key).Val() == "string" {
			if val, err = dc.Rds.Get(dc.Ctx, key).Bytes(); err != nil {
				return err
			}
			if msgpack.Unmarshal(val, dataStruct) != nil {
				return err
			}
			bytes, err := msgpack.Marshal(dataStruct)
			if err != nil {
				return err
			}
			if err = dc.Rds.Set(dc.Ctx, key, bytes, 0).Err(); err != nil {
				return err
			}
		} else if dc.Rds.Type(dc.Ctx, key).Val() == "list" {
			//not impleted yet
			return fmt.Errorf("not impleted yet")

		} else if dc.Rds.Type(dc.Ctx, key).Val() == "set" {
			//not impleted yet
			return fmt.Errorf("not impleted yet")
		} else if dc.Rds.Type(dc.Ctx, key).Val() == "zset" {
			//not impleted yet
			return fmt.Errorf("not impleted yet")
		} else {
			return fmt.Errorf("unknown type")
		}
	}
	return nil
}

// iterate redis value , send each value to DataProcess,
// DataProcess is a function with 3 parameters, key, field, dataStruct
// the dataStruct should be a pointer to  struct
func (dc *DataCtx) DataIterator(match string, dataStruct interface{}, DataProcess func(string, string, interface{})) (err error) {
	var (
		val  []byte
		keys []string = []string{match}
		data map[string]string
	)
	//error check, error if reflect of dataStruct is not a pointer
	if reflect.TypeOf(dataStruct).Kind() != reflect.Ptr {
		return fmt.Errorf("dataStruct must be a pointer")
	}
	//if keyStart end with *,iter scan all key start with keyStart
	if strings.HasSuffix(match, "*") {
		if keys, err = dc.Scan(match, 0, 1024*1024*1024); err != nil {
			return err
		}
	}
	//check type of redis value
	for _, key := range keys {
		if dc.Rds.Type(dc.Ctx, key).Val() == "hash" {
			cmd := dc.Rds.HGetAll(dc.Ctx, key)
			if data, err = cmd.Result(); err != nil {
				return err
			}
			for field, v := range data {
				if msgpack.Unmarshal([]byte(v), dataStruct) != nil {
					return err
				}
				DataProcess(key, field, dataStruct)
			}
		} else if dc.Rds.Type(dc.Ctx, key).Val() == "string" {
			if val, err = dc.Rds.Get(dc.Ctx, key).Bytes(); err != nil {
				return err
			}
			if msgpack.Unmarshal(val, dataStruct) != nil {
				return err
			}
			DataProcess(key, "", dataStruct)
		} else if dc.Rds.Type(dc.Ctx, key).Val() == "list" {
			//not impleted yet
			return fmt.Errorf("not impleted yet")

		}
	}
	return nil
}
