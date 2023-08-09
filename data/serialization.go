package data

import (
	"encoding/json"
	"reflect"

	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu"
)

func (db *Ctx[k, v]) KeyValuesToStrs(keyValue []interface{}) (keyBytes [][]byte, valueBytes [][]byte, err error) {
	var (
		bytes []byte
		key   k
		value v
		ok    bool
	)
	for i, l := 0, len(keyValue); i < l; i += 2 {
		//type check, should be of type k and v
		if key, ok = interface{}(keyValue[i]).(k); !ok {
			log.Info().Any(" key must be of type k", key)
			return nil, nil, saavuu.ErrInvalidField
		}
		if value, ok = interface{}(keyValue[i+1]).(v); !ok {
			log.Info().Any(" value must be of type v", value)
			return nil, nil, saavuu.ErrInvalidField
		}
		//if key is a string, directly append to keyBytes
		if strkey, ok := interface{}(key).(string); ok {
			keyBytes = append(keyBytes, []byte(strkey))
		} else if bytes, err = json.Marshal(key); err != nil {
			return nil, nil, err
		} else {
			keyBytes = append(keyBytes, bytes)
		}

		if bytes, err = msgpack.Marshal(value); err != nil {
			return nil, nil, err
		}
		valueBytes = append(valueBytes, bytes)
	}
	return valueBytes, valueBytes, nil
}

func (db *Ctx[k, v]) ValuesToStrs(values []v) (valueBytes [][]byte, err error) {
	var bytes []byte
	for _, value := range values {
		if bytes, err = msgpack.Marshal(value); err != nil {
			return nil, err
		}
		valueBytes = append(valueBytes, bytes)
	}
	return valueBytes, nil
}

func (db *Ctx[k, v]) strsToKeys(fields []string) (keys []k, err error) {
	if _, ok := interface{}(fields).([]k); ok {
		return interface{}(fields).([]k), nil
	}
	if keys = make([]k, len(fields)); len(fields) == 0 {
		return keys, nil
	}
	keyStruct := reflect.TypeOf((*k)(nil)).Elem()
	isElemPtr := keyStruct.Kind() == reflect.Ptr

	//save all data to mapOut
	for i, val := range fields {
		if isElemPtr {
			keys[i] = reflect.New(keyStruct.Elem()).Interface().(k)
			err = json.Unmarshal([]byte(val), keys[i])
		} else {
			err = json.Unmarshal([]byte(val), &keys[i])
		}
		if err != nil {
			log.Info().AnErr("HKeys: field unmarshal error:", err)
			continue
		}
	}
	return keys, nil
}

// unmarhsal using msgpack
func (db *Ctx[k, v]) strsToValues(fields ...string) (values []v, err error) {
	if values = make([]v, len(fields)); len(fields) == 0 {
		return values, nil
	}
	valueStruct := reflect.TypeOf((*v)(nil)).Elem()
	isElemPtr := valueStruct.Kind() == reflect.Ptr

	//save all data to mapOut
	for i, val := range fields {
		if isElemPtr {
			values[i] = reflect.New(valueStruct.Elem()).Interface().(v)
			err = msgpack.Unmarshal([]byte(val), values[i])
		} else {
			err = msgpack.Unmarshal([]byte(val), &values[i])
		}
		if err != nil {
			log.Info().AnErr("HVals: value unmarshal error:", err)
			continue
		}
	}
	return values, nil
}
