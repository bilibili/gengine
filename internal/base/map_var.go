package base

import (
	"errors"
	"fmt"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/internal/core"
	"reflect"
)

//support map or array
type MapVar struct {
	SourceCode
	Name   string // map name
	Intkey int64  // array index
	Strkey string // map key
	Varkey string // array index or map key
}

func (m *MapVar) Evaluate(dc *context.DataContext, Vars map[string]reflect.Value) (reflect.Value, error) {
	value, e := dc.GetValue(Vars, m.Name)
	if e != nil {
		return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code: %s, %+v", m.LineNum, m.Column, m.Code, e))
	}
	var newValue reflect.Value
	if value.Kind() == reflect.Ptr {
		newValue = value.Elem()
		if newValue.Kind() == reflect.Map {
			keyType := newValue.Type().Key()

			if len(m.Varkey) > 0 {
				key, e := dc.GetValue(Vars, m.Varkey)
				if e != nil {
					return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code: %s, %+v", m.LineNum, m.Column, m.Code, e))
				}
				wantedKey, e := core.GetWantedValue(key, keyType)
				if e != nil {
					return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code: %s, %+v", m.LineNum, m.Column, m.Code, e))
				}

				mv := value.Elem().MapIndex(wantedKey)
				if mv.IsValid() {
					return mv, nil
				} else {
					return reflect.Zero(value.Type().Elem().Elem()), nil
				}
			}

			if len(m.Strkey) > 0 {
				mv := value.Elem().MapIndex(reflect.ValueOf(m.Strkey))
				if mv.IsValid() {
					return mv, nil
				} else {
					return reflect.Zero(value.Type().Elem().Elem()), nil
				}
			}

			wantedKey, e := core.GetWantedValue(reflect.ValueOf(m.Intkey), keyType)
			if e != nil {
				return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code: %s, %+v", m.LineNum, m.Column, m.Code, e))
			}

			mv := value.Elem().MapIndex(wantedKey)
			if mv.IsValid() {
				return mv, nil
			} else {
				return reflect.Zero(value.Type().Elem().Elem()), nil
			}
		}

		if newValue.Kind() == reflect.Slice || newValue.Kind() == reflect.Array {

			if len(m.Varkey) > 0 {
				wantedKey, e := dc.GetValue(Vars, m.Varkey)
				if e != nil {
					return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code: %s, %+v", m.LineNum, m.Column, m.Code, e))
				}
				return value.Elem().Index(int(wantedKey.Int())), nil
			}

			if len(m.Strkey) > 0 {
				return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code: %s, %s", m.LineNum, m.Column, m.Code, " the index of array and slice should not be string"))
			}

			if m.Intkey >= 0 {
				return value.Elem().Index(int(m.Intkey)), nil
			} else {
				return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code %s, Slice or Array index must be non-negative!", m.LineNum, m.Column, m.Code))
			}
		}

	} else {
		newValue = value

		if newValue.Kind() == reflect.Map {
			keyType := newValue.Type().Key()

			if len(m.Varkey) > 0 {
				key, e := dc.GetValue(Vars, m.Varkey)
				if e != nil {
					return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code: %s, %+v", m.LineNum, m.Column, m.Code, e))
				}
				wantedKey, e := core.GetWantedValue(key, keyType)
				if e != nil {
					return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code: %s, %+v", m.LineNum, m.Column, m.Code, e))
				}

				mv := value.MapIndex(wantedKey)
				if mv.IsValid() {
					return mv, nil
				} else {
					return reflect.Zero(value.Type().Elem()), nil
				}
			}

			if len(m.Strkey) > 0 {
				mv := value.MapIndex(reflect.ValueOf(m.Strkey))
				if mv.IsValid() {
					return mv, nil
				} else {
					return reflect.Zero(value.Type().Elem()), nil
				}
			}

			wantedKey, e := core.GetWantedValue(reflect.ValueOf(m.Intkey), keyType)
			if e != nil {
				return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code: %s, %+v", m.LineNum, m.Column, m.Code, e))
			}

			mv := value.MapIndex(wantedKey)
			if mv.IsValid() {
				return mv, nil
			} else {
				return reflect.Zero(value.Type().Elem()), nil
			}
		}

		if newValue.Kind() == reflect.Slice || newValue.Kind() == reflect.Array {

			if len(m.Varkey) > 0 {
				wantedKey, e := dc.GetValue(Vars, m.Varkey)
				if e != nil {
					return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code: %s, %+v", m.LineNum, m.Column, m.Code, e))
				}
				return value.Index(int(wantedKey.Int())), nil
			}

			if len(m.Strkey) > 0 {
				return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code: %s, %s", m.LineNum, m.Column, m.Code, " the index of array and slice should not be string"))
			}

			if m.Intkey >= 0 {
				return value.Index(int(m.Intkey)), nil
			} else {
				return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code %s, Slice or Array index must be non-negative!", m.LineNum, m.Column, m.Code))
			}
		}
	}
	return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code: %s, Evaluate MapVarValue Only support directly-Pointer-Map, directly-Pointer-Slice and directly-Pointer-Array  or Map, Slice and Array in Pointer-Struct!", m.LineNum, m.Column, m.Code))
}

func (m *MapVar) AcceptVariable(name string) error {
	if len(m.Name) == 0 {
		m.Name = name
		return nil
	}

	if len(m.Varkey) == 0 {
		m.Varkey = name
		return nil
	}
	return errors.New("MapVar's Varkey set three times! ")
}

func (m *MapVar) AcceptInteger(i64 int64) error {
	m.Intkey = i64
	return nil
}

func (m *MapVar) AcceptString(str string) error {
	if len(m.Strkey) == 0 {
		m.Strkey = str
		return nil
	}
	return errors.New("MapVar's Strkey set three times! ")
}
