package data

import (
	"encoding/json"
	"reflect"

	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/msgpack/v5"
)

func (db *Ctx[k, v]) toKeys(valStr []string) (keys []k, err error) {
	if _, ok := interface{}(valStr).([]k); ok {
		return interface{}(valStr).([]k), nil
	}
	if keys = make([]k, len(valStr)); len(valStr) == 0 {
		return keys, nil
	}
	keyStruct := reflect.TypeOf((*k)(nil)).Elem()
	isElemPtr := keyStruct.Kind() == reflect.Ptr

	//save all data to mapOut
	for i, val := range valStr {
		if isElemPtr {
			keys[i] = reflect.New(keyStruct.Elem()).Interface().(k)
			err = json.Unmarshal([]byte(val), keys[i])
		} else {
			//if key is type of string, just return string
			if keyStruct.Kind() == reflect.String {
				keys[i] = interface{}(string(val)).(k)
			} else {
				err = json.Unmarshal([]byte(val), &keys[i])
			}
		}
		if err != nil {
			log.Info().AnErr("HKeys: field unmarshal error:", err).Msgf("Key: %s", db.Key)
			continue
		}
	}
	return keys, nil
}

// unmarhsal using msgpack
func (db *Ctx[k, v]) toValues(valStr ...string) (values []v, err error) {
	if values = make([]v, len(valStr)); len(valStr) == 0 {
		return values, nil
	}
	valueStruct := reflect.TypeOf((*v)(nil)).Elem()
	isElemPtr := valueStruct.Kind() == reflect.Ptr

	//save all data to mapOut
	for i, val := range valStr {
		if isElemPtr {
			values[i] = reflect.New(valueStruct.Elem()).Interface().(v)
			err = msgpack.Unmarshal([]byte(val), values[i])
		} else {
			err = msgpack.Unmarshal([]byte(val), &values[i])
		}
		if err != nil {
			log.Info().AnErr("HVals: value unmarshal error:", err).Msgf("Key: %s", db.Key)
			continue
		}
	}
	return values, nil
}
func (db *Ctx[k, v]) toValue(valbytes []byte) (value v, err error) {
	valueStruct := reflect.TypeOf((*v)(nil)).Elem()
	isElemPtr := valueStruct.Kind() == reflect.Ptr
	if isElemPtr {
		value = reflect.New(valueStruct.Elem()).Interface().(v)
		return value, msgpack.Unmarshal(valbytes, value)
	} else {
		err = msgpack.Unmarshal(valbytes, &value)
		return value, err
	}
}

func (db *Ctx[k, v]) toKey(valBytes []byte) (key k, err error) {
	keyStruct := reflect.TypeOf((*k)(nil)).Elem()
	isElemPtr := keyStruct.Kind() == reflect.Ptr
	if isElemPtr {
		key = reflect.New(keyStruct.Elem()).Interface().(k)
		return key, json.Unmarshal(valBytes, key)
	} else {
		//if key is type of string, just return string
		if keyStruct.Kind() == reflect.String {
			return reflect.ValueOf(string(valBytes)).Interface().(k), nil
		}
		err = json.Unmarshal(valBytes, &key)
		return key, err
	}
}
