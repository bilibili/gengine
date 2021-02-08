package test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
)

func Test_input_array(t *testing.T) {

	var rule = `
rule "test" "test" 
begin
	if isNil(InputAndResult.Result["tenant_id"]) {
       InputAndResult.AddStringArray("data",InputAndResult.GetSlice("123","234"))
    }
end
`

	dataContext := context.NewDataContext()
	//init rule engine
	InputAndResult := &InputAndResult{Result: make(map[string][]string)}
	dataContext.Add("InputAndResult", InputAndResult)

	ruleBuilder := builder.NewRuleBuilder(dataContext)

	//读取规则
	e1 := ruleBuilder.BuildRuleFromString(rule)
	if e1 != nil {
		panic(e1)
	}

	eng := engine.NewGengine()
	// true: means when there are many rules， if one rule execute error，continue to execute rules after the occur error rule
	e2 := eng.Execute(ruleBuilder, true)
	if e2 != nil {
		panic(e2)
	}

	//输出在规则中输入的内容
	println(fmt.Sprintf("%+v", InputAndResult.Result["data"]))

}

type InputAndResult struct {
	Result map[string][]string
}

func (input *InputAndResult) AddStringArray(key string, value []string) {
	input.Result[key] = value
}

func (input *InputAndResult) GetSlice(vs ...string) []string {
	var s []string
	for _, v := range vs {
		s = append(s, v)
	}
	return s
}
