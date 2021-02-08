package test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
)

const operator_rules = `
rule "1"
begin
	a := "hello"
	a += "hello"
	println("a1--->", a)
	a += a + "world"
	println("a2--->", a)

	b := 10
	b += b // b = b + b
	println("b1--->", b)

	b = 10
	b += b * 10 // b = b + b * 10
	println("b2--->", b)

	b = 10
	b += b*10 + 9 // b = b + b * 10 + 9
	println("b3--->", b)

	b = 10
	b += b*10 + 9*2 // b = b + b * 10 + 9 * 2
	println("b4--->", b)

	b = 10
	b += b * (10 + 9) * 2 // b = b + b *（10+9）*2
	println("b5--->", b)

	b = 10
	b += (b*10 + 9) * 2 // b = b + (b *10+9）*2
	println("b6--->", b)

	b = 10
	b -= b // b = b - b
	println("b6--->", b)


	b = 10
	b *= b // b = b * b
	println("b7--->", b)

	b = 10
	b *= b * 10 // b = b * (b * 10)
	println("b8--->", b)

	b = 10
	b *= b*10 + 9 // b = b * (b * 10 + 9)
	println("b9--->", b)

	b = 10
	b *= b*10 + 9*2 // b = b * (b * 10 + 9 * 2)
	println("b10--->", b)

	b = 10
	b *= b * (10 + 9) * 2 // b = b * (b *（10+9）*2)
	println("b11--->", b)

	b = 10
	b *= (b*10 + 9) * 2 // b = b * ((b *10+9）*2)
	println("b12--->", b)

	b = 10
	b /=b+ 1+ 20  // b = b / (1 + 20)
	println("b13--->", b)
end
`

func Test_operator(t *testing.T) {

	dataContext := context.NewDataContext()
	//rename and inject
	dataContext.Add("println", fmt.Println)

	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	//读取规则
	e1 := ruleBuilder.BuildRuleFromString(operator_rules)
	if e1 != nil {
		panic(e1)
	}

	eng := engine.NewGengine()
	// true: means when there are many rules， if one rule execute error，continue to execute rules after the occur error rule
	e2 := eng.Execute(ruleBuilder, true)
	if e2 != nil {
		panic(e2)
	}

}

func Test_golang_base(t *testing.T) {

	a := "hello"
	a += a
	println("a1--->", a)
	a += a + "world"
	println("a2--->", a)

	b := 10
	b += b // b = b + b
	println("b1--->", b)

	b = 10
	b += b * 10 // b = b + b * 10
	println("b2--->", b)

	b = 10
	b += b*10 + 9 // b = b + b * 10 + 9
	println("b3--->", b)

	b = 10
	b += b*10 + 9*2 // b = b + b * 10 + 9 * 2
	println("b4--->", b)

	b = 10
	b += b * (10 + 9) * 2 // b = b + b *（10+9）*2
	println("b5--->", b)

	b = 10
	b += (b*10 + 9) * 2 // b = b + (b *10+9）*2
	println("b6--->", b)

	b = 10
	b -= b // b = b - b
	println("b6--->", b)

	b = 10
	b *= b // b = b * b
	println("b7--->", b)

	b = 10
	b *= b * 10 // b = b * (b * 10)
	println("b8--->", b)

	b = 10
	b *= b*10 + 9 // b = b * (b * 10 + 9)
	println("b9--->", b)

	b = 10
	b *= b*10 + 9*2 // b = b * (b * 10 + 9 * 2)
	println("b10--->", b)

	b = 10
	b *= b * (10 + 9) * 2 // b = b * (b *（10+9）*2)
	println("b11--->", b)

	b = 10
	b *= (b*10 + 9) * 2 // b = b * ((b *10+9）*2)
	println("b12--->", b)

	b = 10
	b /= b + 1 + 20 // b = b / (1 + 20)
	println("b13--->", b)

}
