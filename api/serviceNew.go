package api

import (
	"errors"
	"reflect"

	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/data"
	"github.com/yangkequn/saavuu/logger"
)

var ErrBackTo = errors.New("param[\"backTo\"] is not a string")

// Key purpose of ApiNamed is to allow different API to have the same input type
func ApiNamed[i any, o any](ServiceName string, f func(InServiceName i) (ret o, err error)) (retf func(InParam i) (ret o, err error), ctx *Ctx[i, o]) {
	//create Api context
	//Serivce name should Start with "api:"
	ctx = New[i, o](ServiceName)
	ctx.Func = f
	//create a goroutine to process the job
	ProcessOneJob := func(s []byte) (ret interface{}, err error) {
		var (
			out    o
			param  map[string]interface{} = map[string]interface{}{}
			vValue i
		)
		if err = msgpack.Unmarshal(s, &param); err != nil {
			return nil, err
		}
		//process one job
		//check configureation is loaded
		if config.Rds == nil {
			logger.Lshortfile.Panic("config.ParamRedis is nil. Call config.ApiInitial first")
		}

		vType := reflect.TypeOf((*i)(nil)).Elem()
		if vType.Kind() == reflect.Ptr {
			vValue = reflect.New(vType.Elem()).Interface().(i)
			err = data.MapsToStructure(param, vValue)
		} else {
			var valueWithPointer *i = reflect.New(vType).Interface().(*i)
			//MapsToStructure does not support double pointer decoding
			if err = data.MapsToStructure(param, valueWithPointer); err == nil {
				vValue = *valueWithPointer
			}
		}
		//print the unmarshal error
		if err != nil {
			if ctx.Debug {
				logger.Lshortfile.Println(err)
			}
			return nil, err
		}
		if out, err = f(vValue); err != nil {
			return nil, err
		}

		return out, err
	}
	//register Api
	ApiServices[ctx.ServiceName] = &ApiInfo{
		ApiName:                   ctx.ServiceName,
		ApiFuncWithMsgpackedParam: ProcessOneJob,
	}
	//return Api context
	return f, ctx
}

// crate Api context. the created context is used :
//  1. to call api service,using Do() or DoAt()
//  2. to be called by web client or another language client
//
// ServiceName is defined as "In" + ServiceName in the InServiceName parameter
// ServiceName is automatically converted to lower case
func Api[i any, o any](f func(InParam i) (ret o, err error)) (retf func(InParam i) (ret o, err error), ctx *Ctx[i, o]) {
	//get default ServiceName
	var _type reflect.Type
	//take name of type v as key
	for _type = reflect.TypeOf((*i)(nil)); _type.Kind() == reflect.Ptr; _type = _type.Elem() {
	}
	return ApiNamed[i, o](_type.Name(), f)
}
