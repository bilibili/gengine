package test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
	"time"
)

const rule4 = `
rule "111" "111" salience 80
BEGIN
		Print("111")
END
rule "222" "222" salience 10
BEGIN
		Print("222")
END
rule "333" "333" salience 11
BEGIN
	if false {
		Print("333")
	}else{
		Print("444")
	}
END

`

func Print(s string) {
	fmt.Println(s)
}

func Test_Priority(t *testing.T) {
	dataContext := context.NewDataContext()
	dataContext.Add("Print", Print)

	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	start1 := time.Now().UnixNano()
	err := ruleBuilder.BuildRuleFromString(rule4)
	end1 := time.Now().UnixNano()

	println(fmt.Sprintf("rules num:%d, load rules cost time:%d", len(ruleBuilder.Kc.RuleEntities), end1-start1))

	if err != nil {
		panic(err)
	}
	eng := engine.NewGengine()
	start := time.Now().UnixNano()
	err = eng.Execute(ruleBuilder, true)
	end := time.Now().UnixNano()
	if err != nil {
		panic(err)
	}
	println(fmt.Sprintf("execute rule cost %d ns", end-start))

}
