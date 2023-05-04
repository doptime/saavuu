package api

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/data"
	"github.com/yangkequn/saavuu/logger"
)

var ErrBackTo = errors.New("param[\"backTo\"] is not a string")

// crate Api context. the created context is used :
//  1. to call api service,using Do() or DoAt()
//  2. to be called by web client or another language client
//
// ServiceName is defined as "In" + ServiceName in the InServiceName parameter
// ServiceName is automatically converted to lower case
func Api[i any, o any](f func(InServiceName i) (ret o, err error)) (ctx *Ctx[i, o]) {
	//get ServiceName
	var ServiceName string
	_type := reflect.TypeOf((*i)(nil))
	//take name of type v as key
	for _type.Kind() == reflect.Ptr {
		_type = _type.Elem()
	}
	//ensure SerivceName Starts with "In",
	//1. To make sure you know the definition of ServiceName
	//2. By noticing "In" in ServiceName, you can easily find the definition of Service
	if ServiceName = _type.Name(); len(ServiceName) < 3 || ServiceName[0:2] != "In" {
		logger.Lshortfile.Panic("Api: ServiceName is empty")
	}
	//remove "In" from ServiceName
	ServiceName = ServiceName[2:]
	//first byte of ServiceName should be lower case
	if ServiceName[0] >= 'A' && ServiceName[0] <= 'Z' {
		ServiceName = string(ServiceName[0]+32) + ServiceName[1:]
	}
	//Serivce name should Start with "api:"
	ServiceName = "api:" + ServiceName

	//create Api context
	ctx = New[i, o](ServiceName)
	ctx.Func = f
	//create a goroutine to process the job
	ProcessOneJob := func(BackToID string, s []byte) (err error) {
		var (
			out            o
			marshaledBytes []byte
			param          map[string]interface{} = map[string]interface{}{}
		)
		if err = msgpack.Unmarshal(s, &param); err != nil {
			return err
		}
		//process one job
		//check configureation is loaded
		if config.Rds == nil {
			logger.Lshortfile.Panic("config.ParamRedis is nil. Call config.ApiInitial first")
		}

		vType := reflect.TypeOf((*i)(nil)).Elem()
		if vType.Kind() == reflect.Ptr {
			vValue := reflect.New(vType.Elem()).Interface().(i)
			if ctx.Debug {
				//just allow stop here to see the input data
				ctx.Debug = !ctx.Debug
				ctx.Debug = !ctx.Debug
			}
			if err = data.MapsToStructure(param, vValue); err != nil {
				return err
			}
			if out, err = f(vValue); err != nil {
				return err
			}

		} else {
			vValueWithPointer := reflect.New(vType).Interface().(*i)
			if err = data.MapsToStructure(param, vValueWithPointer); err != nil {
				return err
			}
			if out, err = f(*vValueWithPointer); err != nil {
				return err
			}
		}

		//Post Back
		if marshaledBytes, err = msgpack.Marshal(out); err != nil {
			return err
		}
		ctx := context.Background()
		pipline := config.Rds.Pipeline()
		pipline.RPush(ctx, BackToID, marshaledBytes)
		pipline.Expire(ctx, BackToID, time.Second*20)
		_, err = pipline.Exec(ctx)
		return err
	}
	//register Api
	apiServices[ctx.ServiceName] = &ApiInfo{
		ApiName: ctx.ServiceName,
		ApiFunc: ProcessOneJob,
	}
	//return Api context
	return ctx
}
