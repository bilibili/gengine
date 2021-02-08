package test

import (
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
)

const pass_status_rule = `
rule "1" "1"
begin
rs.V = rs.V +1

end

rule "2" "1"
begin
rs.V = rs.V +1
end

rule "3" "1"
begin
rs.V = rs.V +1
end

rule "4" "1"
begin
rs.V = rs.V +1
end

rule "5" "1"
begin
rs.V = rs.V +1
end

`

type Result struct {
	V int
}

//示例： 首先 初始化Result.V=0，然后在每个规则中，对Result.V加1，打印Result.V的最终结果
func Test_pass_status(t *testing.T) {
	dataContext := context.NewDataContext()
	//inject struct
	rs := &Result{V: 0}
	dataContext.Add("rs", rs)

	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	//读取规则
	e1 := ruleBuilder.BuildRuleFromString(pass_status_rule)
	if e1 != nil {
		panic(e1)
	}

	eng := engine.NewGengine()
	e2 := eng.Execute(ruleBuilder, true)
	if e2 != nil {
		panic(e2)
	}

	println("rs.V-->", rs.V)

}
