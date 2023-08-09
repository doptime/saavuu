package data

import (
	"reflect"

	"github.com/rs/zerolog/log"

	"github.com/vmihailenco/msgpack/v5"
)

// get all keys that match the pattern, and return a map of key->value
func (db *Ctx[k, v]) GetAll(match string, mapOut interface{}) (err error) {
	var (
		keys []string = []string{match}
		val  []byte
	)
	mapElem := reflect.TypeOf(mapOut)
	if (mapElem.Kind() != reflect.Map) || (mapElem.Key().Kind() != reflect.String) {
		log.Fatal().Msg("mapOut must be a map[string] struct/interface{}")
	}
	if keys, err = db.Scan(match, 0, 1024*1024*1024); err != nil {
		return err
	}
	var result error
	structSupposed := mapElem.Elem()
	for _, key := range keys {
		if val, result = db.Rds.Get(db.Ctx, key).Bytes(); result != nil {
			err = result
			continue
		}
		obj := reflect.New(structSupposed).Interface()
		if msgpack.Unmarshal(val, obj) != nil {
			err = result
			continue
		} else {
			reflect.ValueOf(mapOut).SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(obj).Elem())
		}
	}
	return result
}

// set each key value of _map to redis string type key value
func (db *Ctx[k, v]) SetAll(_map interface{}) (err error) {
	mapElem := reflect.TypeOf(_map)
	if (mapElem.Kind() != reflect.Map) || (mapElem.Key().Kind() != reflect.String) {
		log.Fatal().Msg("mapOut must be a map[string] struct/interface{}")
	}
	//HSet each element of _map to redis
	var result error
	pipe := db.Rds.Pipeline()
	for _, k := range reflect.ValueOf(_map).MapKeys() {
		v := reflect.ValueOf(_map).MapIndex(k)
		if bytes, err := msgpack.Marshal(v.Interface()); err != nil {
			result = err
		} else {
			pipe.Set(db.Ctx, k.String(), bytes, -1)
		}
	}
	pipe.Exec(db.Ctx)
	return result
}
