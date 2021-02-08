package test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
)

//测试规则报错时支持行号
const line_number_rule = `
rule "aaa" "test line number support when execute error" 
begin
//SetNumber(101)
AAA.SetNumber(s)
a = AAA.SetNumber(SetNumber(101), 888)
//println(a)
end
`

func (a *AAA) SetNumber(i int64) int64 {
	if i == 100 {
		panic("panic i=100")
	}

	return i
}

type AAA struct {
	B int64
	C []int64
}

func (a *AAA) GetC() *CCC {
	return nil
}

type CCC struct {
}

func Test_line_number(t *testing.T) {
	dataContext := context.NewDataContext()
	//inject struct
	//rename and inject
	A := &AAA{}
	dataContext.Add("println", fmt.Println)
	dataContext.Add("AAA", A)

	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	//读取规则
	e1 := ruleBuilder.BuildRuleFromString(line_number_rule)
	if e1 != nil {
		panic(e1)
	}

	eng := engine.NewGengine()
	// true: means when there are many rules， if one rule execute error，continue to execute rules after the occur error rule
	e2 := eng.Execute(ruleBuilder, true)
	if e2 != nil {
		println(fmt.Sprintf("%+v", e2))
	}

}
