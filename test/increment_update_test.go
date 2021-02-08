package test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
	"time"
)

const all_rules = `
rule "helo" "zzz" salience 100
begin
println("hello")
a = 100
b = 1000
println(a+b)
end

rule "word" "www" salience 50
begin
println("word")
a = 10
b = 9
println(a+b)
end
`

const incremental_replace_rule = `
rule "word" "www" salience 30
begin
println("word+++++")
a = 10
b = 90000
println(a+b)
end
`

const incremental_replace_rule_1 = `
rule "word" "www" salience 60
begin
println("word+++++11")
a = 10
b = 989798700
println(a+b)

end
`

const incremental_add_rule = `
rule "iii" "www" salience 50
begin
println("iii+++++")
a = 10
b = 989111
println(a+b)
end
`

func Test_increment_1(t *testing.T) {

	dataContext := context.NewDataContext()
	dataContext.Add("println", fmt.Println)

	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	//resolve rules from string
	er1 := ruleBuilder.BuildRuleFromString(all_rules)
	if er1 != nil {
		println(fmt.Sprintf("%+v", er1))
	}

	eng := engine.NewGengine()
	e1 := eng.Execute(ruleBuilder, false)
	if e1 != nil {
		panic(e1)
	}

	println("------------------1------------------base")

	er2 := ruleBuilder.BuildRuleWithIncremental(incremental_replace_rule)

	if er2 != nil {
		panic(er2)
	}

	e2 := eng.Execute(ruleBuilder, true)
	if e2 != nil {
		panic(e2)
	}

	println("------------------2------------------replace update")

	er3 := ruleBuilder.BuildRuleWithIncremental(incremental_add_rule)

	if er3 != nil {
		panic(er3)
	}

	e3 := eng.Execute(ruleBuilder, true)
	if e3 != nil {
		panic(e3)
	}

	println("------------------3------------------add update")

	er4 := ruleBuilder.BuildRuleWithIncremental(incremental_replace_rule_1)

	if er4 != nil {
		panic(er4)
	}

	e4 := eng.Execute(ruleBuilder, true)
	if e4 != nil {
		panic(e4)
	}

	println("------------------4------------------differen salience update")

	e5 := eng.Execute(ruleBuilder, true)
	if e5 != nil {
		panic(e5)
	}

	println("------------------5------------------differen salience update")

	//resolve rules from string
	er6 := ruleBuilder.BuildRuleFromString(all_rules)
	if er6 != nil {
		println(fmt.Sprintf("%+v", er1))
	}

	e7 := eng.Execute(ruleBuilder, true)
	if e7 != nil {
		panic(e7)
	}
	println("------------------7------------------all update")

}

func Test_increment_1_concurrent(t *testing.T) {

	dataContext := context.NewDataContext()
	dataContext.Add("println", fmt.Println)

	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	e1 := ruleBuilder.BuildRuleWithIncremental(incremental_replace_rule)
	if e1 != nil {
		panic(e1)
	}

	eng := engine.NewGengine()
	e3 := eng.Execute(ruleBuilder, true)
	if e3 != nil {
		panic(e3)
	}

	go func() {
		for i := 0; i < 1000; i++ {
			e1 := ruleBuilder.BuildRuleFromString(all_rules)
			if e1 != nil {
				panic(e1)
			}
		}
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			e1 := ruleBuilder.BuildRuleWithIncremental(incremental_add_rule)
			if e1 != nil {
				panic(e1)
			}
		}
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			e1 := ruleBuilder.BuildRuleWithIncremental(incremental_replace_rule)
			if e1 != nil {
				panic(e1)
			}
		}
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			e3 := eng.Execute(ruleBuilder, true)
			if e3 != nil {
				panic(e3)
			}
		}
	}()

	time.Sleep(2 * time.Second)

}

func Test_increment_pool(t *testing.T) {

	t1 := time.Now()
	apis := make(map[string]interface{})
	apis["println"] = fmt.Println
	pool, e1 := engine.NewGenginePool(1, 2, 1, all_rules, apis)
	if e1 != nil {
		println(fmt.Sprintf("e1: %+v", e1))
	}
	println("build pool cost time:", time.Since(t1), "ns")
	reqest := &Reqest{}
	data := make(map[string]interface{})
	data["Req"] = reqest

	e2, _ := pool.ExecuteRulesWithMultiInputWithSpecifiedEM(data)
	if e2 != nil {
		panic(e2)
	}
	println("------------------1------------------base")

	er1 := pool.UpdatePooledRulesIncremental(incremental_replace_rule)
	if er1 != nil {
		panic(er1)
	}

	er2, _ := pool.ExecuteRulesWithMultiInputWithSpecifiedEM(data)
	if er2 != nil {
		panic(er2)
	}
	println("------------------2------------------replace update")

	er3 := pool.UpdatePooledRulesIncremental(incremental_add_rule)
	if er3 != nil {
		panic(er3)
	}

	er4, _ := pool.ExecuteRulesWithMultiInputWithSpecifiedEM(data)
	if er4 != nil {
		panic(er4)
	}

	println("------------------3------------------add update")

	er5 := pool.UpdatePooledRulesIncremental(incremental_replace_rule_1)
	if er5 != nil {
		panic(er5)
	}

	er6, _ := pool.ExecuteRulesWithMultiInputWithSpecifiedEM(data)
	if er6 != nil {
		panic(er6)
	}

	println("------------------4------------------differen salience update")

	/*	er7 := pool.UpdatePooledRulesIncremental(incremental_replace_rule_1)
		if er7!=nil {
			panic(er5)
		}*/

	er8, _ := pool.ExecuteRulesWithMultiInputWithSpecifiedEM(data)
	if er8 != nil {
		panic(er6)
	}

	println("------------------5------------------differen salience update")
}

func Test_increment_pool_concurrent(t *testing.T) {

	t1 := time.Now()
	apis := make(map[string]interface{})
	apis["println"] = fmt.Println
	pool, e1 := engine.NewGenginePool(1, 2, 1, incremental_replace_rule_1, apis)
	if e1 != nil {
		println(fmt.Sprintf("e1: %+v", e1))
	}
	println("build pool cost time:", time.Since(t1), "ns")
	reqest := &Reqest{}
	data := make(map[string]interface{})
	data["Req"] = reqest

	go func() {
		for i := 0; i < 10000; i++ {
			er2, _ := pool.ExecuteRulesWithMultiInputWithSpecifiedEM(data)
			if er2 != nil {
				panic(er2)
			}
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			er5 := pool.UpdatePooledRulesIncremental(incremental_replace_rule)
			if er5 != nil {
				panic(er5)
			}
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			er5 := pool.UpdatePooledRules(all_rules)
			if er5 != nil {
				panic(er5)
			}
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			er5 := pool.UpdatePooledRulesIncremental(incremental_add_rule)
			if er5 != nil {
				panic(er5)
			}
		}
	}()

	time.Sleep(10 * time.Second)

}
