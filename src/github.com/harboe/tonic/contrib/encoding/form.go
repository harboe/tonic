package encoding

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/harboe/tonic/contrib/binding"
)

type formEncoder struct{}

func (e formEncoder) ContentType() string {
	return "application/x-www-form-urlencoded"
}

func (e formEncoder) Decode(req *http.Request, v interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	if err := mapForm(v, req.Form); err != nil {
		return err
	}
	return binding.Validate(v)
}

func mapForm(ptr interface{}, form map[string][]string) error {
	typ := reflect.TypeOf(ptr).Elem()
	formStruct := reflect.ValueOf(ptr).Elem()
	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		if inputFieldName := typeField.Tag.Get("form"); inputFieldName != "" {
			structField := formStruct.Field(i)
			if !structField.CanSet() {
				continue
			}

			inputValue, exists := form[inputFieldName]
			if !exists {
				continue
			}
			numElems := len(inputValue)
			if structField.Kind() == reflect.Slice && numElems > 0 {
				sliceOf := structField.Type().Elem().Kind()
				slice := reflect.MakeSlice(structField.Type(), numElems, numElems)
				for i := 0; i < numElems; i++ {
					if err := setWithProperType(sliceOf, inputValue[i], slice.Index(i)); err != nil {
						return err
					}
				}
				formStruct.Elem().Field(i).Set(slice)
			} else {
				if err := setWithProperType(typeField.Type.Kind(), inputValue[0], structField); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func setWithProperType(valueKind reflect.Kind, val string, structField reflect.Value) error {
	switch valueKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val == "" {
			val = "0"
		}
		intVal, err := strconv.Atoi(val)
		if err != nil {
			return err
		} else {
			structField.SetInt(int64(intVal))
		}
	case reflect.Bool:
		if val == "" {
			val = "false"
		}
		boolVal, err := strconv.ParseBool(val)
		if err != nil {
			return err
		} else {
			structField.SetBool(boolVal)
		}
	case reflect.Float32:
		if val == "" {
			val = "0.0"
		}
		floatVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return err
		} else {
			structField.SetFloat(floatVal)
		}
	case reflect.Float64:
		if val == "" {
			val = "0.0"
		}
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		} else {
			structField.SetFloat(floatVal)
		}
	case reflect.String:
		structField.SetString(val)
	}
	return nil
}
