package test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
)

const set_single_value_rule = `
rule "single_value_test"
begin

//px =111

pi = 100
pi = piv

pf = 101.7
pu = 10

pb = false
pb = pb
pb = pbv

ps = "my world"
ps = psv

st1 =ss
st2 = in.S

end
`

type Single struct {
	Hello string
}

type In struct {
	S Single
}

func Test_single_value_set(t *testing.T) {
	dataContext := context.NewDataContext()
	dataContext.Add("println", fmt.Println)

	//反例，不可设置
	//var px int
	//dataContext.Add("px", px)

	//正例，都可设置
	var pi int
	dataContext.Add("pi", &pi)
	dataContext.Add("piv", 10)

	var pf float32
	dataContext.Add("pf", &pf)

	var pu uint
	dataContext.Add("pu", &pu)

	var pb bool
	pb = false
	dataContext.Add("pb", &pb)
	dataContext.Add("pbv", true)

	var ps string
	dataContext.Add("ps", &ps)
	dataContext.Add("psv", "zezeze")

	var S1 Single
	dataContext.Add("st1", &S1)

	ss := Single{Hello: "world"}
	dataContext.Add("ss", ss)

	var S2 Single
	dataContext.Add("st2", &S2)
	in := In{S: Single{Hello: "happy ending"}}
	dataContext.Add("in", &in)

	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	//读取规则
	e1 := ruleBuilder.BuildRuleFromString(set_single_value_rule)
	if e1 != nil {
		panic(e1)
	}

	eng := engine.NewGengine()
	// true: means when there are many rules， if one rule execute error，continue to execute rules after the occur error rule
	e2 := eng.Execute(ruleBuilder, true)
	if e2 != nil {
		panic(e2)
	}

	//println("outer--->", px)

	println("outer--->", pi)
	println("outer--->", pf)
	println("outer--->", pu)
	println("outer--->", pb)
	println("outer--->", ps)
	println("outer--->", S1.Hello)
	println("outer--->", S2.Hello)
}
