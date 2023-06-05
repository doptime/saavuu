package rds

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/yangkequn/saavuu/logger"
)

func FieldsToSlice(fields interface{}) (fieldsString []string, err error) {
	//make sure fields should be a slice
	fieldsType := reflect.TypeOf(fields)
	if fieldsType.Kind() != reflect.Slice {
		logger.Lshortfile.Println("fields must be a slice")
		return nil, errors.New("fields must be a slice")
	}
	//if  fields is []string, return directly
	if fieldsString = fields.([]string); fieldsString != nil {
		return fields.([]string), nil
	}
	//now fieldsElem is not []string, marshal each field to string
	//约定，来自客户端的fields，如果是[]string，则是真实的fields,那么就不需要再次marshal
	fieldsElem := reflect.ValueOf(fields)
	//marshal each field to string
	fieldsString = make([]string, 0, fieldsElem.Len())
	for i := 0; i < fieldsElem.Len(); i++ {
		b, err := json.Marshal(reflect.ValueOf(fields).Index(i).Interface())
		if err != nil {
			logger.Lshortfile.Println("HMGET: field marshal error:", err)
			continue
		}
		fieldsString = append(fieldsString, string(b))
	}
	return fieldsString, nil
}
