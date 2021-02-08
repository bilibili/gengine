package pool

import (
	"fmt"
	"github.com/bilibili/gengine/engine"
	"testing"
	"time"
)

const pool_mix_model_rule = `
rule "best" "best"  salience 100
begin
	println("Ps.P....", Ps.P)
	Ps.P = true
end

rule "better" "better"   salience 99
begin

if Ps.P {
	println("better....")

	Ps.R = true
	println("better....",Ps.R)
}
end


rule "good" "good"   salience 98
begin

if Ps.P {
	println("good....")
	Ps.R = false
	println("good....",Ps.R)
}
end
`

type Ps struct {
	P bool
	R bool
}

func Test_mix_model(t *testing.T) {

	apis := make(map[string]interface{})

	pool, e1 := engine.NewGenginePool(1, 3, 3, pool_mix_model_rule, apis)
	if e1 != nil {
		panic(e1)
	}

	println("pool.GetExecModel()==", pool.GetExecModel())
	data := make(map[string]interface{})
	Ps := &Ps{}
	data["Ps"] = Ps
	data["println"] = fmt.Println
	/*	e,_ := pool.ExecuteSelectedRulesMixModelWithMultiInput(data, []string{"best", "better", "good"})
		if e != nil {
			panic(e)
		}
		println("Ps.R=", Ps.R)*/

	e, _ := pool.ExecuteSelectedWithSpecifiedEM(data, []string{"best", "better", "good"})
	if e != nil {
		panic(e)
	}
	println("Ps.R=", Ps.R)

	time.Sleep(3 * time.Second)

}
