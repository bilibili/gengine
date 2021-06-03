package test

import (
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
)

func Test_nil(t *testing.T) {
	type Hello struct {
		X interface{}
	}
	He := &Hello{
		X: nil, //nil在此
	}

	mx := make(map[string]interface{})

	dc := context.NewDataContext()
	dc.Add("He", He)
	dc.Add("mx", &mx)

	rb := builder.NewRuleBuilder(dc)
	rule := `
rule "test" begin 
mx["666"] = He.X //在这里把nil赋值给interface{},
end`
	e1 := rb.BuildRuleFromString(rule)
	if e1 != nil {
		panic(e1)
	}
	gengine := engine.NewGengine()
	e2 := gengine.Execute(rb, false)
	if e2 != nil {
		panic(e2)
	}

	//如果在规则设置interface{}为nil是成功的，那么ok == true，同时会输出"yes"
	if i, ok := mx["666"]; ok {
		if i == nil {
			println("yes")
		}
	}
	//但也要记住一点，在golang中，是不允许有不为任何类型的nil的存在，言下之意就是，golang能处理的所有的nil，都是有类型的
}
