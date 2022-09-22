package saavuu

import (
	"github.com/mitchellh/mapstructure"
	"github.com/vmihailenco/msgpack/v5"
)

func MapsToStructure(parmIn map[string]interface{}, outStruct interface{}) (err error) {
	msgPack, ok := parmIn["MsgPack"].([]byte)
	if ok {
		return msgpack.Unmarshal(msgPack, outStruct)
	}
	return mapstructure.Decode(parmIn, &outStruct)
}
