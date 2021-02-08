package complex

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
	"time"
)

//求并集
func union(slice1, slice2 []string) []string {
	m := make(map[string]int)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 0 {
			slice1 = append(slice1, v)
		}
	}
	return slice1
}

//求交集
func intersect(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

//求差集
func difference(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	inter := intersect(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}

	for _, value := range slice1 {
		times, _ := m[value]
		if times == 0 {
			nn = append(nn, value)
		}
	}
	return nn
}

func Test_sets(t *testing.T) {

	s1 := []string{"saa", "sss", "xxxx"}
	s2 := []string{"sa1", "sss2", "xxxx"}

	/*	s3 := union(s1, s2)
		for _,v := range s3 {
			println(v)
		}*/

	/*	s4 := intersect(s1, s2)
		for _,v := range s4 {
			println(v)
		}*/

	s5 := difference(s1, s2)
	for _, v := range s5 {
		println(v)
	}
}

type Data struct {
	S1 []string
	S2 []string
	S3 []string
}

const rule = `
rule "测试交叉并" "rule desc"
begin
data.S3 = difference(data.S1, data.S2)
end
`

func exec() {

	data := &Data{
		S1: []string{"111", "2222", "333"},
		S2: []string{"1111", "222", "333"},
		S3: []string{},
	}

	dataContext := context.NewDataContext()
	dataContext.Add("data", data)
	dataContext.Add("union", union)
	dataContext.Add("intersect", intersect)
	dataContext.Add("difference", difference)

	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	//读取规则
	start1 := time.Now().UnixNano()
	err := ruleBuilder.BuildRuleFromString(rule)
	end1 := time.Now().UnixNano()

	println(fmt.Sprintf("rules num:%d, load rules cost time:%d ns", len(ruleBuilder.Kc.RuleEntities), end1-start1))

	if err != nil {
		panic(err)
	} else {
		eng := engine.NewGengine()

		start := time.Now().UnixNano()
		// true: means when there are many rules， if one rule execute error，continue to execute rules after the occur error rule
		err := eng.Execute(ruleBuilder, true)
		if err != nil {
			panic(err)
		}
		for _, v := range data.S3 {
			println(v)
		}
		end := time.Now().UnixNano()
		println(fmt.Sprintf("execute rule cost %d ns", end-start))
	}
}

func Test_engine(t *testing.T) {
	exec()

}
