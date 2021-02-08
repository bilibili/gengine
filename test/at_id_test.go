package test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
)

/**
use '@id',you can get rule name in rule content
*/
const atIDRule = `
rule "测试规则名称1" "rule desc" salience 10
begin
	println(@id)
end
rule " 100 " "rule desc" salience 20
begin
	x = @id
	println(x)
end
`

func Test_at_id(t *testing.T) {
	dataContext := context.NewDataContext()
	dataContext.Add("println", fmt.Println)

	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	//resolve rules from string
	err := ruleBuilder.BuildRuleFromString(atIDRule)
	if err != nil {
		panic(err)
	}

	eng := engine.NewGengine()
	err = eng.Execute(ruleBuilder, false)
	if err != nil {
		panic(err)
	}
}

type Data struct {
	M map[string]string
}

func (d *Data) exe() {
	println("hhhh")
}

func Test_need(t *testing.T) {

	//var data1 *Data
	//data1 = nil
	//println()
	//data1.exe()

	data := &Data{
		M: make(map[string]string),
	}
	data.M["a"] = "b"

	data1 := &Data{
		M: make(map[string]string),
	}
	data1.M["b"] = "c"

	dx := make(map[string]*Data)
	dx["hello"] = data
	d1 := dx["hello"]
	println("1--->", len(dx), d1.M["a"])
	dx["hello"] = data1
	d2 := dx["hello"]

	//delete(dx, "hello")
	println("1--->", len(dx), d2.M["a"])
	println(data.M["a"])

}
