package main

import (
	"fmt"
	"testing"

	"github.com/yangkequn/saavuu/data"
	"github.com/yangkequn/saavuu/test"
)

func TestHMGET(t *testing.T) {
	// Test code goes here
	var _data = data.NewStruct[string, *test.InDemo]()
	_data.HSet("field1", &test.InDemo{Data: "value1"})
	_data.HSet("field2", &test.InDemo{Data: "value2"})
	_data.HSet("field3", &test.InDemo{Data: "value3"})
	if values, err := _data.HMGET([]string{"field1", "field2", "field3"}); err != nil {
	} else {
		for _, v := range values {
			fmt.Println(v)
		}
	}
	//remove test key
	_data.Rds.Del(_data.Ctx, _data.Key)
}
