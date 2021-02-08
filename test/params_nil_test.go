package test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
)

const params_nil_rule = `
rule "1" "2"
begin

st = SetSt()
x = GetSt(st)

end
`

type SS struct {
	Str string
	Sl  []int64
	Ab  AB
}

type AB struct {
	Sx string
}

func SetSt() *SS {
	return nil
}

func GetSt(s *SS) *SS {
	if s == nil {
		println("xxxxxx")
	} else {
		println("yyyyyy", len(s.Sl), fmt.Sprintf("%+v", s.Ab), s.Sl == nil, s.Ab.Sx == "")
	}

	return s
}

func Test_p_b(t *testing.T) {

	st := SetSt()
	_ = GetSt(st)

}

func Test_params_nil(t *testing.T) {

	dataContext := context.NewDataContext()
	//inject struct
	//rename and inject

	dataContext.Add("SetSt", SetSt)
	dataContext.Add("GetSt", GetSt)

	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	//读取规则
	e1 := ruleBuilder.BuildRuleFromString(params_nil_rule)
	if e1 != nil {
		panic(e1)
	}

	eng := engine.NewGengine()
	// true: means when there are many rules， if one rule execute error，continue to execute rules after the occur error rule
	e2 := eng.Execute(ruleBuilder, true)
	if e2 != nil {
		panic(e2)
	}

}
