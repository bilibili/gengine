package test

import (
	"fmt"
	"github.com/bilibili/gengine/engine"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

type Reqest struct {
	Data float32
}

const (
	pool_rule = `
rule "测试规则名称1" "rule desc"
begin
	Req.Data = 10 + 7 + 8
	//sleep()   
	 print(Req.Data)
end
rule "1"
begin
	Req.Data = 10 + 7 + 8
	//sleep()
	//print("1")
end
rule "2" "rule desc"
begin
	Req.Data = 10 + 7 + 8
	//sleep()
	//	   print("2")
end`
)

func Test_pool_select_rules(t *testing.T) {

	apis := make(map[string]interface{})
	apis["print"] = fmt.Println
	pool, e1 := engine.NewGenginePool(1, 2, 1, pool_rule, apis)
	if e1 != nil {
		panic(e1)
	}

	reqest := &Reqest{}
	data := make(map[string]interface{})
	data["Req"] = reqest

	e2, _ := pool.ExecuteSelectedRules(data, []string{"测试规则名称1", "1"})
	if e2 != nil {
		panic(e2)
	}

	println("pool.GetRulesNumber()---->", pool.GetRulesNumber())
	sal, e1 := pool.GetRuleSalience("2")
	if e1 != nil {
		panic(e1)
	}
	println("sal --->", sal)

	desc, e2 := pool.GetRuleDesc("1")
	if e2 != nil {
		panic(e2)
	}
	println("desc--->", desc)

	exist := pool.IsExist([]string{"333"})
	println(fmt.Sprintf("rule 333 exist--->%+v", exist))

}

//test no rules in pool
/*func Test_pool_no_rules(t *testing.T) {

	t1 := time.Now()
	pool, e1 := engine.NewGenginePool(3, 5, 1, "", nil)
	if e1 != nil {
		panic(e1)
	}
	e2, _ := pool.ExecuteRules("", nil, "", nil)
	if e2 != nil {
		panic(e2)
	}

	println("cost time:", time.Since(t1), "ns")
}*/

func Test_once(t *testing.T) {

	apis := make(map[string]interface{})
	apis["print"] = fmt.Println
	pool, e1 := engine.NewGenginePool(1, 2, 1, pool_rule, apis)
	if e1 != nil {
		panic(e1)
	}

	reqest := &Reqest{}

	t1 := time.Now()
	e2, _ := pool.ExecuteRulesWithSpecifiedEM("Req", reqest, "", nil)
	if e2 != nil {
		panic(e2)
	}
	println("build pool cost time:", time.Since(t1), "ns")
}

func Sleep() {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(1000)
	time.Sleep(time.Nanosecond * time.Duration(i))
}

func Test_pool_with_rules_for_goruntine(t *testing.T) {

	poolMinLen := int64(1)
	poolMaxLen := int64(2)
	max := int64(0)
	min := int64(1000000)
	cnt := int64(0)

	g1 := int64(0)
	g2 := int64(0)
	g3 := int64(0)
	g4 := int64(0)
	g5 := int64(0)

	t1 := time.Now()
	apis := make(map[string]interface{})
	apis["print"] = fmt.Println
	apis["sleep"] = Sleep
	pool, e1 := engine.NewGenginePool(poolMinLen, poolMaxLen, 2, pool_rule, apis)
	if e1 != nil {
		panic(e1)
	}
	println("build pool cost time:", time.Since(t1), "ns")

	go func() {
		for {
			t2 := time.Now()
			reqest := &Reqest{Data: 1}
			e2, _ := pool.ExecuteRulesWithSpecifiedEM("Req", reqest, "", nil)
			if e2 != nil {
				println(fmt.Sprintf("e2: %+v", e2))
			}
			duration := time.Since(t2)
			if int64(duration) > max {
				atomic.StoreInt64(&max, int64(duration))
			}
			if int64(duration) < min {
				atomic.StoreInt64(&min, int64(duration))
			}
			atomic.AddInt64(&cnt, 1)
			g1++
			//println("1 exec cost time:", , "ns\n")
		}
	}()

	go func() {
		for {
			t2 := time.Now()
			reqest := &Reqest{Data: 1}
			e2, _ := pool.ExecuteRulesWithSpecifiedEM("Req", reqest, "", nil)
			if e2 != nil {
				println(fmt.Sprintf("e2: %+v", e2))
			}
			duration := time.Since(t2)
			if int64(duration) > max {
				atomic.StoreInt64(&max, int64(duration))
			}
			if int64(duration) < min {
				atomic.StoreInt64(&min, int64(duration))
			}
			atomic.AddInt64(&cnt, 1)
			g2++
			//println("2 exec cost time:", time.Since(t2), "ns\n")
		}
	}()

	go func() {
		for {
			t2 := time.Now()
			reqest := &Reqest{Data: 1}
			e2, _ := pool.ExecuteRulesWithSpecifiedEM("Req", reqest, "", nil)
			if e2 != nil {
				println(fmt.Sprintf("e2: %+v", e2))
			}
			duration := time.Since(t2)
			if int64(duration) > max {
				atomic.StoreInt64(&max, int64(duration))
			}
			if int64(duration) < min {
				atomic.StoreInt64(&min, int64(duration))
			}
			atomic.AddInt64(&cnt, 1)
			g3++
			//println("3 exec cost time:", time.Since(t2), "ns\n")
		}
	}()

	go func() {
		for {
			t2 := time.Now()
			reqest := &Reqest{Data: 1}
			e2, _ := pool.ExecuteRulesWithSpecifiedEM("Req", reqest, "", nil)
			if e2 != nil {
				println(fmt.Sprintf("e2: %+v", e2))
			}
			duration := time.Since(t2)
			if int64(duration) > max {
				atomic.StoreInt64(&max, int64(duration))
			}
			if int64(duration) < min {
				atomic.StoreInt64(&min, int64(duration))
			}
			g4++
			atomic.AddInt64(&cnt, 1)
			//println("4 exec cost time:", time.Since(t2), "ns\n")
		}
	}()

	go func() {
		for {
			t2 := time.Now()
			reqest := &Reqest{Data: 1}
			e2, _ := pool.ExecuteRulesWithSpecifiedEM("Req", reqest, "", nil)
			if e2 != nil {
				println(fmt.Sprintf("e2: %+v", e2))
			}
			duration := time.Since(t2)
			if int64(duration) > max {
				atomic.StoreInt64(&max, int64(duration))
			}
			if int64(duration) < min {
				atomic.StoreInt64(&min, int64(duration))
			}
			atomic.AddInt64(&cnt, 1)
			g5++
			//println("5 exec cost time:", time.Since(t2), "ns\n")
		}
	}()

	go func() {
		i := 0
		for {
			time.Sleep(1 * time.Second)
			i++
			println("poolMinLen=", poolMinLen, ", poolMaxLen=", poolMaxLen, ", sort", i, ", min: ", min, "ns, max: ", max, "ns, request-QPS:", int(cnt)/i, ", g1:", g1, ",g2:", g2, ",g3:", g3, ",g4:", g4, ",g5:", g5)
		}

	}()

	println("test 10 seconds...")
	time.Sleep(5 * time.Second)
}
