package math

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
	"time"
)

// in golang
//-------- just not want to lost data, if the transformation will not lose data, it will transfer success ,else it will be failed  -------------------------
func Test_Type(t *testing.T) {

	//right
	var x1 int
	x1 = 7.0
	println(x1)

	//wrong
	var x2 int
	//x2 = 7.1
	println(x2)

	y := -12.4
	var x3 = int64(y) //transfer success, but lose data
	var x4 = int64(11)
	var x5 float64
	x5 = float64(x3) + float64(x4)
	println(x5)

	var x6 uint
	x6 = 7
	println(x6)

	var x7 uint
	x7 = 7
	println(x7)

	var x8 uint64
	x8 = 8
	println(x8)

	var x9 float32
	x9 = float32(x8)
	println(x9)

	//var x10 int32
	//x10 = 9.869
	//println(x10)

	var x11 uint32
	x11 = uint32(y)
	println(x11)

}

//
type IntNum struct {
	Int   int
	Int8  int8
	Int16 int16
	Int32 int32
	Int64 int64
}

type FloatNum struct {
	Float32 float32
	Float64 float64
}

type UintNum struct {
	Uint   uint
	Uint8  uint8
	Uint16 uint16
	Uint32 uint32
	Uint64 uint64
}

const math_rule = `
rule "math" "math and number type test"
begin
//IntNum.Int = 100.0  //error: 100.0 resolved as  float，but IntNum.Int is integer type
//println(IntNum.Int)
FloatNum.Float32 = 100 //right
println(FloatNum.Float32)
end
`

func Test_math(t *testing.T) {
	math_exec()
}

func math_exec() {

	IntNum := &IntNum{}
	FloatNum := &FloatNum{}
	UintNum := &UintNum{}

	dataContext := context.NewDataContext()
	dataContext.Add("IntNum", IntNum)
	dataContext.Add("FloatNum", FloatNum)
	dataContext.Add("UintNum", UintNum)
	dataContext.Add("println", fmt.Println)

	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	//resolve rules from string
	start1 := time.Now().UnixNano()
	err := ruleBuilder.BuildRuleFromString(math_rule)
	end1 := time.Now().UnixNano()

	println(fmt.Sprintf("rules num:%d, load rules cost time:%d ns", len(ruleBuilder.Kc.RuleEntities), end1-start1))

	if err != nil {
		panic(err)
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
