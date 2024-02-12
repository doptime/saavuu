package api

import (
	"net/url"
	"reflect"

	"github.com/gorilla/schema"
	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/config"
	"github.com/yangkequn/saavuu/specification"
)

// crate ApiFun. the created Api can be used as normal function:
//
//	f := func(InParam *InDemo) (ret string, err error) , this is logic function
//	options. there are 2 poosible options:
//		1. api.Name("ServiceName")  //set the ServiceName of the Api. which is string. default is the name of the InParameter type but with "In" removed
//		2. api.DB("RedisDatabaseName")  //set the DB name of the job. default is the name of the function
//
// ServiceName is defined as "In" + ServiceName in the InParameter
// ServiceName is automatically converted to lower case
func Api[i any, o any](f func(InParameter i) (ret o, err error), options ...string) (retf func(InParam i) (ret o, err error)) {
	var (
		ServiceName string
		DBName      string
		ctx         *Ctx[i, o]
	)

	if ServiceName, DBName = optionsDecode(options...); len(ServiceName) == 0 {
		ServiceName = specification.TypeName((*i)(nil))
	}

	log.Debug().Str("Api service create start. name", ServiceName).Send()
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
					var form url.Values = map[string][]string{}
					if err = msgpack.Unmarshal(datapack.JsonPack, &form); err == nil {
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
					var form url.Values = map[string][]string{}
					if err = msgpack.Unmarshal(datapack.JsonPack, &form); err == nil {
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
		DBName:                    DBName,
		ApiFuncWithMsgpackedParam: ProcessOneJob,
	})
	log.Debug().Str("ApiNamed service created completed!", ServiceName).Send()
	//return Api context
	return f
}
