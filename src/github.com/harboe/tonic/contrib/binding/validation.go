package binding

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
)

type Validator func(v interface{}) bool
type Validators map[string]Validator

var validators = Validators{
	"required": func(v interface{}) bool {
		typ := reflect.TypeOf(v)
		zero := reflect.Zero(typ).Interface()

		return !reflect.DeepEqual(zero, v)
	},
	"min": func(v interface{}) bool {
		log.Println(v)
		return false
	},
	"max": func(v interface{}) bool {
		log.Println(v)
		return false
	},
}

func Validate(v interface{}, parents ...string) error {
	if v == nil {
		return nil
	}

	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	switch typ.Kind() {
	case reflect.Map:
		return validateMap(typ, val, v, parents...)
	case reflect.Slice:
		return validateSlice(typ, val, v, parents...)
	case reflect.Struct:
		return validateStruct(typ, val, v, parents...)
	}

	return nil
}

func validateMap(typ reflect.Type, val reflect.Value, v interface{}, parents ...string) error {
	for _, key := range val.MapKeys() {
		elmVal := val.MapIndex(key).Interface()

		// ensure the map doesn't contains empty keys.
		if validator, ok := validators["required"]; ok {
			if ok := validator(key.Interface()); !ok {
				return errors.New("validation: " + strings.Join(parents, ".") + " contains an empty key")
			}
		}

		if err := Validate(elmVal, append(parents, fmt.Sprintf("[%v]", key))...); err != nil {
			return err
		}
	}

	return nil
}

func validateSlice(typ reflect.Type, val reflect.Value, v interface{}, parents ...string) error {
	for i := 0; i < val.Len(); i++ {
		elmVal := val.Index(i).Interface()

		if err := Validate(elmVal, append(parents, fmt.Sprintf("[%v]", i))...); err != nil {
			return err
		}
	}

	return nil
}

func validateStruct(typ reflect.Type, val reflect.Value, v interface{}, parents ...string) error {
	for i := 0; i < typ.NumField(); i++ {
		v := val.Field(i)
		f := typ.Field(i)

		if !v.CanInterface() {
			continue
		}

		tag := f.Tag.Get("binding")
		val := v.Interface()

		// log.Println(f.Name, "=>", val, "tag:", tag, f.Tag)
		if err := validators.Validate(tag, val); err != nil {
			return err // errors.New("validation: " + tag + " " + strings.Join(append(parents, f.Name), "."))
		}

		// if validator, ok := validators[tag]; ok {
		// 	if ok := validator(val); !ok {
		// 		return errors.New("validation: " + tag + " " + strings.Join(append(parents, f.Name), "."))
		// 	}
		// }

		if err := Validate(val, append(parents, f.Name)...); err != nil {
			return err
		}
	}

	return nil
}

func (v Validators) Validate(binding string, val interface{}) error {

	return nil
}
