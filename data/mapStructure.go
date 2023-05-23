package data

import (
	"reflect"

	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

func MapsToStructure(parmIn map[string]interface{}, outStruct interface{}) (err error) {
	var (
		bytes, bytes_web []byte
		ok               bool
	)
	//the source of MsgPack if from web client
	if bytes_web, ok = parmIn["MsgPack"].([]byte); ok {
		delete(parmIn, "MsgPack")
	}
	if bytes, err = msgpack.Marshal(parmIn); err != nil {
		return err
	}
	if err = msgpack.Unmarshal(bytes, outStruct); err != nil {
		return err
	}
	if ok {
		//msgpackBytes should unmarshal after mapstructure.Decode
		//allowing JWT_*** to Cover the value in outStruct
		if err = msgpack.Unmarshal(bytes_web, outStruct); err != nil {
			return err
		}
		parmIn["MsgPack"] = bytes_web
	}
	return nil
}

func MarshalSlice(members ...interface{}) (ret [][]byte, err error) {
	var bytes []byte
	ret = make([][]byte, len(members))
	for i, member := range members {
		if bytes, err = msgpack.Marshal(member); err != nil {
			return nil, err
		}
		ret[i] = bytes
	}
	return ret, nil
}

func (db *Ctx[v]) UnmarshalToSlice(members []string) (out []v, err error) {
	out = make([]v, 0, len(members))
	//unmarshal each member in cmd.Result() using msgpack,to the type of element of out
	elemType := reflect.TypeOf(out).Elem()
	//don't set elemType to elemType.Elem() again, because out is a slice of pointer
	for _, member := range members {
		elem := reflect.New(elemType).Interface()
		if err := msgpack.Unmarshal([]byte(member), elem); err != nil {
			return out, err
		}
		out = append(out, *elem.(*v))
	}

	return out, nil
}

func (db *Ctx[v]) UnmarshalRedisZ(members []redis.Z) (out []v, scores []float64, err error) {
	var (
		str string
		ok  bool
	)
	out = make([]v, 0, len(members))
	//unmarshal each member in cmd.Result() using msgpack,to the type of element of out
	elemType := reflect.TypeOf(out).Elem()
	scores = make([]float64, len(members))
	for i, member := range members {
		if str, ok = member.Member.(string); !ok || str == "" {
			continue
		}
		elem := reflect.New(elemType).Interface()
		if err := msgpack.Unmarshal([]byte(str), elem); err != nil {
			return nil, nil, err
		}
		out = append(out, *elem.(*v))

		scores[i] = member.Score
	}
	return out, scores, nil
}
