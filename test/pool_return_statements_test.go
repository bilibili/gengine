package test

import (
	"fmt"
	"github.com/bilibili/gengine/engine"
	"testing"
	"time"
)

func random() int {
	return time.Now().Nanosecond()
}

func Test_pool_return_statments(t *testing.T) {

	ruleName := "test_pool_return"
	rule := `rule "` + ruleName + `"  
			begin
				return random()
			end
			`

	apis := make(map[string]interface{})
	apis["print"] = fmt.Println
	apis["random"] = random
	pool, e1 := engine.NewGenginePool(1, 2, 1, rule, apis)
	if e1 != nil {
		println(fmt.Sprintf("e1: %+v", e1))
	}

	data := make(map[string]interface{})
	e2, rrm1 := pool.ExecuteRulesWithMultiInputWithSpecifiedEM(data)
	if e2 != nil {
		panic(e2)
	}

	i1 := rrm1[ruleName]
	ix1 := i1.(int)
	println("ix1--->", ix1)

	e3, rrm2 := pool.ExecuteRulesWithMultiInputWithSpecifiedEM(data)
	if e3 != nil {
		panic(e2)
	}

	i2 := rrm2[ruleName]
	ix2 := i2.(int)
	println("ix2--->", ix2)

	i11 := rrm1[ruleName]
	ix11 := i11.(int)
	println("ix11--->", ix11)

	e2, rrm3 := pool.ExecuteRulesWithMultiInputWithSpecifiedEM(data)
	if e2 != nil {
		panic(e2)
	}

	i4 := rrm3[ruleName]
	ix4 := i4.(int)
	println("ix4---->", ix4)

	i22 := rrm2[ruleName]
	ix22 := i22.(int)
	println("ix22--->", ix22)

	i111 := rrm1[ruleName]
	ix111 := i111.(int)
	println("ix111--->", ix111)

	if ix1 != ix11 && ix1 != ix111 && ix2 != ix22 {
		panic("返回值收集状态隔离异常...")
	}
}
