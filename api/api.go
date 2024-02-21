package api

import (
	"context"
	"encoding/json"
	"fmt"
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
			MsgpackBody []byte
			JsonBody    []byte
			Form        url.Values
		}
		var (
			in       i
			pIn      *i
			datapack DataPacked
		)
		//check configureation is loaded
		if config.Rds == nil {
			log.Panic().Msg("config.ParamRedis is nil.")
		}
		// case double pointer decoding
		if vType := reflect.TypeOf((*i)(nil)).Elem(); vType.Kind() == reflect.Ptr {
			in = reflect.New(vType.Elem()).Interface().(i)
			pIn = &in
		} else {
			pIn = reflect.New(vType).Interface().(*i)
			in = *pIn
		}
		// the base principle of decoding is delay the decoding to the latest moment, so that each element can be decoded to the right type

		//step 1, try to unmarshal jwt
		if err = msgpack.Unmarshal(s, in); err != nil {
			return nil, err
		}
		//try to takeout DataPacked
		if err = msgpack.Unmarshal(s, &datapack); err != nil {
			return nil, fmt.Errorf("msgpack.Unmarshal DataPacked error %s", err)
		}
		//step 3, try to unmarshal MsgPack
		if len(datapack.MsgpackBody) > 0 {
			if err = msgpack.Unmarshal(datapack.MsgpackBody, in); err != nil {
				return nil, fmt.Errorf("msgpack.Unmarshal MsgpackBody error %s", err)
			}
		}
		//step 4, unmarshal Form
		if len(datapack.Form) > 0 {
			if err = schema.NewDecoder().Decode(in, datapack.Form); err != nil {
				return nil, fmt.Errorf("schema.NewDecoder().Decode error %s", err)
			}
		}
		//step 4, unmarshal JsonPack
		if len(datapack.JsonBody) > 0 {
			if err = json.Unmarshal(datapack.JsonBody, in); err != nil {
				return nil, fmt.Errorf("msgpack.Unmarshal JsonBody error %s", err)
			}
		}
		return f(in)
	}
	//register Api
	apiInfo := &ApiInfo{
		ApiName:                   option.ApiName,
		DbName:                    option.DbName,
		ApiFuncWithMsgpackedParam: ProcessOneJob,
		Ctx:                       context.Background(),
	}
	ApiServices.Set(option.ApiName, apiInfo)
	funcPtr := reflect.ValueOf(f).Pointer()
	fun2ApiInfoMap.Store(funcPtr, apiInfo)
	log.Debug().Str("ApiNamed service created completed!", option.ApiName).Send()
	//return Api context
	return f
}
