package map_slice_array

import (
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
)

const return_asm_rule string = `
rule "asm" "test"
begin


MyValue.M=GetMap()
MyValue.A = GetArray()
MyValue.S = GetSlice()
MyValue.SC = GetStruct()
MyValue.In = GetInterface()
MyValue.Ch = GetChan()
MyValue.Cp = GetComplex128()
MyValue.Fu = GetFunc()
end
`

type MyValue struct {
	M  map[string]string
	A  [4]int64
	S  []string
	SC MyStruct
	In interface{}
	Ch chan string
	Cp complex64
	Fu func(int)
}

type MyStruct struct {
}

func GetMap() map[string]string {
	return make(map[string]string)
}

func GetArray() [4]int64 {
	return [4]int64{1, 3, 4, 9}
}

func GetSlice() []string {
	return []string{"hello", "world"}
}

func GetStruct() MyStruct {
	return MyStruct{}
}

func GetInterface() interface{} {

	return 10
}

func GetChan() chan string {
	ch := make(chan string, 10)
	return ch
}

func GetComplex64() complex64 {
	return complex(1, 3)
}

func GetComplex128() complex128 {
	return complex(1, 3)
}

func GetFunc() func(int) {
	return func(i int) {

	}
}

func Test_return_array_map_slice(t *testing.T) {

	dc := context.NewDataContext()
	dc.Add("GetMap", GetMap)
	dc.Add("GetArray", GetArray)
	dc.Add("GetSlice", GetSlice)
	dc.Add("GetStruct", GetStruct)
	dc.Add("MyValue", &MyValue{})
	dc.Add("GetInterface", GetInterface)
	dc.Add("GetChan", GetChan)
	dc.Add("GetComplex64", GetComplex64)
	dc.Add("GetComplex128", GetComplex128)
	dc.Add("GetFunc", GetFunc)
	rb := builder.NewRuleBuilder(dc)
	e := rb.BuildRuleFromString(return_asm_rule)
	if e != nil {
		panic(e)
	}
	gengine := engine.NewGengine()

	e = gengine.Execute(rb, true)
	if e != nil {
		panic(e)
	}
}
