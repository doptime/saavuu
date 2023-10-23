package https

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	"github.com/yangkequn/saavuu/config"

	"github.com/vmihailenco/msgpack/v5"
)

// Parse parses the request.
func ToStruct(r *http.Request, v interface{}) error {
	//parse json body
	if r.ContentLength > 0 && (r.Header.Get("Content-Type") == "application/json") {
		json.NewDecoder(r.Body).Decode(v)
	} else if r.ContentLength > 0 && (r.Header.Get("Content-Type") == "application/octet-stream") {
		msgpack.NewDecoder(r.Body).Decode(v)
	} else if r.ContentLength > 0 && (r.Header.Get("Content-Type") == "multipart/form-data") {
		m := map[string]interface{}{}
		//parse form multipart form data
		if err := r.ParseMultipartForm(config.Cfg.Http.MaxBufferSize); err == nil {
			for k, v := range r.MultipartForm.Value {
				if len(v) == 1 {
					m[k] = v[0]
				} else {
					m[k] = v
				}
			}
			for k, v := range r.MultipartForm.File {
				m[k] = v
			}
		}
		//Convert map to struct
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
			return errors.New("v is not a pointer to a struct")
		}
		rt := rv.Type().Elem()
		for i := 0; i < rt.NumField(); i++ {
			field := rt.Field(i)
			fieldName := field.Name
			//set value according to map
			if value, ok := m[fieldName]; ok {
				rv.Elem().Field(i).Set(reflect.ValueOf(value))
			}
		}
	}
	return nil
}
