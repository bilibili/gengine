package test

import (
	"encoding/json"
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"strings"
	"testing"
)

func JsonString(s string) {

	ms := make(map[string]map[string]string)

	//this is very important!!!
	bytes := []byte(strings.Replace(s, "\\", "", -1))

	e := json.Unmarshal(bytes, &ms)
	if e != nil {
		panic(e)
	}
	println(ms["hello"]["never"])
}

func Test_hello_222(t *testing.T) {

	/*	ms := make(map[string]map[string]string)
		m := make(map[string]string)
		m["never"] = "give up"
		ms["hello"] = m
		bytes, _ := json.Marshal(&ms)
		fmt.Println(string(bytes))
	*/
	ruleStr := `
rule "test" "test_d" salience 100
begin
	//a = "{"hello":{"never":"give up"}}" 不可以
	a = "{\"hello\":{\"never\":\"give up\"}}"	//可以
	println(a)
	jsonString(a)
end
`

	dataContext := context.NewDataContext()
	dataContext.Add("println", fmt.Println)
	dataContext.Add("jsonString", JsonString)
	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	//resolve rules from string
	err := ruleBuilder.BuildRuleFromString(ruleStr)

	if err != nil {
		panic(err)
	}
	eng := engine.NewGengine()
	// true: means when there are many rules， if one rule execute error，continue to execute rules after the occur error rule
	err = eng.Execute(ruleBuilder, true)
	if err != nil {
		panic(err)
	}
}
