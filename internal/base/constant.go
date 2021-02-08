package base

import (
	"github.com/bilibili/gengine/context"
	"reflect"
)

type Constant struct {
	ConstantValue reflect.Value
}

func (cons *Constant) AcceptString(str string) error {
	cons.ConstantValue = reflect.ValueOf(str)
	return nil
}

func (cons *Constant) Evaluate(dc *context.DataContext, Vars map[string]reflect.Value) (reflect.Value, error) {
	return cons.ConstantValue, nil
}

func (cons *Constant) AcceptInteger(i64 int64) error {
	cons.ConstantValue = reflect.ValueOf(i64)
	return nil
}

//receive rule's name
func (cons *Constant) AcceptName(name string) error {
	cons.ConstantValue = reflect.ValueOf(name)
	return nil
}

func (cons *Constant) AcceptId(id int64) error {
	cons.ConstantValue = reflect.ValueOf(id)
	return nil
}

//receive rule's description
func (cons *Constant) AcceptDesc(desc string) error {
	cons.ConstantValue = reflect.ValueOf(desc)
	return nil
}

func (cons *Constant) AcceptSalience(sal int64) error {
	cons.ConstantValue = reflect.ValueOf(sal)
	return nil
}
