package api

import (
	"fmt"
	"reflect"
	"strings"
)

// fieldsToCheck checks if fields with the `nonempty` tag or `nonzero` tag are empty or zero
func fieldsToCheck(v reflect.Type) (indexes []int) {
	var (
		tag   string
		field reflect.StructField
		k     reflect.Kind
		// _v    bit 0: check nonempty, bit 1: check nonzero, bit 2+: index value of the field
		_v int
	)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	// Check if the passed interface is a pointer to a struct
	if vk := v.Kind(); vk != reflect.Struct {
		return nil
	}
	for i := 0; i < v.NumField(); i++ {
		field, k, _v = v.Field(i), v.Field(i).Type.Kind(), 0

		if tag = field.Tag.Get("mapstructure"); tag == "" {
			continue
		}
		if _i := strings.LastIndex(tag, "nonempty"); _i > 0 && tag[_i-1] == ',' || tag[_i-1] == ' ' && field.Type.Kind() == reflect.String {
			_v = 1
		}
		if _i := strings.LastIndex(tag, "nonzero"); _i > 0 && (tag[_i-1] == ',' || tag[_i-1] == ' ') && (k >= reflect.Int && k <= reflect.Float64) {
			_v = _v | 2
		}
		if _v > 0 {
			indexes = append(indexes, (i<<2)|_v)
		}
	}
	return indexes

}
func checkNonEmpty(s interface{}, indexesToCheck []int) error {
	v := reflect.ValueOf(s)

	// Check if the passed interface is a pointer to a struct
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return nil
	}

	v = v.Elem() // Dereference the pointer to get the struct
	for _, index := range indexesToCheck {
		field_i := index >> 2
		field := v.Field(field_i)
		fieldKind := field.Kind()

		if checkNonEmpty := (index & 1) == 1; checkNonEmpty && (field.Kind() == reflect.String) && field.String() == "" {
			return fmt.Errorf("%s should not be empty", v.Type().Field(field_i).Name)
		} else if nonZero, checkNonZero := true, index&2 == 2; checkNonZero {
			if fieldKind >= reflect.Int && fieldKind <= reflect.Int64 {
				nonZero = nonZero && field.Int() != 0
			} else if fieldKind >= reflect.Uint && fieldKind <= reflect.Uint64 {
				nonZero = nonZero && field.Uint() != 0
			} else if fieldKind >= reflect.Float32 && fieldKind <= reflect.Float64 {
				nonZero = nonZero && field.Float() != 0
			}
			if !nonZero {
				return fmt.Errorf("%s should not be zero", v.Type().Field(field_i).Name)
			}
		}
	}

	return nil
}
