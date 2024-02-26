package api

import (
	"context"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/msgpack/v5"
	"github.com/yangkequn/saavuu/aopt"
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
func Api[i any, o any](f func(InParameter i) (ret o, err error), option ...aopt.Setter) (retf func(InParam i) (ret o, err error)) {
	var (
		options               *aopt.ApiOptions = aopt.MergeOptions(option...)
		NonEmptyOrZeroToCheck []int
	)
	if len(options.Name) > 0 {
		options.Name = specification.ApiName(options.Name)
	}
	if len(options.Name) == 0 {
		options.Name = specification.ApiNameByType((*i)(nil))
	}
	if len(options.Name) == 0 {
		log.Error().Str("service misnamed", options.Name).Send()
	}

	if _, ok := specification.DisAllowedServiceNames[options.Name]; ok {
		log.Error().Str("service misnamed", options.Name).Send()
	}

	log.Debug().Str("Api service create start. name", options.Name).Send()
	NonEmptyOrZeroToCheck = fieldsToCheck(reflect.TypeOf(new(i)).Elem())

	//create a goroutine to process one job
	ProcessOneJob := func(s []byte) (ret interface{}, err error) {
		var (
			in   i
			pIn  interface{}
			_map map[string]interface{} = map[string]interface{}{}
			//datapack DataPacked
		)
		//check configureation is loaded
		if config.Rds == nil {
			log.Panic().Msg("config.ParamRedis is nil.")
		}
		// case double pointer decoding
		if vType := reflect.TypeOf((*i)(nil)).Elem(); vType.Kind() == reflect.Ptr {
			pIn = reflect.New(vType.Elem()).Interface()
			in = pIn.(i)
		} else {
			pIn = reflect.New(vType).Interface()
			in = *pIn.(*i)
		}

		//type conversion of form data (from url parameter or post form)
		if err = msgpack.Unmarshal(s, &_map); err != nil {
			return nil, err
		}
		//mapstructure support type conversion
		if err = mapstructure.Decode(_map, pIn); err != nil {
			return nil, err
		}
		if len(NonEmptyOrZeroToCheck) > 0 {
			if err = checkNonEmpty(pIn, NonEmptyOrZeroToCheck); err != nil {
				return nil, err
			}
		}

		return f(in)
	}
	//register Api
	apiInfo := &ApiInfo{
		Name:                      options.Name,
		DataSource:                options.DataSource,
		ApiFuncWithMsgpackedParam: ProcessOneJob,
		Ctx:                       context.Background(),
	}
	ApiServices.Set(options.Name, apiInfo)
	funcPtr := reflect.ValueOf(f).Pointer()
	fun2ApiInfoMap.Store(funcPtr, apiInfo)
	APIGroupByDataSource.Upsert(options.DataSource, []string{}, func(exist bool, valueInMap, newValue []string) []string {
		return append(valueInMap, options.Name)
	})
	log.Debug().Str("ApiNamed service created completed!", options.Name).Send()
	//return Api context
	return f
}
