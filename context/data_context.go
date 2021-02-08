package context

import (
	"errors"
	"fmt"
	"github.com/bilibili/gengine/internal/core"
	"path"
	"path/filepath"
	"plugin"
	"reflect"
	"strings"
	"sync"
)

type DataContext struct {
	lockVars sync.Mutex
	lockBase sync.Mutex
	base     map[string]reflect.Value
}

func NewDataContext() *DataContext {
	dc := &DataContext{
		base: make(map[string]reflect.Value),
	}
	dc.loadInnerUDF()
	return dc
}

func (dc *DataContext) loadInnerUDF() {
	dc.Add("isNil", core.IsNil)
}

func (dc *DataContext) Add(key string, obj interface{}) {
	dc.lockBase.Lock()
	defer dc.lockBase.Unlock()
	dc.base[key] = reflect.ValueOf(obj)
}

func (dc *DataContext) Del(keys ...string) {
	if len(keys) == 0 {
		return
	}
	dc.lockBase.Lock()
	defer dc.lockBase.Unlock()

	for _, key := range keys {
		delete(dc.base, key)
	}
}

//plugin_exportName_apiName.so
// _ is a separator
//plugin is prefix
//exportName is user export in plugin file
//apiName is plugin used in gengine
func (dc *DataContext) PluginLoader(absolutePathOfSO string) (string, plugin.Symbol, error) {

	plg, err := plugin.Open(absolutePathOfSO)
	if err != nil {
		return "", nil, err
	}

	_, file := filepath.Split(absolutePathOfSO)
	if path.Ext(file) != ".so" {
		return "", nil, errors.New(fmt.Sprintf("%s is not a plugin file", absolutePathOfSO))
	}

	fileWithOutExt := strings.ReplaceAll(file, ".so", "")

	splits := strings.Split(fileWithOutExt, "_")
	if len(splits) != 3 || !strings.HasPrefix(file, "plugin_") {
		return "", nil, errors.New(fmt.Sprintf("the plugin file name(%s) is not fit for need! ", absolutePathOfSO))
	}

	exportName := splits[1]
	apiName := splits[2]

	exportApi, err := plg.Lookup(exportName)
	if err != nil {
		return "", nil, err
	}

	dc.lockBase.Lock()
	defer dc.lockBase.Unlock()

	dc.base[apiName] = reflect.ValueOf(exportApi)
	return apiName, exportApi, nil
}

func (dc *DataContext) Get(key string) (reflect.Value, error) {
	dc.lockBase.Lock()
	v, ok := dc.base[key]
	dc.lockBase.Unlock()
	if ok {
		return v, nil
	} else {
		return reflect.ValueOf(nil), errors.New(fmt.Sprintf("NOT FOUND key :%s ", key))
	}
}

/**
execute the injected functions
function execute supply multi return values, but simplify ,just return one value
*/
func (dc *DataContext) ExecFunc(Vars map[string]reflect.Value, funcName string, parameters []reflect.Value) (reflect.Value, error) {
	dc.lockBase.Lock()
	v, ok := dc.base[funcName]
	dc.lockBase.Unlock()

	if ok {
		args := core.ParamsTypeChange(v, parameters)
		res := v.Call(args)
		raw, e := core.GetRawTypeValue(res)
		if e != nil {
			return reflect.ValueOf(nil), e
		}
		return raw, nil
	}

	dc.lockVars.Lock()
	vv, vok := Vars[funcName]
	dc.lockVars.Unlock()
	if vok {
		args := core.ParamsTypeChange(vv, parameters)
		res := vv.Call(args)
		raw, e := core.GetRawTypeValue(res)
		if e != nil {
			return reflect.ValueOf(nil), e
		}
		return raw, nil
	}
	return reflect.ValueOf(nil), errors.New(fmt.Sprintf("NOT FOUND function \"%s\"", funcName))
}

/**
execute the struct's functions
function execute supply multi return values, but simplify ,just return one value
*/
func (dc *DataContext) ExecMethod(Vars map[string]reflect.Value, methodName string, args []reflect.Value) (reflect.Value, error) {
	structAndMethod := strings.Split(methodName, ".")
	//Dimit rule
	if len(structAndMethod) != 2 {
		return reflect.ValueOf(nil), errors.New(fmt.Sprintf("Not supported call, just support struct.method call, now length is %d", len(structAndMethod)))
	}

	dc.lockBase.Lock()
	v, ok := dc.base[structAndMethod[0]]
	dc.lockBase.Unlock()

	if ok {
		res, err := core.InvokeFunction(v, structAndMethod[1], args)
		if err != nil {
			return reflect.ValueOf(nil), err
		}
		return res, nil
	}

	dc.lockVars.Lock()
	vv, vok := Vars[structAndMethod[0]]
	dc.lockVars.Unlock()
	if vok {
		res, err := core.InvokeFunction(vv, structAndMethod[1], args)
		if err != nil {
			return reflect.ValueOf(nil), err
		}
		return res, nil
	}
	return reflect.ValueOf(nil), errors.New(fmt.Sprintf("Not found method: %s", methodName))
}

/**
get the value user set
*/
func (dc *DataContext) GetValue(Vars map[string]reflect.Value, variable string) (reflect.Value, error) {
	if strings.Contains(variable, ".") {
		//in dataContext
		structAndField := strings.Split(variable, ".")
		//Dimit rule
		if len(structAndField) != 2 {
			return reflect.ValueOf(nil), errors.New(fmt.Sprintf("Not supported Field, just support struct.field, now length is %d", len(structAndField)))
		}

		dc.lockBase.Lock()
		v, ok := dc.base[structAndField[0]]
		dc.lockBase.Unlock()

		if ok {
			return core.GetStructAttributeValue(v, structAndField[1])
		}

		//for return struct or struct ptr
		dc.lockVars.Lock()
		obj, ok := Vars[structAndField[0]]
		dc.lockVars.Unlock()
		if ok {
			return core.GetStructAttributeValue(obj, structAndField[1])
		}
	} else {
		//user set
		dc.lockBase.Lock()
		v, ok := dc.base[variable]
		dc.lockBase.Unlock()

		if ok {
			return v, nil
		}
		//in RuleEntity
		dc.lockVars.Lock()
		res, rok := Vars[variable]
		dc.lockVars.Unlock()
		if rok {
			return res, nil
		}

	}
	return reflect.ValueOf(nil), errors.New(fmt.Sprintf("Did not found variable : %s ", variable))
}

func (dc *DataContext) SetValue(Vars map[string]reflect.Value, variable string, newValue reflect.Value) error {
	if strings.Contains(variable, ".") {
		structAndField := strings.Split(variable, ".")
		//Dimit rule
		if len(structAndField) != 2 {
			return errors.New(fmt.Sprintf("just support struct.field, now length is %d", len(structAndField)))
		}

		dc.lockBase.Lock()
		v, ok := dc.base[structAndField[0]]
		dc.lockBase.Unlock()

		if ok {
			return core.SetAttributeValue(v, structAndField[1], newValue)
		} else {
			dc.lockVars.Lock()
			vv, vok := Vars[structAndField[0]]
			dc.lockVars.Unlock()
			if vok {
				return core.SetAttributeValue(vv, structAndField[1], newValue)
			}
		}
	} else {
		dc.lockBase.Lock()
		v, ok := dc.base[variable]
		dc.lockBase.Unlock()
		if ok {
			return core.SetSingleValue(v, variable, newValue)
		} else {
			//in RuleEntity
			dc.lockVars.Lock()
			Vars[variable] = newValue
			dc.lockVars.Unlock()
			return nil
		}
	}
	return errors.New(fmt.Sprintf("setValue not found \"%s\" error.", variable))
}

func (dc *DataContext) SetMapVarValue(Vars map[string]reflect.Value, mapVarName, mapVarStrkey, mapVarVarkey string, mapVarIntkey int64, setValue reflect.Value) error {

	value, e := dc.GetValue(Vars, mapVarName)
	if e != nil {
		return e
	}

	if value.Kind() == reflect.Ptr {
		newValue := value.Elem()
		valueType := newValue.Type().Elem()

		if newValue.Kind() == reflect.Map {
			keyType := newValue.Type().Key()
			if len(mapVarVarkey) > 0 {
				key, e := dc.GetValue(Vars, mapVarVarkey)
				if e != nil {
					return e
				}
				wantedKey, e := core.GetWantedValue(key, keyType)
				if e != nil {
					return e
				}

				wantedValue, e := core.GetWantedValue(setValue, valueType)
				if e != nil {
					return e
				}
				value.Elem().SetMapIndex(wantedKey, wantedValue)
				return nil
			}

			if len(mapVarStrkey) > 0 {
				wantedValue, e := core.GetWantedValue(setValue, valueType)
				if e != nil {
					return e
				}
				value.Elem().SetMapIndex(reflect.ValueOf(mapVarStrkey), wantedValue)
				return nil
			}

			//int key
			wantedKey, e := core.GetWantedValue(reflect.ValueOf(mapVarIntkey), keyType)
			if e != nil {
				return e
			}
			wantedValue, e := core.GetWantedValue(setValue, valueType)
			if e != nil {
				return e
			}
			value.Elem().SetMapIndex(wantedKey, wantedValue)
			return nil
		}

		if newValue.Kind() == reflect.Slice || newValue.Kind() == reflect.Array {
			if len(mapVarVarkey) > 0 {
				key, e := dc.GetValue(Vars, mapVarVarkey)
				if e != nil {
					return e
				}
				wantedValue, e := core.GetWantedValue(setValue, valueType)
				if e != nil {
					return e
				}
				value.Elem().Index(int(key.Int())).Set(wantedValue)
				return nil
			}

			if len(mapVarStrkey) > 0 {
				return errors.New(fmt.Sprintf("the index of array or slice should not be string, now is str \"%s\"", mapVarStrkey))
			}

			if mapVarIntkey >= 0 {
				wantedValue, e := core.GetWantedValue(setValue, valueType)
				if e != nil {
					return e
				}
				value.Elem().Index(int(mapVarIntkey)).Set(wantedValue)
				return nil
			} else {
				return errors.New("Slice or Array index  must be non-negative!")
			}
		}

	} else {
		newValue := value
		valueType := newValue.Type().Elem()

		if newValue.Kind() == reflect.Map {
			keyType := newValue.Type().Key()
			if len(mapVarVarkey) > 0 {
				key, e := dc.GetValue(Vars, mapVarVarkey)
				if e != nil {
					return e
				}
				wantedKey, e := core.GetWantedValue(key, keyType)
				if e != nil {
					return e
				}

				wantedValue, e := core.GetWantedValue(setValue, valueType)
				if e != nil {
					return e
				}
				value.SetMapIndex(wantedKey, wantedValue)
				return nil
			}

			if len(mapVarStrkey) > 0 {
				wantedValue, e := core.GetWantedValue(setValue, valueType)
				if e != nil {
					return e
				}
				value.SetMapIndex(reflect.ValueOf(mapVarStrkey), wantedValue)
				return nil
			}

			//int key
			wantedKey, e := core.GetWantedValue(reflect.ValueOf(mapVarIntkey), keyType)
			if e != nil {
				return e
			}
			wantedValue, e := core.GetWantedValue(setValue, valueType)
			if e != nil {
				return e
			}
			value.SetMapIndex(wantedKey, wantedValue)
			return nil

		}

		if newValue.Kind() == reflect.Slice || newValue.Kind() == reflect.Array {
			if len(mapVarVarkey) > 0 {
				key, e := dc.GetValue(Vars, mapVarVarkey)
				if e != nil {
					return e
				}
				wantedValue, e := core.GetWantedValue(setValue, valueType)
				if e != nil {
					return e
				}
				value.Index(int(key.Int())).Set(wantedValue)
				return nil
			}

			if len(mapVarStrkey) > 0 {
				return errors.New(fmt.Sprintf("the index of array or slice should not be string, now is str \"%s\"", mapVarStrkey))
			}

			if mapVarIntkey >= 0 {
				wantedValue, e := core.GetWantedValue(setValue, valueType)
				if e != nil {
					return e
				}
				value.Index(int(mapVarIntkey)).Set(wantedValue)
				return nil
			} else {
				return errors.New("Slice or Array index  must be non-negative!")
			}
		}
	}

	return errors.New(fmt.Sprintf("unspport type, mapVarName =%s", mapVarName))
}

func (dc *DataContext) makeArray(value interface{}) {
	//todo
}
