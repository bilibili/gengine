package test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
)

const rules_execute string = `
rule "1" "1"
begin
println("----1-------")
end

rule "2" "2"
begin
println("----2-------")
end

rule "3" "3"
begin
println("----3-------")
end

rule "4" "4"
begin
println("----4-------")
end
`

func Test_execute(t *testing.T) {

	dataContext := context.NewDataContext()
	dataContext.Add("println", fmt.Println)

	ruleBuilder := builder.NewRuleBuilder(dataContext)
	e := ruleBuilder.BuildRuleFromString(rules_execute)
	if e != nil {
		panic(fmt.Sprintf("build rules err:%+v", e))
	}

	gengine := engine.NewGengine()
	println("==================0=================")
	e1 := gengine.Execute(ruleBuilder, false)
	if e1 != nil {
		panic(fmt.Sprintf("Execute:%+v", e1))
	}
	println("==================1=================")

	sTag := &engine.Stag{StopTag: false}
	e2 := gengine.ExecuteWithStopTagDirect(ruleBuilder, true, sTag)
	if e2 != nil {
		panic(fmt.Sprintf("ExecuteWithStopTagDirect:%+v", e2))
	}
	println("==================2=================")

	e3 := gengine.ExecuteConcurrent(ruleBuilder)
	if e3 != nil {
		panic(fmt.Sprintf("ExecuteConcurrent:%+v", e3))
	}
	println("==================3=================")

	e4 := gengine.ExecuteMixModel(ruleBuilder)
	if e4 != nil {
		panic(fmt.Sprintf("ExecuteMixModel:%+v", e4))
	}
	println("==================4=================")

	sTag1 := &engine.Stag{StopTag: false}
	e5 := gengine.ExecuteMixModelWithStopTagDirect(ruleBuilder, sTag1)
	if e5 != nil {
		panic(fmt.Sprintf("ExecuteMixModelWithStopTagDirect:%+v", e5))
	}
	println("==================5=================")

	e6 := gengine.ExecuteSelectedRules(ruleBuilder, []string{"1", "2", "6"})
	if e6 != nil {
		panic(fmt.Sprintf("ExecuteSelectedRules:%+v", e6))
	}
	println("==================6=================")

	e7 := gengine.ExecuteSelectedRulesWithControl(ruleBuilder, true, []string{"1", "2", "3"})
	if e7 != nil {
		panic(fmt.Sprintf("ExecuteSelectedRulesWithControl:%+v", e7))
	}

	println("==================7=================")
	sTag2 := &engine.Stag{StopTag: false}
	e8 := gengine.ExecuteSelectedRulesWithControlAndStopTag(ruleBuilder, true, sTag2, []string{"1", "2", "3", "4"})
	if e8 != nil {
		panic(fmt.Sprintf("ExecuteSelectedRulesWithControlAndStopTag:%+v", e8))
	}
	println("==================8=================")
	e9 := gengine.ExecuteSelectedRulesConcurrent(ruleBuilder, []string{"1", "2", "3", "4"})
	if e9 != nil {
		panic(fmt.Sprintf("ExecuteSelectedRulesConcurrent:%+v", e9))
	}
	println("==================9=================")

	e10 := gengine.ExecuteSelectedRulesMixModel(ruleBuilder, []string{"1", "2", "3", "4", "5"})
	if e10 != nil {
		panic(fmt.Sprintf("ExecuteSelectedRulesMixModel:%+v", e10))
	}

	println("===================10================")

	e11 := gengine.ExecuteSelectedRulesWithControlAsGivenSortedName(ruleBuilder, true, []string{"3", "2", "1", "4"})
	if e11 != nil {
		panic(fmt.Sprintf("ExecuteSelectedRulesWithControlSortByName:%+v", e11))
	}
	println("===================11================")

	sTag3 := &engine.Stag{StopTag: false}
	e12 := gengine.ExecuteSelectedRulesWithControlAndStopTagAsGivenSortedName(ruleBuilder, true, sTag3, []string{"2", "1", "4", "3"})
	if e12 != nil {
		panic(fmt.Sprintf("ExecuteSelectedRulesWithControlAndStopTagSortByNames:%+v", e12))
	}
	println("===================12================")
}
