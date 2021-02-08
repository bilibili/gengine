package test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
	"time"
)

type Person struct {
	Name string
	Age  int64
}

func getPerson(n string, a int64) *Person {
	return &Person{
		Name: n,
		Age:  a,
	}
}

type Req struct {
	//Data string
}

func GetPool(req *Req) {

	println("hello....")
}

/*func Test_Struct(t *testing.T)  {
	 p := getPerson("777", 5)
	 println(p.Age)
}*/

const rule_s = `
rule "test_struct_return" "test" 
begin
	//p = getPerson("hello", 100)
	//Sout(p.Age)
	GetPool(req)
end
`

func exe_struct() {
	/**
	不要注入除函数和结构体指针以外的其他类型(如变量)
	*/
	dataContext := context.NewDataContext()
	//inject struct
	//rename and inject

	req := &Req{}
	dataContext.Add("Sout", fmt.Println)
	dataContext.Add("getPerson", getPerson)
	dataContext.Add("GetPool", GetPool)
	dataContext.Add("req", req)

	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	//读取规则
	start1 := time.Now().UnixNano()
	err := ruleBuilder.BuildRuleFromString(rule_s)
	end1 := time.Now().UnixNano()

	println(fmt.Sprintf("rules num:%d, load rules cost time:%d ns", len(ruleBuilder.Kc.RuleEntities), end1-start1))

	if err != nil {
		println(fmt.Sprintf("err:%s ", err))
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

func Test_struct(t *testing.T) {
	exe_struct()
}
