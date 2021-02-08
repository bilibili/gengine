package test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"

	"testing"
)

type Res struct {
	RiskLevel string
	X         map[string]string
}

func (r *Res) Handle(s string) string {
	println("Handle--->", s)
	return s
}

func NewFunc() {
	println("newFunc....")
}

func (r *Res) GetFunc() func() {
	return NewFunc
}

type CTX struct {
	Res *Res
}

func (result *Result) MyF(s string) {

}

func Test_hello(t *testing.T) {

	result := &Res{
		RiskLevel: "hello",
		X:         make(map[string]string),
	}
	ctx := &CTX{Res: result}

	dataContext := context.NewDataContext()
	dataContext.Add("ctx", ctx)
	dataContext.Add("println", fmt.Println)

	ruleBuilder := builder.NewRuleBuilder(dataContext)
	e1 := ruleBuilder.BuildRuleFromString(`
rule "1" 
begin
result = ctx.Res
result.RiskLevel = "E"
result.X["hello"]= "world"
result.Handle("xxxxxx") //二级数据附带的函数

myFunc = result.GetFunc() //二级数据返回一个函数
myFunc() //调用这个返回的函数

end
`)
	if e1 != nil {
		panic(e1)
	}

	gengine := engine.NewGengine()
	e := gengine.Execute(ruleBuilder, true)
	if e != nil {
		panic(e)
	}

	println("---x->", ctx.Res.RiskLevel, ctx.Res.X["hello"])

}

func Test_in_in(t *testing.T) {
	type Result struct {
		RiskLevel string
	}

	type CTX struct {
		Result Result
	}

	result := Result{
		RiskLevel: "hello",
	}
	ctx := &CTX{Result: result}

	dataContext := context.NewDataContext()
	dataContext.Add("ctx", ctx)
	dataContext.Add("println", fmt.Println)

	ruleBuilder := builder.NewRuleBuilder(dataContext)
	e1 := ruleBuilder.BuildRuleFromString(`
rule "1" 
begin
result = ctx.Result
result.RiskLevel = "E"
end
`)
	if e1 != nil {
		panic(e1)
	}

	gengine := engine.NewGengine()
	e := gengine.Execute(ruleBuilder, true)
	if e != nil {
		panic(e)
	}

	println("---x->", ctx.Result.RiskLevel)
}
