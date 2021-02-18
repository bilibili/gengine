package test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
)

type B struct {
	Name string
	Mp   map[string]string
	Ar   [5]string
	Sl   []int
}

func (b *B) Meth(s string) {
	println("--->", s)
}

func Test_three_level_call(t *testing.T) {

	type A struct {
		N  string
		Ma map[string]string
		B  *B
	}

	rule := `
rule "three level call"
begin
	//conc{
	A.B.Name = "xiaoMing"
	A.B.Mp["hello"] = "world"
	A.B.Ar[1] = "Calo"
	A.B.Sl[2] = 3
	x = A.B.Sl[0]
	A.B.Meth(A.B.Ar[1])

	A.N = "kakaka"
	A.Ma["a"] = "b"
	//}
	println(A.B.Name, A.B.Mp["hello"], A.B.Ar[1], A.B.Sl[2], x, A.N, A.Ma["a"])

	if A.B.Sl[0] == 0 {
		println(true)
	}
end
`

	b := B{
		Name: "",
		Mp:   make(map[string]string),
		Ar:   [5]string{},
		Sl:   make([]int, 6),
	}
	pA := &A{
		N:  "",
		Ma: make(map[string]string),
		B:  &b,
	}

	dataContext := context.NewDataContext()
	dataContext.Add("println", fmt.Println)
	dataContext.Add("A", pA)

	ruleBuilder := builder.NewRuleBuilder(dataContext)
	e := ruleBuilder.BuildRuleFromString(rule)
	if e != nil {
		panic(e)
	}

	gengine := engine.NewGengine()
	e = gengine.Execute(ruleBuilder, true)
	if e != nil {
		panic(e)
	}

	println(pA.B.Name, pA.B.Mp["hello"], pA.B.Ar[1], pA.B.Sl[2])
}
