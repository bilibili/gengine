package iter

import (
	"fmt"
	"reflect"
)

// Iteration 迭代器
type Iteration interface {
	Next() bool
	Key() reflect.Value
}

type sliceIter struct {
	max int
	cur int
}

func (sl *sliceIter) Next() bool {
	return sl.cur < sl.max
}

func (sl *sliceIter) Key() reflect.Value {
	value := sl.cur
	sl.cur++
	return reflect.ValueOf(value)
}

type mapIter struct {
	keys []reflect.Value
	cur  int
}

func (mp *mapIter) Next() bool {
	return mp.cur < len(mp.keys)
}

func (mp *mapIter) Key() reflect.Value {
	if !mp.Next() {
		return reflect.ValueOf(nil)
	}
	value := mp.keys[mp.cur]
	mp.cur++
	return value
}

type dmIter struct {
	keys []reflect.Value
	cur  int
}

func (dm *dmIter) Next() bool {
	return dm.cur < len(dm.keys)
}

func (dm *dmIter) Key() reflect.Value {
	if !dm.Next() {
		return reflect.ValueOf(nil)
	}
	value := dm.keys[dm.cur]
	dm.cur++
	return value
}

// NewIteration 新建迭代器
func NewInter(value reflect.Value) (Iteration, error) {
	if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
		return &sliceIter{max: value.Len()}, nil
	}
	if value.Kind() == reflect.Map {
		return &mapIter{keys: value.MapKeys()}, nil
	}
	return nil, fmt.Errorf("value kind: %v, type: %v is not iterable", value.Kind(), value.Type())
}
