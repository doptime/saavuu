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

// crate Api
// ServiceName is defined as "In" + ServiceName in the first parameter
func Api[i any, o any](f func(InServiceName i) (ret o, err error)) (ctx *Ctx[i, o]) {
	//get ServiceName
	var ServiceName string
	_type := reflect.TypeOf((*i)(nil))
	//take name of type v as key
	for _type.Kind() == reflect.Ptr {
		_type = _type.Elem()
	}

	if ServiceName = _type.Name(); len(ServiceName) < 3 || ServiceName[0:2] != "In" {
		logger.Lshortfile.Panic("Api: ServiceName is empty")
	}
	ServiceName = ServiceName[2:]
	//first byte of ServiceName should be lower case
	if ServiceName[0] >= 'A' && ServiceName[0] <= 'Z' {
		ServiceName = string(ServiceName[0]+32) + ServiceName[1:]
	}
	//create Api context
	ctx = New[i, o](ServiceName)
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
		if config.ParamRds == nil {
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
		pipline := config.ParamRds.Pipeline()
		pipline.RPush(ctx, BackToID, marshaledBytes)
		pipline.Expire(ctx, BackToID, time.Second*6)
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
