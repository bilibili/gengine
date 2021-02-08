package test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
	"time"
)

/**
use '@desc',you can get rule description in rule content
*/
const atDescRule = `
rule "rule name 1" "我是一个测试用的描述信息1" salience 100
begin
  desc = @desc
  Print(desc)
  Print(@name + " : " + @desc)
end

rule "rule name 2" //"我是描述，desc" salience 10
begin
  desc = @desc
  Print(desc)
  Print(@name + " : " + @desc)
end
`

func Test_AtDesc(t *testing.T) {
	dataContext := context.NewDataContext()
	// dataContext.Add("Print", PrintName)
	dataContext.Add("Print", fmt.Println)

	//init rule engine

	ruleBuilder := builder.NewRuleBuilder(dataContext)

	//resolve rules from string
	start1 := time.Now().UnixNano()
	err := ruleBuilder.BuildRuleFromString(atDescRule)
	end1 := time.Now().UnixNano()

	println(fmt.Sprintf("rules num:%d, load rules cost time:%d ns", len(ruleBuilder.Kc.RuleEntities), end1-start1))

	if err != nil {
		panic(err)
	}
	eng := engine.NewGengine()
	// true: means when there are many rules， if one rule execute error，continue to execute rules after the occur error rule
	err = eng.Execute(ruleBuilder, true)
	if err != nil {
		panic(err)
	}

}
