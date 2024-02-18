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
func Api[i any, o any](f func(InParameter i) (ret o, err error), options ...Option) (retf func(InParam i) (ret o, err error)) {
	var (
		option *Options = optionsMerge(options...)
	)
	if len(option.ApiName) > 0 {
		option.ApiName = specification.ApiName(option.ApiName)
	}
	if len(option.ApiName) == 0 {
		option.ApiName = specification.ApiNameByType((*i)(nil))
	}
	if len(option.ApiName) == 0 {
		log.Error().Str("service misnamed", option.ApiName).Send()
	}

	if _, ok := specification.DisAllowedServiceNames[option.ApiName]; ok {
		log.Error().Str("service misnamed", option.ApiName).Send()
	}

	log.Debug().Str("Api service create start. name", option.ApiName).Send()
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
	apiInfo := &ApiInfo{
		ApiName:                   option.ApiName,
		DbName:                    option.DbName,
		ApiFuncWithMsgpackedParam: ProcessOneJob,
	}
	ApiServices.Set(option.ApiName, apiInfo)
	fun2ApiInfo.Store(&f, apiInfo)
	log.Debug().Str("ApiNamed service created completed!", option.ApiName).Send()
	//return Api context
	return f
}
