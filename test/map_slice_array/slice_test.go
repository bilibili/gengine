package map_slice_array

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"reflect"
	"testing"
	"time"
)

type SS struct {
	MI []int
	MM *[]int
}

const S_1 = `
rule "slice test" "slice dec"
begin

a = 1
//calculate
SS.MI[1] = 22 + 2 - 5 * 6 / 3
println("SS.MI[1]-----> ",SS.MI[1])
println("SS.MI[a]-----> ",SS.MI[a])

b = 3
SS.MM[3]=66666
println("SS.MM[3]---->", SS.MM[3])
println("SS.MM[b]---->", SS.MM[b])

a = 1
S[a] = 11111
println("S[1]---->",S[1])
println("S[a]---->",S[a])


SSS[1] = 3333
println("SSS[1]---->",SSS[1])
println("SSS[a]---->",SSS[a])
end
`

func Test_s1(t *testing.T) {
	SS := &SS{
		MI: []int{1, 2, 3, 4},
		MM: &[]int{9, 1, 34, 5},
	}

	var S []int
	S = []int{1, 2, 3}

	var SSS []int
	SSS = []int{1, 2, 3}

	dataContext := context.NewDataContext()
	dataContext.Add("PrintName", fmt.Println)
	dataContext.Add("SS", SS)
	dataContext.Add("S", &S)
	dataContext.Add("SSS", SSS)
	dataContext.Add("println", fmt.Println)

	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	//读取规则
	start1 := time.Now().UnixNano()
	err := ruleBuilder.BuildRuleFromString(S_1)
	end1 := time.Now().UnixNano()

	println(fmt.Sprintf("rules num:%d, load rules cost time:%d ns", len(ruleBuilder.Kc.RuleEntities), end1-start1))

	if err != nil {
		println(fmt.Sprintf("err:%s ", err))
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

func Test_slice_un(t *testing.T) {
	//可以设置
	x := []int{1, 2, 34}
	reflect.ValueOf(x).Index(1).Set(reflect.ValueOf(5))
	println("x--->", x[1])
}
