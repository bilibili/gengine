package test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
)

func GetNum(i int) int {
	fmt.Println("GetNum：" + strconv.Itoa(i))
	return i
}

const shortCircuitLogicRule = `
rule "short-circuit test"
begin
		if 7 == GetNum(8) || 6 == GetNum(6) || 7 == GetNum(7) {
			println("hit first")
		}else{
			println("hit second")
		}

        if 7 == GetNum(8) && 6 == GetNum(6)  {
			println("hit third")
		}else{
			println("hit fourth")
		}
end
`

func Test_expression_logic(t *testing.T) {
	//printMock := NewUnitPrintMock()
	dataContext := context.NewDataContext()
	dataContext.Add("GetNum", GetNum)
	dataContext.Add("println", fmt.Println)
	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)
	//读取规则
	start1 := time.Now().UnixNano()
	err := ruleBuilder.BuildRuleFromString(shortCircuitLogicRule)
	if err != nil {
		fmt.Printf("Build failed\n")
		panic(err)
	}
	end1 := time.Now().UnixNano()
	println(fmt.Sprintf("rules num:%d, load rules cost time:%d ns", len(ruleBuilder.Kc.RuleEntities), end1-start1))

	eng := engine.NewGengine()
	// true: means when there are many rules， if one rule execute error，continue to execute rules after the occur error rule
	err = eng.Execute(ruleBuilder, true)
	if err != nil {
		panic(err)
	}
}
