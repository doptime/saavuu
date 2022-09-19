package tools

import "github.com/vmihailenco/msgpack/v5"

func MapToStruct(in map[string]interface{}, out interface{}) (err error) {
	bytes, err := msgpack.Marshal(in)
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(bytes, out)
}
