package data

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu"
)

func (db *Ctx[k, v]) toKeyStr(key k) (keyStr string, err error) {
	//if k == nil {
	vv := reflect.ValueOf(key)
	if !vv.IsValid() || (vv.Kind() == reflect.Ptr && vv.IsNil()) {
		return keyStr, saavuu.ErrInvalidField
	}
	//if key is a string, directly append to keyBytes
	if strkey, ok := interface{}(key).(string); ok {
		return strkey, nil
	}
	if keyBytes, err := json.Marshal(key); err != nil {
		return keyStr, err
	} else {
		return string(keyBytes), nil
	}
}
func (db *Ctx[k, v]) toValueStr(value v) (valueStr string, err error) {
	//marshal with msgpack
	//nil value is allowed
	if bytes, err := msgpack.Marshal(value); err != nil {
		return valueStr, err
	} else {
		return string(bytes), nil
	}
}

func (db *Ctx[k, v]) toValueStrs(values []v) (valueStrs []string, err error) {
	var bytes []byte
	for _, value := range values {
		if bytes, err = msgpack.Marshal(value); err != nil {
			return nil, err
		}
		valueStrs = append(valueStrs, string(bytes))
	}
	return valueStrs, nil
}
func (db *Ctx[k, v]) toKeyStrs(keys ...k) (KeyStrs []string, err error) {
	var keyStr string
	for _, key := range keys {
		if keyStr, err = db.toKeyStr(key); err != nil {
			return nil, err
		}
		KeyStrs = append(KeyStrs, keyStr)
	}
	return KeyStrs, nil
}

func (db *Ctx[k, v]) toKeyValueStrs(keyValue ...interface{}) (keyValStrs []string, err error) {
	var (
		key              k
		value            v
		strkey, strvalue string
	)
	if len(keyValue) == 0 {
		return keyValStrs, fmt.Errorf("key value is nil")
	}
	// if key value is a map, convert it to key value slice
	if kvMap, ok := keyValue[0].(map[k]v); ok {
		for key, value := range kvMap {
			if strkey, err = db.toKeyStr(key); err != nil {
				return nil, err
			}
			if strvalue, err = db.toValueStr(value); err != nil {
				return nil, err
			}
			keyValStrs = append(keyValStrs, strkey, strvalue)
		}
	} else if l := len(keyValue); l%2 == 0 {
		for i := 0; i < l; i += 2 {
			//type check, should be of type k and v
			if key, ok = interface{}(keyValue[i]).(k); !ok {
				log.Info().Any(" key must be of type k", key)
				return nil, saavuu.ErrInvalidField
			}
			if value, ok = interface{}(keyValue[i+1]).(v); !ok {
				log.Info().Any(" value must be of type v", value)
				return nil, saavuu.ErrInvalidField
			}
			if strkey, err = db.toKeyStr(key); err != nil {
				return nil, err
			}
			if strvalue, err = db.toValueStr(value); err != nil {
				return nil, err
			}

			keyValStrs = append(keyValStrs, strkey, strvalue)
		}
	} else {
		return nil, saavuu.ErrInvalidField
	}
	return keyValStrs, nil
}
