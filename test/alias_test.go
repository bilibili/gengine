package test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
)

type RoomScore struct {
	RoomId int64
	Score  float64
	Pctr   float64
	Small  int64
}

type Rs RoomScore

const alias_rule = `
rule "1" "2" 
begin
Rs.Small = 1000
println(Rs.Small)
end
`

func Test_alias(t *testing.T) {

	Rs := &Rs{}

	dataContext := context.NewDataContext()
	dataContext.Add("Rs", Rs)
	dataContext.Add("println", fmt.Println)

	ruleBuilder := builder.NewRuleBuilder(dataContext)
	e1 := ruleBuilder.BuildRuleFromString(alias_rule)
	if e1 != nil {
		panic(e1)
	}

	gengine := engine.NewGengine()
	e := gengine.Execute(ruleBuilder, true)
	if e != nil {
		panic(e)
	}
}
