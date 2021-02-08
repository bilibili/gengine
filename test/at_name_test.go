package test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
	"time"
)

func PrintName(name string) {
	fmt.Println(name)
}

/**
use '@name',you can get rule name in rule content
*/
const atname_rule = `
rule "测试规则名称1" "rule desc"
begin
  va = @name
  PrintName(va)
  PrintName(@name)
end
rule "rule name" "rule desc"
begin
  va = @name
  PrintName(va)
  PrintName(@name)
end
`

func exec() {
	dataContext := context.NewDataContext()
	dataContext.Add("PrintName", PrintName)

	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	//resolve rules from string
	start1 := time.Now().UnixNano()
	err := ruleBuilder.BuildRuleFromString(atname_rule)
	end1 := time.Now().UnixNano()

	println(fmt.Sprintf("rules num:%d, load rules cost time:%d ns", len(ruleBuilder.Kc.RuleEntities), end1-start1))

	if err != nil {
		panic(err)
	}
	eng := engine.NewGengine()

	start := time.Now().UnixNano()
	// true: means when there are many rules， if one rule execute error，continue to execute rules after the occur error rule
	err = eng.Execute(ruleBuilder, true)
	end := time.Now().UnixNano()
	if err != nil {
		panic(err)
	}
	println(fmt.Sprintf("execute rule cost %d ns", end-start))

}

func Test_AtName(t *testing.T) {
	exec()
}
