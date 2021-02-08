package test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
)

const n_m_model_rules = `
rule "100" "最高优先级" salience 100 
begin
	println("a")
	data.Count += 2
end

rule "98"  salience 97 
begin
	println("b")
	data.Count += 2
end

rule "90" salience 90
begin
	println("c")
	data.Count += 2
end

rule "89" salience 89
begin
	println("d")	
	println(data.Count)
end

rule "77" salience 85
begin
	println("e")
	s = "hello world"
	println(s)
	return s
end

rule "50"  salience 50 
begin
	println("f")
	println("one more rule")
end


`

//测试代码框架
func nmFramework(n, m, em int, names []string) {
	type Data struct {
		Count int
	}

	data := &Data{Count: 0}
	dataContext := context.NewDataContext()
	dataContext.Add("data", data)
	dataContext.Add("println", fmt.Println)
	ruleBuilder := builder.NewRuleBuilder(dataContext)
	e1 := ruleBuilder.BuildRuleFromString(n_m_model_rules)
	if e1 != nil {
		panic(e1)
	}

	eng := engine.NewGengine()
	var e2 error
	if em == 1 {
		e2 = eng.ExecuteNSortMConcurrent(n, m, ruleBuilder, true)
	}
	if em == 2 {
		e2 = eng.ExecuteNConcurrentMSort(n, m, ruleBuilder, true)
	}
	if em == 3 {
		e2 = eng.ExecuteNConcurrentMConcurrent(n, m, ruleBuilder, true)
	}
	if em == 4 {
		e2 = eng.ExecuteSelectedNSortMConcurrent(n, m, ruleBuilder, true, names)
	}
	if em == 5 {
		e2 = eng.ExecuteSelectedNConcurrentMSort(n, m, ruleBuilder, true, names)
	}
	if em == 6 {
		e2 = eng.ExecuteSelectedNConcurrentMConcurrent(n, m, ruleBuilder, true, names)
	}

	if e2 != nil {
		panic(e2)
	}

	//如果有规则有返回值, 可以这样获取
	//resultMap, _ := eng.GetRulesResultMap()

	//i := resultMap["77"] //获取规则名为77的返回值
	//s := i.(string)
	//println("get return -->", s )
}

func Test_n_sort_m_concurrent(t *testing.T) {
	//ExecuteNSortMConcurrent
	nmFramework(3, 3, 1, []string{})
}

func Test_n_concurrent_m_sort(t *testing.T) {
	//ExecuteNConcurrentMSort
	nmFramework(3, 3, 2, []string{})
}

func Test_n_concurrent_m_concurrent(t *testing.T) {
	//ExecuteNConcurrentMConcurrent
	nmFramework(3, 3, 3, []string{})
}

func Test_selected_n_sort_m_concurrent(t *testing.T) {
	//ExecuteSelectedNSortMConcurrent
	nmFramework(3, 2, 4, []string{"100", "98", "89", "90", "77"})
}

func Test_selected_n_concurrent_m_sort(t *testing.T) {
	//ExecuteSelectedNConcurrentMSort
	nmFramework(3, 3, 5, []string{"100", "98", "89", "90", "77", "50"})
}

func Test_selected_n_concurrent_m_concurrent(t *testing.T) {
	//ExecuteSelectedNConcurrentMConcurrent
	nmFramework(3, 3, 5, []string{"100", "98", "89", "90", "77", "50"})
}
