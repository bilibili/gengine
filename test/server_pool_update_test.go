package test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
	"time"
)

const rule_to_update1 = `
rule "3" "sid=3" salience 100
begin
sid =3
rc = room.Recall("a")
ft = room.Feature(rc)
fl = room.Filter(ft)
md = room.Model(fl)
res = room.Rerank(md)
result.Add(rc, ft, fl, md, res)
end
`

const rule_to_update2 = `
rule "3" "sid=3" salience 100
begin
sid =3
rc = room.Recall("a")
ft = room.Feature(rc)
fl = room.Filter(ft)
md = room.Model(fl)
res = room.Rerank(md)
result.Add(rc, ft, fl, md, res)
//println("hellow")
end
rule "23" "xxx"
begin
rc = room.Recall("a")
end
`

type Room struct {
}

func (r *Room) Recall(s string) int {
	time.Sleep(3 * time.Millisecond)
	return 1
}

func (r *Room) Feature(i int) int {
	time.Sleep(7 * time.Millisecond)
	return 2
}

func (r *Room) Filter(i int) int {
	time.Sleep(4 * time.Millisecond)
	return 3
}

func (r *Room) Model(i int) int {
	time.Sleep(20 * time.Millisecond)
	return 4
}

func (r *Room) Rerank(i int) int {
	time.Sleep(20 * time.Millisecond)
	return 5
}

type ResultSet struct {
	//sync.Mutex
	M map[int]int
}

func (r *ResultSet) Add(ii ...int) {
	//r.Lock()
	//println("--------",fmt.Sprintf("len=%d, ii=%+v",len(ii), ii))

	for _, v := range ii {
		r.M[v] = v
	}

	//r.Unlock()
}

type Sengine struct {
	Rb *builder.RuleBuilder
	Eg *engine.Gengine
}

func Test_single_engine(t *testing.T) {

	egs := make(chan *Sengine, 2)

	for i := 0; i < 2; i++ {
		dc := context.NewDataContext()
		rb := builder.NewRuleBuilder(dc)
		e := rb.BuildRuleFromString(rule_to_update1)
		if e != nil {
			panic(e)
		}
		eg := engine.NewGengine()
		egs <- &Sengine{
			Rb: rb,
			Eg: eg,
		}
	}

	go exe2(egs)
	go exe2(egs)
	go exe2(egs)
	go exe2(egs)
	go exe2(egs)

	time.Sleep(20 * time.Second)

}

func Test_server_pool(t *testing.T) {

	apis := make(map[string]interface{})
	apis["println"] = fmt.Println
	start := time.Now()
	pool, e1 := engine.NewGenginePool(1, 2, 1, rule_to_update1, apis)
	fmt.Println("cost time", time.Since(start).Milliseconds())
	if e1 != nil {
		panic(e1)
	}

	go exe1(pool)
	go exe1(pool)
	go exe1(pool)
	go exe1(pool)
	go exe1(pool)

	go update(pool)
	go update(pool)
	time.Sleep(20 * time.Second)

}

func exe1(pool *engine.GenginePool) {
	for {
		data := make(map[string]interface{})
		data["room"] = &Room{}
		result := &ResultSet{M: make(map[int]int)}
		data["result"] = result

		sids := []string{"3"}
		e, _ := pool.ExecuteSelectedRulesConcurrent(data, sids)
		if e != nil {
			panic(e)
		}

		var ne []int
		for i := 1; i <= 5; i++ {
			if _, ok := result.M[i]; !ok {
				ne = append(ne, i)
			}
		}

		if len(ne) != 0 {
			panic(fmt.Sprintf("other not exist :%+v", ne))
		} /*else {
			println(fmt.Sprintf("exist :%+v", result.M))
		}*/
	}
}

func update(pool *engine.GenginePool) {
	for i := 0; i < 1000; i++ {
		time.Sleep(3 * time.Second)
		if i%2 == 0 {
			e := pool.UpdatePooledRules(rule_to_update2)
			println("-------1-------index=", i, "number=", pool.GetRulesNumber())
			if e != nil {
				panic("update rule err:" + fmt.Sprintf("%+v", e))
			}
		} else {
			e := pool.UpdatePooledRules(rule_to_update1)
			println("-------2-------index=", i, "number=", pool.GetRulesNumber())
			if e != nil {
				panic("update rule err:" + fmt.Sprintf("%+v", e))
			}
		}
	}
}

func exe2(egs chan *Sengine) {
	for {
		sg := <-egs
		sg.Rb.Dc.Add("room", &Room{})
		sg.Rb.Dc.Add("room", &Room{})
		result := &ResultSet{M: make(map[int]int)}
		sg.Rb.Dc.Add("result", result)
		sids := []string{"3"}
		e := sg.Eg.ExecuteSelectedRulesConcurrent(sg.Rb, sids)
		if e != nil {
			println(fmt.Sprintf("err:%+v", e))
		}

		var ne []int
		for i := 1; i <= 5; i++ {
			if _, ok := result.M[i]; !ok {
				ne = append(ne, i)
			}
		}

		if len(ne) != 0 {
			panic(fmt.Sprintf("other not exist :%+v", ne))
		}

		egs <- sg
	}

}
