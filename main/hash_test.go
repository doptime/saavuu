package main

import (
	"fmt"
	"testing"

	"github.com/yangkequn/saavuu/data"
)

type TestHash struct {
	Name string
}

func TestStringKey(t *testing.T) {
	var (
		keys   []string
		err    error
		values []*TestHash
		value  *TestHash
	)
	//create a http context

	var keyTestInDemo = data.NewStruct[string, *TestHash]()
	if err = keyTestInDemo.HSet("field1", &TestHash{Name: "value1"}); err != nil {
		t.Error(err)
	}
	if value, err = keyTestInDemo.HGet("field1"); err != nil {
		t.Error(err)
	} else if value.Name != "value1" {
		t.Error("value.Name != value1")
	}
	keyTestInDemo.HSet("field2", &TestHash{Name: "value2"})
	values, err = keyTestInDemo.HMGET([]string{"field1", "field2"}...)
	if err != nil {
		t.Error(err)
	}
	if len(values) != 2 {
		t.Error("len(values) != 2")
	}
	//non nil value check
	if values[0] == nil {
		t.Error("values[0] == nil")
	}

	if keys, _ = keyTestInDemo.HKeys(); len(keys) != 2 {
		t.Error(fmt.Errorf("len(keys) != 2, len(keys) = %d", len(keys)))
	}
	keyTestInDemo.HDel("field1", "field2")
	if keys, _ = keyTestInDemo.HKeys(); len(keys) != 0 {
		t.Error(fmt.Errorf("len(keys) != 0, len(keys) = %d", len(keys)))
	}
}
func TestObjectKey(t *testing.T) {
	var (
		keys   []int32
		err    error
		values []*TestHash
		value  *TestHash
	)
	//create a http context

	var keyTestInDemo = data.NewStruct[int32, *TestHash]()
	if err = keyTestInDemo.HSet(123, &TestHash{Name: "value1"}); err != nil {
		t.Error(err)
	}
	if value, err = keyTestInDemo.HGet(123); err != nil {
		t.Error(err)
	} else if value.Name != "value1" {
		t.Error("value.Name != value1")
	}
	keyTestInDemo.HSet(456, &TestHash{Name: "value2"})
	values, err = keyTestInDemo.HMGET([]int32{123, 456}...)
	if err != nil {
		t.Error(err)
	}
	if len(values) != 2 {
		t.Error("len(values) != 2")
	}
	//non nil value check
	if values[0] == nil {
		t.Error("values[0] == nil")
	}

	if keys, _ = keyTestInDemo.HKeys(); len(keys) != 2 {
		t.Error(fmt.Errorf("len(keys) != 2, len(keys) = %d", len(keys)))
	}
	keyTestInDemo.HDel(123, 456)
	if keys, _ = keyTestInDemo.HKeys(); len(keys) != 0 {
		t.Error(fmt.Errorf("len(keys) != 0, len(keys) = %d", len(keys)))
	}
}
