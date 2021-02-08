package line_number

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
)

var lineNumberRules = `
rule "line_number"  "when execute error,gengine will give out error"
begin

//Println("golang", "hello", "world" )

//取消斜杠注释，依次测试不同的报错情况
//if Println("golang", "hello") == 100 {
// 	Println("golang", "hello")
//}

ms.X()

end
`

type MyStruct struct {
}

func (m *MyStruct) XX(s string) {
	println("XX")
}

func Println(s1, s2 string) bool {
	println(s1, s2)
	return false
}

func Test_number(t *testing.T) {
	dataContext := context.NewDataContext()
	//注入自定义函数
	dataContext.Add("Println", Println)
	ms := &MyStruct{}
	dataContext.Add("ms", ms)

	ruleBuilder := builder.NewRuleBuilder(dataContext)
	e1 := ruleBuilder.BuildRuleFromString(lineNumberRules)
	if e1 != nil {
		panic(e1)
	}

	eng := engine.NewGengine()
	// true: means when there are many rules， if one rule execute error，continue to execute rules after the occur error rule
	e2 := eng.Execute(ruleBuilder, true)
	if e2 != nil {
		println(fmt.Sprintf("%+v", e2))
	}

}
