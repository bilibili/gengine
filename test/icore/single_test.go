package icore

import (
	"reflect"
	"testing"
)

type Req struct {
	//Data string
}

func GetPool(req *Req) {

	println("succ")
}

func Test_vz(t *testing.T) {

	var x interface{}
	x = GetPool

	req := &Req{}
	var y interface{}
	y = req

	println(reflect.TypeOf(y).Kind().String())
	v := reflect.ValueOf(x)
	args := make([]reflect.Value, 0)

	values := append(args, reflect.ValueOf(y))

	v.Call(values)
}
