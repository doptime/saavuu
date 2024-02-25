package test

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

	var keyTestInDemo = data.New[string, *TestHash]()
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
	if values, err = keyTestInDemo.HMGET([]string{"field1", "field2"}...); err != nil {
		t.Error(err)
	}
	if len(values) != 2 {
		t.Error("len(values) != 3")
	}

	if keys, _ = keyTestInDemo.HKeys(); len(keys) != 2 {
		t.Error(fmt.Errorf("len(keys) != 2, len(keys) = %d", len(keys)))
	}
	keyTestInDemo.HDel("field1", "field2")
	if keys, _ = keyTestInDemo.HKeys(); len(keys) != 0 {
		t.Error(fmt.Errorf("len(keys) != 0, len(keys) = %d", len(keys)))
	}
}
func TestStringKey2(t *testing.T) {
	var (
		keys   []*string
		err    error
		values []*TestHash
		value  *TestHash
		k1, k2 string    = "key1", "key2"
		v1, v2 *TestHash = &TestHash{Name: "value1"}, &TestHash{Name: "value2"}
	)
	//create a http context

	var keyTestInDemo = data.New[*string, *TestHash]()
	if err = keyTestInDemo.HSet(&k1, v1); err != nil {
		t.Error(err)
	}
	if value, err = keyTestInDemo.HGet(&k1); err != nil {
		t.Error(err)
	} else if value.Name != "value1" {
		t.Error("value.Name != value1")
	}

	if err = keyTestInDemo.HSet(&k2, v2); err != nil {
		t.Error(err)
	}
	values, err = keyTestInDemo.HMGET([]*string{&k1, &k2}...)
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
	keyTestInDemo.HDel(&k1, &k2)
	if keys, _ = keyTestInDemo.HKeys(); len(keys) != 0 {
		t.Error(fmt.Errorf("len(keys) != 0, len(keys) = %d", len(keys)))
	}
}
func TestObjectKey(t *testing.T) {
	type Key struct {
		Id   int32
		Name string
	}
	var (
		keys  []*Key
		err   error
		value *TestHash
		k1    = &Key{Id: 123, Name: "field1"}
		k2    = &Key{Id: 456, Name: "field2"}
		v1    = &TestHash{Name: "value1"}
		v2    = &TestHash{Name: "value2"}
	)
	//create a http context

	var keyTestInDemo = data.New[*Key, *TestHash]()
	if err = keyTestInDemo.HSet(k1, v1); err != nil {
		t.Error(err)
	}
	if value, err = keyTestInDemo.HGet(k1); err != nil {
		t.Error(err)
	} else if value.Name != "value1" {
		t.Error("value.Name != value1")
	}

	if err = keyTestInDemo.HSet(k1, v1, k2, v2); err != nil {
		t.Error(err)
	}

	if keys, _ = keyTestInDemo.HKeys(); len(keys) != 2 {
		t.Error(fmt.Errorf("len(keys) != 2, len(keys) = %d", len(keys)))
	}
	keyTestInDemo.HDel(k1, k2)
	if keys, _ = keyTestInDemo.HKeys(); len(keys) != 0 {
		t.Error(fmt.Errorf("len(keys) != 0, len(keys) = %d", len(keys)))
	}
	if err = keyTestInDemo.HSet(map[*Key]*TestHash{k1: v1, k2: v2}); err != nil {
		t.Error(err)
	}
	if keys, _ = keyTestInDemo.HKeys(); len(keys) != 2 {
		t.Error(fmt.Errorf("len(keys) != 2, len(keys) = %d", len(keys)))
	}
	keyTestInDemo.HDel(k1, k2)
	if keys, _ = keyTestInDemo.HKeys(); len(keys) != 0 {
		t.Error(fmt.Errorf("len(keys) != 0, len(keys) = %d", len(keys)))
	}
}
