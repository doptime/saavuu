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

func (ctx *Ctx[v]) Serve(f func(paramIn v) (ret interface{}, err error)) {
	ProcessOneJob := func(BackToID string, s []byte) (err error) {
		var (
			out            interface{}
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

		vType := reflect.TypeOf((*v)(nil)).Elem()
		if vType.Kind() == reflect.Ptr {
			vValue := reflect.New(vType.Elem()).Interface().(v)
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
			vValueWithPointer := reflect.New(vType).Interface().(*v)
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
	apiServices[ctx.ServiceName] = &ApiInfo{
		ApiName: ctx.ServiceName,
		ApiFunc: ProcessOneJob,
	}
}
