package test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
)

func Test_dag(t *testing.T) {

	row, column := 3, 4
	var answer [][]int


	for i := 0; i < row; i++ {
		inline := make([]int, column)
		answer = append(answer, inline)
	}
	fmt.Println("line len=",len(answer), "row len=", len(answer[1]))

	// 方法1
	answer1 := make([][]int, row)
	for i := range answer1 {
		answer1[i] = make([]int, column)
	}
	fmt.Println(answer1)

}

func hello(x int, y string){
	print(x, y)
}

const dag_rules = `
rule "1"
begin
print(1, ",")
end

rule "2"
begin
print(2, ",")
end

rule "3"
begin
print(3, ",")
end

rule "4"
begin
print(4, ",")
end

rule "5"
begin
print(5, ",")
end

rule "6"
begin
print(6, ",")
end

rule "7"
begin
print(7, ",")
end

rule "8"
begin
print(8, ",")
end

rule "9"
begin
print(9, ",")
end

rule "10"
begin
print(10, ",")
end

rule "11"
begin
print(11, ",")
end

rule "12"
begin
print(12, ",")
end
`

func makeDAG() [][]string {
	names := make([][]string, 5)

	namesCol1 := make([]string, 3)
	namesCol1[0] = "1"
	namesCol1[1] = "2"
	namesCol1[2] = "3"

	/*	namesCol1[3] = "4"
		namesCol1[4] = "5"
		namesCol1[5] = "6"
		namesCol1[6] = "7"
		namesCol1[7] = "8"
		namesCol1[8] = "9"
		namesCol1[9] = "10"
		namesCol1[10] = "11"
		namesCol1[11] = "12"
	*/
	names[0] = namesCol1

	namesCol2 := make([]string, 1)
	namesCol2[0] = "4"
	names[1] = namesCol2

	namesCol3 := make([]string, 5)
	namesCol3[0] = "5"
	namesCol3[1] = "6"
	namesCol3[2] = "7"

	//add the rules not exist
	namesCol3[3] = "100"
	namesCol3[4] = "200"
	names[2] = namesCol3


	namesCol4 := make([]string, 5)
	namesCol4[0] = "8"
	namesCol4[1] = "9"
	namesCol4[2] = "10"
	namesCol4[3] = "11"
	namesCol4[4] = "12"
	names[3] = namesCol4
	return names
}



func Test_dag_run(t *testing.T) {

	dataContext := context.NewDataContext()
	dataContext.Add("print", hello)

	ruleBuilder := builder.NewRuleBuilder(dataContext)
	e1 := ruleBuilder.BuildRuleFromString(dag_rules)
	if e1 != nil {
		panic(e1)
	}

	gengine := engine.NewGengine()
	names := makeDAG()

	e := gengine.ExecuteDAGModel(ruleBuilder, names)
	if e != nil {
		panic(e)
	}
}

func Test_pool_dag_run(t *testing.T) {

	apis := make(map[string]interface{})
	apis["print"] = hello

	pool,e1 := engine.NewGenginePool(1,2,1,dag_rules,apis)
	if e1 != nil {
		panic(e1)
	}

	names := makeDAG()

	e, _ := pool.ExecuteDAGModel(names, make(map[string]interface{}))
	if e != nil {
		panic(e)
	}
}



