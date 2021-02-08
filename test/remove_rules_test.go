package test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
)

var rules = `
rule "1" 
begin
println("1")
end
rule "2" 
begin
println("2")
end
rule "3" 
begin
println("3")
end
rule "4" 
begin
println("4")
end
rule "5" 
begin
println("5")
end
rule "6" 
begin
println("6")
end
`

func Test_single_remove_rules(t *testing.T) {

	dataContext := context.NewDataContext()
	dataContext.Add("println", fmt.Println)
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	e := ruleBuilder.BuildRuleFromString(rules)

	if e != nil {
		panic(e)
	}

	isExist := ruleBuilder.IsExist([]string{"1", "2", "3", "4", "5", "6", "7"})
	fmt.Println(fmt.Sprintf("%+v,len=%d", isExist, len(ruleBuilder.Kc.RuleEntities)))

	e = ruleBuilder.RemoveRules([]string{"3", "5"})
	if e != nil {
		panic(e)
	}

	isExist = ruleBuilder.IsExist([]string{"1", "2", "3", "4", "5", "6", "7"})
	fmt.Println(fmt.Sprintf("%+v,len=%d", isExist, len(ruleBuilder.Kc.RuleEntities)))

}



func Test_pool_remove_rules(t *testing.T) {

	apis := make(map[string]interface{})
	pool, e := engine.NewGenginePool(1, 10, 1, rules, apis)
	if e != nil {
		panic(e)
	}

	isExist := pool.IsExist([]string{"1", "2", "3", "4", "5", "6", "7"})
	fmt.Println(fmt.Sprintf("%+v", isExist))

	e = pool.RemoveRules([]string{"3", "5"})
	if e != nil {
		panic(e)
	}
	isExist = pool.IsExist([]string{"1", "2", "3", "4", "5", "6", "7"})
	fmt.Println(fmt.Sprintf("%+v", isExist))
}
