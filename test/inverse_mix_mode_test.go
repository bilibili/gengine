package test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
	"time"
)

const inverse_rules = `
rule "1000" "most priority" salience 1000
begin
cal.Data = 5

end

rule "998" "lower priority" salience 998
begin
cal.Name = "hello world"
end

rule "996" "lowest priority" salience 996
begin
println("cal.Data-->", cal.Data)
println("cal.Name-->", cal.Name) 
end
`

type Calculate struct {
	Data int
	Name string
}

func Test_inverse(t *testing.T) {

	dataContext := context.NewDataContext()
	//inject struct
	calculate := &Calculate{Data: 0}
	dataContext.Add("cal", calculate)
	//rename and inject
	dataContext.Add("println", fmt.Println)

	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	//读取规则
	e1 := ruleBuilder.BuildRuleFromString(inverse_rules)
	if e1 != nil {
		panic(e1)
	}

	eng := engine.NewGengine()

	e2 := eng.ExecuteInverseMixModel(ruleBuilder)
	if e2 != nil {
		panic(e2)
	}
}

func Test_inverse_with_selected(t *testing.T) {

	dataContext := context.NewDataContext()
	//inject struct
	calculate := &Calculate{Data: 0}
	dataContext.Add("cal", calculate)
	//rename and inject
	dataContext.Add("println", fmt.Println)

	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	//读取规则
	e1 := ruleBuilder.BuildRuleFromString(inverse_rules)
	if e1 != nil {
		panic(e1)
	}

	eng := engine.NewGengine()

	e2 := eng.ExecuteSelectedRulesInverseMixModel(ruleBuilder, []string{"996", "1000"})
	if e2 != nil {
		panic(e2)
	}
}

func Test_inverse_pool(t *testing.T) {

	t1 := time.Now()
	apis := make(map[string]interface{})
	apis["println"] = fmt.Println
	pool, e1 := engine.NewGenginePool(1, 2, 4, inverse_rules, apis)
	if e1 != nil {
		println(fmt.Sprintf("e1: %+v", e1))
	}
	println("build pool cost time:", time.Since(t1), "ns")
	calculate := &Calculate{Data: 0}
	data := make(map[string]interface{})
	data["cal"] = calculate

	e2, _ := pool.ExecuteInverseMixModel(data)
	if e2 != nil {
		panic(e2)
	}
}

func Test_inverse_select_pool(t *testing.T) {

	t1 := time.Now()
	apis := make(map[string]interface{})
	apis["println"] = fmt.Println
	pool, e1 := engine.NewGenginePool(1, 2, 4, inverse_rules, apis)
	if e1 != nil {
		println(fmt.Sprintf("e1: %+v", e1))
	}
	println("build pool cost time:", time.Since(t1), "ns")
	calculate := &Calculate{Data: 0}
	data := make(map[string]interface{})
	data["cal"] = calculate

	e2, _ := pool.ExecuteSelectedRulesInverseMixModel(data, []string{"996", "1000"})
	if e2 != nil {
		panic(e2)
	}
}

func Test_select_pool(t *testing.T) {

	t1 := time.Now()
	apis := make(map[string]interface{})
	apis["println"] = fmt.Println
	pool, e1 := engine.NewGenginePool(1, 2, 4, inverse_rules, apis)
	if e1 != nil {
		println(fmt.Sprintf("e1: %+v", e1))
	}
	println("build pool cost time:", time.Since(t1), "ns")
	calculate := &Calculate{Data: 0}
	data := make(map[string]interface{})
	data["cal"] = calculate

	e2, _ := pool.ExecuteSelectedWithSpecifiedEM(data, []string{"996", "1000"})
	if e2 != nil {
		panic(e2)
	}
}

func Test_base_pool(t *testing.T) {

	t1 := time.Now()
	apis := make(map[string]interface{})
	apis["println"] = fmt.Println
	pool, e1 := engine.NewGenginePool(1, 2, 4, inverse_rules, apis)
	if e1 != nil {
		println(fmt.Sprintf("e1: %+v", e1))
	}
	println("build pool cost time:", time.Since(t1), "ns")
	calculate := &Calculate{Data: 0}
	data := make(map[string]interface{})
	data["cal"] = calculate

	//e2 := pool.ExecuteRulesWithMultiInputAndStopTag(data, &engine.Stag{StopTag:false})
	e2, _ := pool.ExecuteRulesWithMultiInputWithSpecifiedEM(data)
	if e2 != nil {
		panic(e2)
	}
}
