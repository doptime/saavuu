package tools

import (
	"github.com/mitchellh/mapstructure"
	"github.com/vmihailenco/msgpack/v5"
)

func MapsToStructure(parmIn map[string]interface{}, outStruct interface{}) (err error) {
	msgPack, ok := parmIn["MsgPack"].([]byte)
	if ok {
		delete(parmIn, "MsgPack")
	}
	if err = msgpack.Unmarshal(msgPack, outStruct); err != nil {
		return err
	}
	if err = mapstructure.Decode(parmIn, outStruct); err != nil {
		return err
	}
	if ok {
		parmIn["MsgPack"] = msgPack
	}
	return nil
}
