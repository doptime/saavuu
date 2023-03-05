package data

import (
	"fmt"
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

var ErrOutSliceType = fmt.Errorf("out should be *[] Type")

func UnmarshalStrings(members []string, outSlice interface{}) (err error) {
	//out should be *[] Type
	if reflect.TypeOf(outSlice).Kind() != reflect.Ptr || reflect.TypeOf(outSlice).Elem().Kind() != reflect.Slice {
		return ErrOutSliceType
	}
	//unmarshal each member in cmd.Result() using msgpack,to the type of element of out
	elemType := reflect.TypeOf(outSlice).Elem().Elem()
	//don't set elemType to elemType.Elem() again, because out is a slice of pointer
	for _, member := range members {
		elem := reflect.New(elemType).Interface()
		if err := msgpack.Unmarshal([]byte(member), &elem); err != nil {
			return err
		}
		//append elem to out, elem is a pointer
		//the following code error: interface {}(string) "reflect.Set: value of type *map[string]interface {} is not assignable to type map[string]interface {}"
		//reflect.ValueOf(out).Elem().Set(reflect.Append(reflect.ValueOf(out).Elem(), reflect.ValueOf(elem)))
		reflect.ValueOf(outSlice).Elem().Set(reflect.Append(reflect.ValueOf(outSlice).Elem(), reflect.ValueOf(elem).Elem()))
	}

	return nil
}

func UnmarshalRedisZ(members []redis.Z, outSlice interface{}) (scores []float64, err error) {
	var (
		str string
		ok  bool
	)
	//out should be *[] Type
	if reflect.TypeOf(outSlice).Kind() != reflect.Ptr || reflect.TypeOf(outSlice).Elem().Kind() != reflect.Slice {
		return nil, ErrOutSliceType
	}
	//unmarshal each member in cmd.Result() using msgpack,to the type of element of out
	elemType := reflect.TypeOf(outSlice).Elem().Elem()
	scores = make([]float64, len(members))
	for i, member := range members {
		if str, ok = member.Member.(string); !ok || str == "" {
			continue
		}
		elem := reflect.New(elemType)
		if err := msgpack.Unmarshal([]byte(str), elem.Interface()); err != nil {
			return nil, err
		}
		//append elem to out, elem is a pointer
		reflect.ValueOf(outSlice).Elem().Set(reflect.Append(reflect.ValueOf(outSlice).Elem(), elem.Elem()))

		scores[i] = member.Score
	}
	return scores, nil
}
func MarshalRedisZ(members ...redis.Z) {
	for i := range members {
		if members[i].Member != nil {
			members[i].Member, _ = msgpack.Marshal(members[i].Member)
		}
	}
}
