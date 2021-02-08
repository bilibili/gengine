package test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
)

func Test_createBaseline(t *testing.T) {

	rule := `rule "3"  "1.1.1-b,1.1-c" salience 5 
begin
	i = P.ProjectInformation["ProjectType"]
	println("map is nil--->", i)
	if P.ProjectInformation["ProjectType"] =="MIS系统,内部平台"{
		return @desc
	}
end
`

	dataContext := context.NewDataContext()
	dataContext.Add("getInt", getInt)
	//input要带上括号
	dataContext.Add("P", input())
	dataContext.Add("println", fmt.Println)

	ruleBuilder := builder.NewRuleBuilder(dataContext)
	e1 := ruleBuilder.BuildRuleFromString(rule)
	if e1 != nil {
		panic(e1)
	}

	eng := engine.NewGengine()
	e2 := eng.Execute(ruleBuilder, true)
	if e2 != nil {
		println(e2)
	}

	resultMap, _ := eng.GetRulesResultMap()

	r, ok := resultMap["3"]
	if ok {
		println("return ---->", r.(string))
	}

}

type P struct {
	ProjectInformation map[string]string
}

func input() *P {
	return &P{ProjectInformation: map[string]string{
		"ProjectType":   "MIS系统,内部平台",
		"TargetUser":    "合作伙伴,政府",
		"ServerInvoker": "前段页面,手机H5,手机APP",
	}}
}
