package test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
)

const expr_params_rule = `
rule "name test" "i can" salience 10
begin

x = 100
z = 200
isOld = x > 10 + 20

println(isOld)

println(!true)
y = false

ExpParam(x +20, !y)

end
`

func ExpParam(i int, b bool) {
	println("input:", i, b)
}

func Test_expr_params(t *testing.T) {

	dataContext := context.NewDataContext()
	dataContext.Add("println", fmt.Println)
	dataContext.Add("ExpParam", ExpParam)

	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	//读取规则
	e1 := ruleBuilder.BuildRuleFromString(expr_params_rule)
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
