package core

import (
	"reflect"
)

func IsNil(v interface{}) bool {
	if v == nil {
		return true
	}
	value := reflect.ValueOf(v)
	kind := value.Kind()
	switch kind {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice: //指针类型
		if value.IsNil() {
			return true
		}
	case reflect.String, reflect.Bool: //基础类型
		return false
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64: //基础类型
		return false
	case reflect.Float32, reflect.Float64: //基础类型
		return false
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64: //基础类型
		return false
	case reflect.Complex64, reflect.Complex128: //基础类型
		return false
	}

	if kind == reflect.Struct { //普通结构体
		for i := 0; i < value.NumField(); i++ {
			if !value.Field(i).IsZero() {
				return false
			}
		}
		return true
	}

	if !value.IsValid() {
		return true
	}

	if value == reflect.ValueOf(nil) {
		return true
	}

	return false
}
