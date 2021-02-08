package test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"reflect"
	"testing"
)

//in golang, you can't println(nil)
const nil_test_rule = `
rule "1" "test return nil" 
begin
s = live.GetStringPtr()
live.SetStringPtr(s)

s1 = live.GetString()

b = live.GetBoolPtr()
live.SetBoolPtr(b)

i = live.GetIntPtr()
live.SetIntPtr(i)

u = live.GetUintPtr()
live.SetUintPtr(u)

f = live.GetFloatPtr()
live.SetFloatPtr(f)

e = live.GetEverPtr()
live.SetEverPtr(e)

ap = live.GetArrayPtr()
live.SetArrayPtr(ap)

a = live.GetArray()


sp = live.GetSlicePtr()
live.SetSlicePtr(sp)

sl = live.GetSlice()
live.SetSlice(sl)

mp = live.GetMapPtr()
live.SetMapPtr(mp)

m = live.GetMap()
live.SetMap(m)

c = live.GetChan()
live.SetChan(c)

ff = live.Setfun()
live.Getfun(ff)

ii = live.GetInterf()
live.SetInterf(ii)
//live.x()
end
`

type Live struct {
}

type Ever struct {
}

func (l *Live) GetStringPtr() *string {
	println("GetStringPtr")
	return nil
}

func (l *Live) SetStringPtr(s *string) {
	println("SetStringPtr... ")
}

//can't return nil
/*func (l *Live)GetString() string {
	return nil
}
*/

func (l *Live) GetString() string {
	var s string
	println("GetString")
	return s
}

func (l *Live) GetBoolPtr() *bool {
	println("GetBoolPtr")
	return nil
}

func (l *Live) SetBoolPtr(b *bool) {
	if b == nil {
		println("----------b is nil-------")
	}
	println("SetBoolPtr...")
}

//can't return nil
/*func (l *Live)GetBool() bool {
	return nil
}*/

func (l *Live) GetIntPtr() *int {
	println("GetIntPtr")
	return nil
}

func (l *Live) SetIntPtr(i *int) {
	println("SetIntPtr.....")
}

//can't return nil
/*func (l *Live)GetInt() int {
	return nil
}
*/

func (l *Live) GetUintPtr() *uint {
	println("GetUintPtr")
	return nil
}

func (l *Live) SetUintPtr(u *uint) {
	println("SetUintPtr...")
}

//can't return nil
/*func (l *Live)GetUint() uint {
	return nil
}*/

func (l *Live) GetFloatPtr() *float64 {
	println("GetFloatPtr")
	return nil
}

func (l *Live) SetFloatPtr(f *float64) {
	println("SetFloatPtr")
}

//can't return nil
/*func (l *Live)GetFloat() float64 {
	return nil
}*/

func (l *Live) GetEverPtr() *Ever {
	println("GetEverPtr")
	return nil
}

func (l *Live) SetEverPtr(e *Ever) {
	println("SetEverPtr...")
}

//can't return nil
/*func (l *Live)GetEver() Ever {
	return nil
}*/

func (l *Live) GetArrayPtr() *[4]int64 {
	println("GetArrayPtr")
	return nil
}

func (l *Live) SetArrayPtr(x *[4]int64) {
	println("SetArrayPtr....")
}

//can't return nil
/*func (l *Live) GetArray() [4]int64 {
	return nil
}*/

func (l *Live) GetArray() [4]int64 {
	var x [4]int64
	println("GetArray")
	return x
}

func (l *Live) GetSlicePtr() *[]int64 {
	println("GetSlicePtr")
	return nil
}

func (l *Live) SetSlicePtr(s *[]int64) {
	println("SetSlicePtr....")
}

func (l *Live) GetSlice() []int64 {
	println("GetSlice")
	return nil
}

func (l *Live) SetSlice(s []int64) {
	println("SetSlice....")
}

func (l *Live) GetMapPtr() *map[int64]string {
	println("GetMapPtr")
	return nil
}

func (l *Live) SetMapPtr(m *map[int64]string) {
	println("setMapPtr....")
}

func (l *Live) GetMap() map[int64]string {
	println("GetMap")
	return nil
}

func (l *Live) SetMap(m map[int64]string) {
	println("SetMap....")
}

func (l *Live) GetChan() chan string {
	println("GetChan")
	return nil
}

func (l *Live) SetChan(c chan string) {
	println("SetChan....")
}

func (l *Live) Setfun() func() {
	println("Setfun")
	return nil
}

func (l *Live) Getfun(f func()) {

	println("Getfun...")
}

func (l *Live) GetInterf() interface{} {
	println("GetInterf...")
	return nil
}

func (l *Live) SetInterf(i interface{}) {
	if i == nil {
		println("interface is nil")
		return
	}
	println("SetInterf...", reflect.ValueOf(i).Int())
}

func Test_return_nil(t *testing.T) {

	dataContext := context.NewDataContext()
	//inject struct
	//dataContext.Add("User", user)
	//rename and inject
	dataContext.Add("println", fmt.Println)
	dataContext.Add("live", &Live{})

	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	//读取规则
	e1 := ruleBuilder.BuildRuleFromString(nil_test_rule)
	if e1 != nil {
		panic(e1)
	}

	eng := engine.NewGengine()
	// true: means when there are many rules， if one rule execute error，continue to execute rules after the occur error rule
	e2 := eng.Execute(ruleBuilder, true)
	if e2 != nil {
		panic(e2)
	}

}
