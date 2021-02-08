package map_slice_array

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"github.com/bilibili/gengine/internal/core"
	"testing"
)

var sliRules = `
rule "1"
begin

//基础类型int，不存在返回默认值
a = contextInt[1]
println("a---> ", a)
println("int isNil  a--->", isNil(a)) // 基础类型永远是false
println("--------------------------------------------------")

//基础类型string，不存在返回默认值
B = 1
b = contextStr[B] //8?
println("b---> ", b)
println("string isNil  b--->",isNil(b)) //基础类型永远是false
println("--------------------------------------------------")

//非指针结构体
x = 0 // 试试 x =3 
c = contextObj[x]
println("c---> ", c)
println("struct isNil  c--->",isNil(c))
println("--------------------------------------------------")


//ptr value值本为nil
d = contextPtr[1]
println("d---> ", d)
println("ptr struct isNil d--->",isNil(d))
println("--------------------------------------------------")

//ptr value值不为nil
e = contextPtr[2]
println("e---> ", e)
println("ptr struct isNil e--->",isNil(e))
println("--------------------------------------------------")

//ptr value值不存在的值
f = contextPtr[0]
println("f---> ", f)
println("ptr struct isNil f--->",isNil(f))
println("--------------------------------------------------")


//bool值不存在的值
g = contextBool[0] //1 ?
println("g---> ", g)
println("ptr struct isNil g--->",isNil(g))
println("--------------------------------------------------")

//map存在
h = contextMap[2] 
println("h---> ", h)
println("ptr struct isNil h--->",isNil(h))
println("--------------------------------------------------")

//map不存在
i = contextMap[8] 
println("i---> ", i)
println("ptr struct isNil i--->",isNil(i))
println("--------------------------------------------------")

end
`

func Test_slice_nil_value(t *testing.T) {

	type Request struct {
		S string
		I int
	}

	//struct
	contextObj := make([]Request, 8)
	contextObj[0] = Request{
		S: "a", //此处有值，和无值，最终判别nil是不一样的!!!
	}

	//ptr
	contextPtr := make([]*Request, 4)
	contextPtr[1] = nil
	contextPtr[2] = &Request{}

	//map
	contextMap := make([]map[string]string, 9)
	contextMap[1] = nil
	contextMap[2] = make(map[string]string)

	//bool
	contextBool := make([]bool, 7)
	contextBool[1] = true

	//int
	contextInt := make([]int, 9)
	contextInt[1] = 0

	//string
	contextStr := make([]string, 8)

	dataContext := context.NewDataContext()
	//结构体
	dataContext.Add("contextObj", &contextObj)
	//指针
	dataContext.Add("contextPtr", &contextPtr)
	//map
	dataContext.Add("contextMap", &contextMap)
	//基础类型值
	dataContext.Add("contextInt", &contextInt)
	dataContext.Add("contextStr", &contextStr)
	dataContext.Add("contextBool", &contextBool)
	dataContext.Add("println", fmt.Println)
	dataContext.Add("isNil", core.IsNil)

	ruleBuilder := builder.NewRuleBuilder(dataContext)
	err := ruleBuilder.BuildRuleFromString(sliRules)
	if err != nil {
		panic(err)
	}

	eng := engine.NewGengine()
	err = eng.Execute(ruleBuilder, true)
	if err != nil {
		panic(err)
	}

}
