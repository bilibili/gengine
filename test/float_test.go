package test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
)

const floatRule = `
rule "float" "float test"
begin
x1 =  0.9
println(x1)

x2 = .7
println(x2)

x3 = 7E-10
println(x3)

x4 = .7E-10
println(x4)

x5 = .1e7
println(x5)

end
`

func Test_float(t *testing.T) {

	x1 := .1e7
	println("x1--->", x1)

	x2 := 6e7
	println("x2--->", x2)
	dataContext := context.NewDataContext()
	dataContext.Add("println", fmt.Println)

	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	//resolve rules from string
	err := ruleBuilder.BuildRuleFromString(floatRule)
	if err != nil {
		panic(err)
	}

	eng := engine.NewGengine()
	err = eng.Execute(ruleBuilder, false)
	if err != nil {
		panic(err)
	}

}
