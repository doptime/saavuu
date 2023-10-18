package api

import (
	"net/url"
	"reflect"

	"github.com/gorilla/schema"
	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/config"
)

// Key purpose of ApiNamed is to allow different API to have the same input type
func ApiNamed[i any, o any](ServiceName string, f func(InParameter i) (ret o, err error)) (retf func(InParam i) (ret o, err error), ctx *Ctx[i, o]) {
	log.Debug().Str("ApiNamed service created start!", ServiceName).Send()
	//create Api context
	//Serivce name should Start with "api:"
	ctx = New[i, o](ServiceName)
	ctx.Func = f
	//create a goroutine to process one job
	ProcessOneJob := func(s []byte) (ret interface{}, err error) {
		type DataPacked struct {
			MsgPack  []byte
			JsonPack []byte
		}
		var (
			in       i
			datapack DataPacked
		)
		//check configureation is loaded
		if config.Rds == nil {
			log.Panic().Msg("config.ParamRedis is nil.")
		}

		//step 1, try to unmarshal MsgPack
		if err = msgpack.Unmarshal(s, &datapack); err == nil {
			// case double pointer decoding
			if vType := reflect.TypeOf((*i)(nil)).Elem(); vType.Kind() == reflect.Ptr {
				in = reflect.New(vType.Elem()).Interface().(i)
				//step 2, try to unmarshal jwt
				err = msgpack.Unmarshal(s, in)
				//step 3, try to unmarshal MsgPack
				if err == nil && len(datapack.MsgPack) > 0 {
					err = msgpack.Unmarshal(datapack.MsgPack, in)
				}
				//step 4, unmarshal JsonPack
				if err == nil && len(datapack.JsonPack) > 0 {
					var form url.Values = make(map[string][]string)
					err = msgpack.Unmarshal(datapack.JsonPack, form)
					if err == nil {
						err = schema.NewDecoder().Decode(in, form)
					}
				}

			} else {
				var pIn *i = reflect.New(vType).Interface().(*i)
				//step 2, try to unmarshal jwt
				err = msgpack.Unmarshal(s, pIn)
				//step 3, try to unmarshal MsgPackƒ
				if err == nil && len(datapack.MsgPack) > 0 {
					err = msgpack.Unmarshal(datapack.MsgPack, pIn)
				}
				//step 4, unmarshal JsonPack
				if err == nil && len(datapack.JsonPack) > 0 {
					var form url.Values = make(map[string][]string)
					err = msgpack.Unmarshal(datapack.JsonPack, &form)
					if err == nil {
						err = schema.NewDecoder().Decode(pIn, form)
					}
				}
				in = *pIn
			}
		}
		if err != nil {
			//print the unmarshal error
			log.Debug().AnErr("ProcessOneJob unmarshal", err).Send()
			return nil, err
		}
		return f(in)
	}
	//register Api
	ApiServices.Set(ctx.ServiceName, &ApiInfo{
		ApiName:                   ctx.ServiceName,
		ApiFuncWithMsgpackedParam: ProcessOneJob,
	})
	log.Debug().Str("ApiNamed service created completed!", ServiceName).Send()
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
	log.Debug().Msg("Api service create start")
	//get default ServiceName
	var _type reflect.Type
	//take name of type v as key
	for _type = reflect.TypeOf((*i)(nil)); _type.Kind() == reflect.Ptr; _type = _type.Elem() {
	}
	typeName := _type.Name()
	log.Debug().Str("Api service create start. name", typeName).Send()
	retf, ctx = ApiNamed[i, o](typeName, f)
	log.Debug().Str("Api service create end. name", typeName).Send()
	return retf, ctx
}
