package specification

import (
	"reflect"

	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/msgpack/v5"
)

func MarshalApiInput(paramIn interface{}) (out []byte, err error) {
	//ensure the paramIn is a map or struct
	paramType := reflect.TypeOf(paramIn)
	if paramType.Kind() == reflect.Struct {
	} else if paramType.Kind() == reflect.Map {
	} else if paramType.Kind() == reflect.Ptr && (paramType.Elem().Kind() == reflect.Struct || paramType.Elem().Kind() == reflect.Map) {
	} else {
		log.Info().Msg("RdsApiBasic param should be a map or struct")
		return nil, err
	}

	if out, err = msgpack.Marshal(paramIn); err != nil {
		return nil, err
	}
	return out, nil
}
