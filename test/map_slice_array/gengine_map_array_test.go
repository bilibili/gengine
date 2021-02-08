package map_slice_array

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
	"time"
)

type MapArray struct {
	Mx map[string]bool
	Ax [3]int
	Sx []string
}

const ma_rule = `
rule "测试规则" "rule desc"
begin
  
  x = Ma.Mx["hello"]
  
  PrintName(x)
  
  Ma.Mx["hello"] = false
  b = "your"
  
  Ma.Mx[b]= true
  
  y = Ma.Mx	
  PrintName("------",y["hello"])

  if x {
     PrintName("Single data")
  }	

  if 2 == 2 {
     PrintName("true == true")
  }
   
  if x == true {
   PrintName("haha")
  } 

  if !x {
     PrintName("haha")
  }else{
     PrintName("!x")
  }
  
  xx = Ma.Ax[2]
  PrintName(xx) 
  Ma.Ax[2] = 300011111
  PrintName(Ma.Ax[2]) 
 
  yy = Ma.Ax
  PrintName(yy[1]) 
  
  
  if yy[2] == 20000 {
     PrintName("20000")
  }

  z = Ma.Sx[1]
  PrintName("z--1--->",z) 

//you can read data from zz,but you can set value to zz
  zz = Ma.Sx
  if zz[2] == "kkkk"{
     PrintName("z--2--->","kkkk") 
  }


  a = 2
  Ma.Sx[a] = "MMMM"
  PrintName("z--3-->", Ma.Sx[a]) 

end
`

func Test_map_array(t *testing.T) {

	Ma := &MapArray{
		Mx: map[string]bool{"hello": true},
		Ax: [3]int{1000, 20000, 300},
		Sx: []string{"jjj", "lll", "kkkk"},
	}

	dataContext := context.NewDataContext()
	dataContext.Add("PrintName", fmt.Println)
	dataContext.Add("Ma", Ma)

	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	//读取规则
	start1 := time.Now().UnixNano()
	err := ruleBuilder.BuildRuleFromString(ma_rule)
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

func Test_unptr_map(t *testing.T) {

	Ma := make(map[int]string)

	dataContext := context.NewDataContext()
	dataContext.Add("PrintName", fmt.Println)
	dataContext.Add("Ma", Ma)

	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	err := ruleBuilder.BuildRuleFromString(`
rule "1"
begin
a = 1
Ma[a] = "xxx"
end
`)

	if err != nil {
		panic(err)
	}
	eng := engine.NewGengine()
	// true: means when there are many rules， if one rule execute error，continue to execute rules after the occur error rule
	err = eng.Execute(ruleBuilder, true)
	if err != nil {
		panic(err)
	}
	println(Ma[1])
}
