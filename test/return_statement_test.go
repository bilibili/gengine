package test

import (
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
)

//程序基本骨架
//program framework
func framework(rule string) map[string]interface{} {
	dataContext := context.NewDataContext()
	dataContext.Add("getInt", getInt)
	ruleBuilder := builder.NewRuleBuilder(dataContext)
	e1 := ruleBuilder.BuildRuleFromString(rule)
	if e1 != nil {
		panic(e1)
	}

	eng := engine.NewGengine()
	e2 := eng.Execute(ruleBuilder, true)
	if e2 != nil {
		println(e2)
	}

	resultMap, _ := eng.GetRulesResultMap()
	return resultMap
}

func getInt() int {
	return 666
}

//test
func Test_return_nil_1(t *testing.T) {

	//无返回值,返回nil
	ruleName := "return_in_statements"
	rule := `rule  "` + ruleName + `"
			 begin
			  return
      		 end	
			`
	returnResultMap := framework(rule)

	//无返回值，或者说是返回nil
	r := returnResultMap[ruleName]
	if r == nil {
		println("return is nil")
	}
}

//test
func Test_return_int64(t *testing.T) {

	ruleName := `return_in_statements`
	rule := `rule  "` + ruleName + `"
			 begin
			  return 56 //此处写的是数字
      		 end	
			`
	returnResultMap := framework(rule)

	//返回int值,且为int64
	r := returnResultMap[ruleName]
	i := r.(int64) //如果写的是数字, 则类型为int64
	println("return--->", i)
}

//test
func Test_return_function_value_int(t *testing.T) {

	ruleName := `return_in_statements`
	rule := `rule  "` + ruleName + `"
			 begin
			  return getInt() //此处是函数
      		 end	
			`
	returnResultMap := framework(rule)

	//返回int值
	i := returnResultMap[ruleName]
	ix := i.(int) //强类型转换, 如果是函数返回值, 函数返回的类型是什么，就是什么
	//如果小类型(大类型:整形, 小类型分为int,int8,int16,int32,int64)不太确定是什么,可以这样:
	//ix := reflect.ValueOf(i).Int() //这也是golang让人不爽的地方

	println("return--->", ix)
}

//test
func Test_return_float64(t *testing.T) {

	ruleName := `return_in_statements`
	rule := `rule  "` + ruleName + `"
			 begin
			  return 56.6
      		 end	
			`
	returnResultMap := framework(rule)

	//返回float64
	r := returnResultMap[ruleName]
	f := r.(float64)
	println("return--->", f)
}

//test
func Test_return_bool(t *testing.T) {

	ruleName := `return_in_statements`
	rule := `rule  "` + ruleName + `"
			 begin
			  return true
      		 end	
			`
	returnResultMap := framework(rule)

	//返回bool值
	r := returnResultMap[ruleName]
	b := r.(bool)
	println("return--->", b)
}

//test
func Test_return_string(t *testing.T) {

	ruleName := `return_in_statements`
	rule := `rule  "` + ruleName + `"
			 begin
			  return "hello world"
      		 end	
			`
	returnResultMap := framework(rule)

	//返回string值
	r := returnResultMap[ruleName]
	s := r.(string)
	println("return--->", s)
}

//test
func Test_return_logic_expression_bool(t *testing.T) {

	ruleName := `return_in_statements`
	rule := `rule  "` + ruleName + `"
			 begin
			  return 5 > 4 + 3
      		 end	
			`
	returnResultMap := framework(rule)

	//返回bool值
	r := returnResultMap[ruleName]
	b := r.(bool)
	println("return--->", b)
}

//test
func Test_return_logic_expression_int64(t *testing.T) {

	ruleName := `return_in_statements`
	rule := `rule  "` + ruleName + `"
			 begin
			  return 5 + 6 * 6
      		 end	
			`
	returnResultMap := framework(rule)

	//返回int64值
	r := returnResultMap[ruleName]
	i := r.(int64)
	println("return--->", i)
}

//test
func Test_return_variable_int64(t *testing.T) {

	ruleName := `return_in_statements`
	rule := `rule  "` + ruleName + `"
			 begin
			    b = 5 + 6 * 6
				return b
      		 end	
			`
	returnResultMap := framework(rule)

	r := returnResultMap[ruleName]
	i := r.(int64)
	println("return--->", i)
}

//test
func Test_return_simple_if_return_int64(t *testing.T) {

	ruleName := `return_in_statements`
	rule := `rule  "` + ruleName + `"
			 begin
			  if false { //修改此处，会得到不同的返回结果
				  return 5 + 6 * 6			
    		  }
			  return 22
      		 end	
			`
	returnResultMap := framework(rule)

	r := returnResultMap[ruleName]
	i := r.(int64)
	println("return--->", i)
}

//test 这里的表现和golang return略有不同
func Test_return_complex_if_return_int64(t *testing.T) {

	ruleName := `return_in_statements`
	rule := `rule  "` + ruleName + `"
			 begin
			  a = 290	//修改这里,就可以得到不同的结果
			  if a < 100 {
				  return 5 + 6 * 6			
    		  }else if a >= 100 && a < 200  {
  				  return "hello world"	
              }else {
				  return true	
			  }
      		 end	
			`
	returnResultMap := framework(rule)
	r := returnResultMap[ruleName]

	// a < 100时 返回int64
	//i := r.(int64)

	//a >= 100 && a < 200时 返回string
	//i := r.(string)

	//a >= 200
	i := r.(bool)

	println("return--->", i)
}
