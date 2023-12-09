package data

import (
	"github.com/rs/zerolog/log"

	"github.com/vmihailenco/msgpack/v5"
)

// get all keys that match the pattern, and return a map of key->value
func (db *Ctx[k, v]) GetAll(match string) (mapOut map[k]v, err error) {
	var (
		keys []string = []string{match}
		val  []byte
	)
	if keys, err = db.Scan(match, 0, 1024*1024*1024); err != nil {
		return nil, err
	}
	mapOut = make(map[k]v, len(keys))
	var result error
	for _, key := range keys {
		if val, result = db.Rds.Get(db.Ctx, key).Bytes(); result != nil {
			err = result
			continue
		}
		k, err := db.toKey([]byte(key))
		if err != nil {
			log.Info().AnErr("GetAll: key unmarshal error:", err).Msgf("Key: %s", db.Key)
			continue
		}
		v, err := db.toValue(val)
		if err != nil {
			log.Info().AnErr("GetAll: value unmarshal error:", err).Msgf("Key: %s", db.Key)
			continue
		}
		mapOut[k] = v
	}
	return mapOut, err
}

// set each key value of _map to redis string type key value
func (db *Ctx[k, v]) SetAll(_map map[k]v) (err error) {
	//HSet each element of _map to redis
	var (
		result error
		bytes  []byte
		keyStr string
	)
	pipe := db.Rds.Pipeline()
	for k, v := range _map {
		if bytes, err = msgpack.Marshal(v); err != nil {
			return err
		}
		if keyStr, err = db.toKeyStr(k); err != nil {
			return err
		}

		pipe.Set(db.Ctx, keyStr, bytes, -1)
	}
	pipe.Exec(db.Ctx)
	return result
}
