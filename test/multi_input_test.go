package test

import (
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
)

const multi_input_rule = `
rule "1" ""
begin
println("hello", "world")
end
`

func Println(as ...string) {
	for _, v := range as {
		print(v, " ")
	}
}

func Test_multi_input_(t *testing.T) {

	dataContext := context.NewDataContext()
	//inject struct
	//rename and inject

	dataContext.Add("println", Println)

	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	//读取规则
	e1 := ruleBuilder.BuildRuleFromString(multi_input_rule)
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
