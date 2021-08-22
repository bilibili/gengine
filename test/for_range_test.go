package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"github.com/stretchr/testify/assert"
)

const forRule = `
rule "array test" "m dec"
begin

// for statment
for i=1;i<=5;i=i+1{
	println("for i = ", i)
}

// forRange map
debugKey := "keylist"
forRange key := mapInfo{
	if key == "key1" {
		debugKey += "|" + key + ":" + mapInfo[key]
	}
}
println("debugKey is ",debugKey)

// forRange slice
forRange key := testSlice {
	println("slice = ", key, testSlice[key])
}
end
`

func Test_For(t *testing.T) {
	testSlice := []int64{6, 7, 8}
	mapInfo := map[string]string{"key1": "val1", "key2": "val2"}
	printMock := NewUnitPrintMock()
	dataContext := context.NewDataContext()
	dataContext.Add("PrintName", printMock.Debug)
	dataContext.Add("testSlice", testSlice)
	dataContext.Add("println", printMock.Debug)
	dataContext.Add("mapInfo", mapInfo)
	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)
	//读取规则
	start1 := time.Now().UnixNano()
	err := ruleBuilder.BuildRuleFromString(forRule)
	if err != nil {
		fmt.Printf("Build failed\n")
		panic(err)
	}
	end1 := time.Now().UnixNano()
	println(fmt.Sprintf("rules num:%d, load rules cost time:%d ns", len(ruleBuilder.Kc.RuleEntities), end1-start1))

	eng := engine.NewGengine()
	start := time.Now().UnixNano()
	// true: means when there are many rules， if one rule execute error，continue to execute rules after the occur error rule
	err = eng.Execute(ruleBuilder, true)
	end := time.Now().UnixNano()
	if err != nil {
		panic(err)
	}
	println(fmt.Sprintf("execute rule cost %d ns", end-start))

	// 测试下for循环
	assert.Equal(t, "for i = 3", printMock.Get(2))
	assert.Equal(t, "debugKey is keylist|key1:val1", printMock.Get(5))
	assert.Equal(t, "slice = 0 6", printMock.Get(6))
	assert.Equal(t, "slice = 1 7", printMock.Get(7))

}
