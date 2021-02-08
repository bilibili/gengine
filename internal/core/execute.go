package core

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func InvokeFunction(obj reflect.Value, methodName string, parameters []reflect.Value) (reflect.Value, error) {
	objVal := obj

	fun := objVal.MethodByName(methodName)
	if !fun.IsValid() {
		return reflect.ValueOf(nil), errors.New(fmt.Sprintf("NOT FOUND Function: %s", methodName))
	}

	params := ParamsTypeChange(fun, parameters)
	rs := fun.Call(params)
	raw, e := GetRawTypeValue(rs)
	if e != nil {
		return reflect.ValueOf(nil), e
	}
	return raw, nil
}

/**
if want to support multi return ,change this method implements
*/
func GetRawTypeValue(rs []reflect.Value) (reflect.Value, error) {
	if len(rs) == 0 {
		return reflect.ValueOf(nil), nil
	} else {
		return rs[0], nil
	}
}

func GetStructAttributeValue(obj reflect.Value, fieldName string) (reflect.Value, error) {
	stru := obj
	var attrVal reflect.Value
	if stru.Kind() == reflect.Ptr {
		attrVal = stru.Elem().FieldByName(fieldName)
	} else {
		attrVal = stru.FieldByName(fieldName)
	}
	return attrVal, nil
}

/**
set field value
*/
func SetAttributeValue(obj reflect.Value, fieldName string, value reflect.Value) error {
	field := reflect.ValueOf(nil)
	objType := obj.Type()
	objVal := obj
	if objType.Kind() == reflect.Ptr {
		//it points to struct
		if objType.Elem().Kind() == reflect.Struct {
			field = objVal.Elem().FieldByName(fieldName)
		}
	} else {
		//not a pointer.
		if objType.Kind() == reflect.Struct {
			field = objVal.FieldByName(fieldName)
		}
	}

	if field == reflect.ValueOf(nil) {
		return errors.New(fmt.Sprintf("struct has no this field: %s", fieldName))
	}

	if field.CanSet() {
		typeName := value.Type().String()
		switch field.Type().Kind() {
		case reflect.String:
			field.SetString(value.String())
			break
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if strings.HasPrefix(typeName, "uint") {
				field.SetInt(int64(value.Uint()))
				return nil
			}
			if strings.HasPrefix(typeName, "float") {
				field.SetInt(int64(value.Float()))
				return nil
			}
			field.SetInt(value.Int())
			break
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if strings.HasPrefix(typeName, "int") && value.Int() >= 0 {
				field.SetUint(uint64(value.Int()))
				return nil
			}
			if strings.HasPrefix(typeName, "float") && value.Float() >= 0 {
				field.SetUint(uint64(value.Float()))
				return nil
			}
			field.SetUint(value.Uint())
			break
		case reflect.Float32, reflect.Float64:
			if strings.HasPrefix(typeName, "int") {
				field.SetFloat(float64(value.Int()))
				return nil
			}
			if strings.HasPrefix(typeName, "uint") {
				field.SetFloat(float64(value.Uint()))
				return nil
			}
			field.SetFloat(value.Float())
			break
		case reflect.Bool:
			field.SetBool(value.Bool())
			break
		case reflect.Slice:
			field.Set(value)
			break
		case reflect.Map:
			field.Set(value)
			break
		case reflect.Array:
			field.Set(value)
			break
		case reflect.Struct:
			field.Set(value)
			break
		case reflect.Interface:
			field.Set(value)
			break
		case reflect.Chan:
			field.Set(value)
			break
		case reflect.Complex64:
			field.SetComplex(value.Complex())
			break
		case reflect.Complex128:
			field.SetComplex(value.Complex())
			break
		case reflect.Func:
			field.Set(value)
			break
		default:
			return errors.New(fmt.Sprintf("Not support type:%s", field.Type().Kind().String()))
		}
	} else {
		return errors.New(fmt.Sprintf("%s:must Be Assignable, it should be or be in addressable value!", field.Type().Kind().String()))
	}
	return nil
}

//set single value
func SetSingleValue(obj reflect.Value, fieldName string, value reflect.Value) error {

	if obj.Kind() == reflect.Ptr {
		if value.Kind() == reflect.Ptr {
			//both ptr
			value = value.Elem()
		}

		objKind := obj.Elem().Kind()
		valueKind := value.Kind()
		if objKind == valueKind {
			obj.Elem().Set(value)
			return nil
		} else {
			valueKindStr := valueKind.String()
			switch objKind {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if strings.HasPrefix(valueKindStr, "int") {
					obj.Elem().SetInt(value.Int())
					return nil
				}
				if strings.HasPrefix(valueKindStr, "float") {
					obj.Elem().SetInt(int64(value.Float()))
					return nil
				}
				if strings.HasPrefix(valueKindStr, "uint") {
					obj.Elem().SetInt(int64(value.Uint()))
					return nil
				}
				break
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if strings.HasPrefix(valueKindStr, "int") && value.Int() >= 0 {
					obj.Elem().SetUint(uint64(value.Int()))
					return nil
				}
				if strings.HasPrefix(valueKindStr, "float") && value.Float() >= 0 {
					obj.Elem().SetUint(uint64(value.Float()))
					return nil
				}
				if strings.HasPrefix(valueKindStr, "uint") {
					obj.Elem().SetUint(value.Uint())
					return nil
				}
				break
			case reflect.Float32, reflect.Float64:
				if strings.HasPrefix(valueKindStr, "int") {
					obj.Elem().SetFloat(float64(value.Int()))
					return nil
				}
				if strings.HasPrefix(valueKindStr, "float") {
					obj.Elem().SetFloat(value.Float())
					return nil
				}
				if strings.HasPrefix(valueKindStr, "uint") {
					obj.Elem().SetFloat(float64(value.Uint()))
					return nil
				}
				break
			}
			return errors.New(fmt.Sprintf("\"%s\" value type \"%s\" is different from \"%s\" ", fieldName, obj.Elem().Kind().String(), value.Kind().String()))
		}
	} else {
		return errors.New(fmt.Sprintf("\"%s\" value is unassignable!", fieldName))
	}
}

const (
	_int   = 1
	_uint  = 2
	_float = 3
)

/*
number type exchange
*/
func ParamsTypeChange(f reflect.Value, params []reflect.Value) []reflect.Value {
	tf := f.Type()
	if tf.Kind() == reflect.Ptr {
		tf = tf.Elem()
	}
	plen := tf.NumIn()
	for i := 0; i < plen; i++ {
		switch tf.In(i).Kind() {
		case reflect.Int:
			tag := getNumType(params[i])
			if tag == _int {
				params[i] = reflect.ValueOf(int(params[i].Int()))
			} else if tag == _uint {
				params[i] = reflect.ValueOf(int(params[i].Uint()))
			} else {
				params[i] = reflect.ValueOf(int(params[i].Float()))
			}
			break
		case reflect.Int8:
			tag := getNumType(params[i])
			if tag == _int {
				params[i] = reflect.ValueOf(int8(params[i].Int()))
			} else if tag == _uint {
				params[i] = reflect.ValueOf(int8(params[i].Uint()))
			} else {
				params[i] = reflect.ValueOf(int8(params[i].Float()))
			}
			break
		case reflect.Int16:
			tag := getNumType(params[i])
			if tag == _int {
				params[i] = reflect.ValueOf(int16(params[i].Int()))
			} else if tag == _uint {
				params[i] = reflect.ValueOf(int16(params[i].Uint()))
			} else {
				params[i] = reflect.ValueOf(int16(params[i].Float()))
			}
			break
		case reflect.Int32:
			tag := getNumType(params[i])
			if tag == _int {
				params[i] = reflect.ValueOf(int32(params[i].Int()))
			} else if tag == _uint {
				params[i] = reflect.ValueOf(int32(params[i].Uint()))
			} else {
				params[i] = reflect.ValueOf(int32(params[i].Float()))
			}
			break
		case reflect.Int64:
			tag := getNumType(params[i])
			if tag == _int {
				params[i] = reflect.ValueOf(params[i].Int())
			} else if tag == _uint {
				params[i] = reflect.ValueOf(int64(params[i].Uint()))
			} else {
				params[i] = reflect.ValueOf(int64(params[i].Float()))
			}
			break
		case reflect.Uint:
			tag := getNumType(params[i])
			if tag == _int {
				params[i] = reflect.ValueOf(uint(params[i].Int()))
			} else if tag == _uint {
				params[i] = reflect.ValueOf(uint(params[i].Uint()))
			} else {
				params[i] = reflect.ValueOf(uint(params[i].Float()))
			}
			break
		case reflect.Uint8:
			tag := getNumType(params[i])
			if tag == _int {
				params[i] = reflect.ValueOf(uint8(params[i].Int()))
			} else if tag == _uint {
				params[i] = reflect.ValueOf(uint8(params[i].Uint()))
			} else {
				params[i] = reflect.ValueOf(uint8(params[i].Float()))
			}
			break
		case reflect.Uint16:
			tag := getNumType(params[i])
			if tag == _int {
				params[i] = reflect.ValueOf(uint16(params[i].Int()))
			} else if tag == _uint {
				params[i] = reflect.ValueOf(uint16(params[i].Uint()))
			} else {
				params[i] = reflect.ValueOf(uint16(params[i].Float()))
			}
			break
		case reflect.Uint32:
			tag := getNumType(params[i])
			if tag == _int {
				params[i] = reflect.ValueOf(uint32(params[i].Int()))
			} else if tag == _uint {
				params[i] = reflect.ValueOf(uint32(params[i].Uint()))
			} else {
				params[i] = reflect.ValueOf(uint32(params[i].Float()))
			}
			break
		case reflect.Uint64:
			tag := getNumType(params[i])
			if tag == _int {
				params[i] = reflect.ValueOf(uint64(params[i].Int()))
			} else if tag == _uint {
				params[i] = reflect.ValueOf(params[i].Uint())
			} else {
				params[i] = reflect.ValueOf(uint64(params[i].Float()))
			}
			break
		case reflect.Float32:
			tag := getNumType(params[i])
			if tag == _int {
				params[i] = reflect.ValueOf(float32(params[i].Int()))
			} else if tag == _uint {
				params[i] = reflect.ValueOf(float32(params[i].Uint()))
			} else {
				params[i] = reflect.ValueOf(float32(params[i].Float()))
			}
			break
		case reflect.Float64:
			tag := getNumType(params[i])
			if tag == _int {
				params[i] = reflect.ValueOf(float64(params[i].Int()))
			} else if tag == _uint {
				params[i] = reflect.ValueOf(float64(params[i].Uint()))
			} else {
				params[i] = reflect.ValueOf(params[i].Float())
			}
			break
		case reflect.Ptr:
			break
		case reflect.Interface:
			if !reflect.ValueOf(params[i]).IsValid() {
				params[i] = reflect.New(tf.In(i))
			}
		default:
			continue
		}
	}
	return params
}

func getNumType(param reflect.Value) int {
	ts := param.Kind().String()
	if strings.HasPrefix(ts, "int") {
		return _int
	}

	if strings.HasPrefix(ts, "uint") {
		return _uint
	}

	if strings.HasPrefix(ts, "float") {
		return _float
	}

	panic(fmt.Sprintf("it is not number type, type is %s !", ts))
}

func GetWantedValue(newValue reflect.Value, toKind reflect.Type) (reflect.Value, error) {
	if newValue.Kind() == toKind.Kind() {
		return newValue, nil
	}

	switch toKind.Kind() {
	case reflect.Int:
		return reflect.ValueOf(int(newValue.Int())), nil
	case reflect.Int8:
		return reflect.ValueOf(int8(newValue.Int())), nil
	case reflect.Int16:
		return reflect.ValueOf(int16(newValue.Int())), nil
	case reflect.Int32:
		return reflect.ValueOf(int32(newValue.Int())), nil
	case reflect.Int64:
		return newValue, nil

	case reflect.Uint:
		return reflect.ValueOf(uint(newValue.Uint())), nil
	case reflect.Uint8:
		return reflect.ValueOf(uint8(newValue.Uint())), nil
	case reflect.Uint16:
		return reflect.ValueOf(uint16(newValue.Uint())), nil
	case reflect.Uint32:
		return reflect.ValueOf(uint32(newValue.Uint())), nil
	case reflect.Uint64:
		return newValue, nil

	case reflect.Float32:
		return reflect.ValueOf(float32(newValue.Float())), nil
	case reflect.Float64:
		return newValue, nil
	}

	return newValue, nil
}
