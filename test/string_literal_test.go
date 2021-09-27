package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
)

const ruleStringLiteral = `
rule "string_literal_test" "test"
begin


debug("abc")
debug("ab\"c")
debug("ab\\c")
debug("ab\nc")
debug("ab\tc")
debug("ab\u2333c")
debug("ab\\u2333c")
debug("Hello, 世界")
debug("Hello, \u4e16\u754c")
debug("\xe4\xb8\x96\xe7\x95\x8c")
debug("")
//debug("ab
//c")
name := "abc\"def"
debug(name)

end
`

func TestStringLiteral(t *testing.T) {

	dataContext := context.NewDataContext()
	dataContext.Add("debug", t.Log)

	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	//resolve rules from string
	start1 := time.Now().UnixNano()
	err := ruleBuilder.BuildRuleFromString(ruleStringLiteral)
	end1 := time.Now().UnixNano()

	println(fmt.Sprintf("rules num:%d, load rules cost time:%d ns", len(ruleBuilder.Kc.RuleEntities), end1-start1))

	if err != nil {
		panic(err)
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
