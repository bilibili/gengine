package icore

import (
	"github.com/bilibili/gengine/context"
	"reflect"
	"testing"
)

type P struct {
	Mx map[string]string
	Ax [3]int
	Sx []int
}

func Test_in_DataContext(t *testing.T) {
	dataContext := context.NewDataContext()

	var vm = make(map[string]reflect.Value)
	//var mm = make(map[string]string)
	//vm["mx"] = mm

	p := &P{
		Mx: make(map[string]string),
		Ax: [3]int{1, 2, 3},
		Sx: []int{7, 8, 9},
	}
	p.Mx["hello"] = "world"
	dataContext.Add("P", p)

	value_1, e := dataContext.GetValue(vm, "P.Mx")
	if e != nil {
		panic(e)
	}

	//map
	println(reflect.TypeOf(value_1).String())
	println(value_1.MapIndex(reflect.ValueOf("hello")).String())
	println("---------------")

	//array
	value_2, e := dataContext.GetValue(vm, "P.Ax")
	if e != nil {
		panic(e)
	}
	println(value_2.String())
	println(value_2.Index(1).Int())
	println("---------------")

	//slice
	value_3, e := dataContext.GetValue(vm, "P.Sx")
	if e != nil {
		panic(e)
	}
	println(value_3.String())
	println(value_3.Index(1).Int())
	println("---------------")

}

func Test_in_Var(t *testing.T) {

	dataContext := context.NewDataContext()

	var vm = make(map[string]reflect.Value)
	vm["mx"] = reflect.ValueOf(map[string]string{"hello": "world"})
	vm["ax"] = reflect.ValueOf([3]int{1, 2, 3})
	vm["sx"] = reflect.ValueOf([]int{7, 8, 9})
	//dataContext.Add("P", )

	value_1, e := dataContext.GetValue(vm, "mx")
	if e != nil {
		panic(e)
	}

	//map
	println(value_1.String())
	println(value_1.MapIndex(reflect.ValueOf("hello")).String())
	println("---------------")

	//array
	value_2, e := dataContext.GetValue(vm, "ax")
	if e != nil {
		panic(e)
	}
	println(value_2.String())
	println(value_2.Index(1).Int())
	println("---------------")

	//slice
	value_3, e := dataContext.GetValue(vm, "sx")
	if e != nil {
		panic(e)
	}
	println(value_3.String())
	println(value_3.Index(1).Int())
	println("---------------")

	var x int
	println(x)

}
