package api

import (
	"encoding/json"
	"reflect"

	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/logger"
)

// Key purpose of ApiNamed is to allow different API to have the same input type
func ApiNamed[i any, o any](ServiceName string, f func(InParameter i) (ret o, err error)) (retf func(InParam i) (ret o, err error), ctx *Ctx[i, o]) {
	//create Api context
	//Serivce name should Start with "api:"
	ctx = New[i, o](ServiceName)
	ctx.Func = f
	//create a goroutine to process one job
	ProcessOneJob := func(s []byte) (ret interface{}, err error) {
		type DataPacked struct {
			MsgPack  []byte
			JsonPack string
		}
		var (
			in       i
			datapack DataPacked
		)
		//check configureation is loaded
		if config.Rds == nil {
			logger.Lshortfile.Panic("config.ParamRedis is nil. Call config.ApiInitial first")
		}

		//step 1, try to unmarshal MsgPack
		if err = msgpack.Unmarshal(s, &datapack); err == nil {
			// case double pointer decoding
			if vType := reflect.TypeOf((*i)(nil)).Elem(); vType.Kind() == reflect.Ptr {
				in = reflect.New(vType.Elem()).Interface().(i)
				//step 2, try to unmarshal OriginalInputParam
				if err = msgpack.Unmarshal(s, in); err == nil && len(datapack.MsgPack) > 0 {
					//step 3, try to unmarshal MsgPack
					err = msgpack.Unmarshal(datapack.MsgPack, in)
				}
				if err == nil && len(datapack.JsonPack) > 0 {
					err = json.Unmarshal([]byte(datapack.JsonPack), in)
				}

			} else {
				var pIn *i = reflect.New(vType).Interface().(*i)
				//step 2, try to unmarshal OriginalInputParam
				if err = msgpack.Unmarshal(s, pIn); err == nil && len(datapack.MsgPack) > 0 {
					//step 3, try to unmarshal MsgPackƒ
					err = msgpack.Unmarshal(datapack.MsgPack, pIn)
				}
				if err == nil && len(datapack.JsonPack) > 0 {
					err = json.Unmarshal([]byte(datapack.JsonPack), pIn)
				}
				in = *pIn
			}
		}
		if err != nil {
			//print the unmarshal error
			if ctx.Debug {
				logger.Lshortfile.Println(err)
			}
			return nil, err
		}
		return f(in)
	}
	//register Api
	ApiServices[ctx.ServiceName] = &ApiInfo{
		ApiName:                   ctx.ServiceName,
		ApiFuncWithMsgpackedParam: ProcessOneJob,
	}
	//return Api context
	return f, ctx
}

// crate ApiFun, ApiContext. the created Api is used :
//  1. to call api service,using ApiFun() or ApiContext.DoAt()
//  2. Api can be called by web client or another language client
//
// ServiceName is defined as "In" + ServiceName in the InParameter
// ServiceName is automatically converted to lower case
func Api[i any, o any](f func(InParam i) (ret o, err error)) (retf func(InParam i) (ret o, err error), ctx *Ctx[i, o]) {
	//get default ServiceName
	var _type reflect.Type
	//take name of type v as key
	for _type = reflect.TypeOf((*i)(nil)); _type.Kind() == reflect.Ptr; _type = _type.Elem() {
	}
	return ApiNamed[i, o](_type.Name(), f)
}
