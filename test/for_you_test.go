package test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"reflect"
	"strings"
	"testing"
	"time"
)

type Container struct {
	//此处不带任何字段，就可以确保container指针的函数附带的是无状态函数（因为没有状态共享）
	//这样注入的时候也可以少写很多代码
}

// Log
func (c *Container) LogModel() {

}

// nil
func (c *Container) IsNil(a interface{}) bool {

	//暂时简写为这个
	return reflect.ValueOf(a).IsValid()
}

// String
func (c *Container) EqualsIgnoreCase(left, right string) bool {
	return strings.EqualFold(left, right)
}

func (c *Container) HasPrefix(str, prefix string) bool {
	return strings.HasPrefix(str, prefix)
}

func (c *Container) HasSuffix(str, suffix string) bool {
	return strings.HasSuffix(str, suffix)
}

// Contains
func (c *Container) ContainsForArray(arr []interface{}, item interface{}) bool {
	if len(arr) == 0 {
		return false
	}
	for _, value := range arr {
		if value == item {
			return true
		}
	}

	return false
}

func (c *Container) ContainsForString(str, item string) bool {
	return strings.Contains(str, item)
}

// Time
func (c *Container) CurrentTimeOfMs() int64 {
	return time.Now().UnixNano() / 1e6
}

func getRequest() Request {
	m := make(map[string]Request)
	return m["x"]
}

type Request struct {
}

var ruleInitTest = `
rule "rule_init_test" "rule_init"  salience 0
begin
 println(contextInt["a"])
 println(container.CurrentTimeOfMs())
 //println(contextInt["b"]) //此处如果是基础类型，如果没有对应的key，也应该返回对应的默认值
 println(container.EqualsIgnoreCase("a", "b"))
 println(container.IsNil(getRequest()))

end
`

func TestRuleInitTest(t *testing.T) {
	contextInt := make(map[string]int)
	contextString := make(map[string]string)
	contextInt["a"] = 1
	dataContext := context.NewDataContext()
	dataContext.Add("contextInt", &contextInt)
	dataContext.Add("contextString", &contextString)
	dataContext.Add("println", fmt.Println)
	dataContext.Add("getRequest", getRequest)

	//简化为这一行，就可以注入附着的所有函数
	container := &Container{}
	dataContext.Add("container", container)

	ruleBuilder := builder.NewRuleBuilder(dataContext)
	err := ruleBuilder.BuildRuleFromString(ruleInitTest)
	if err != nil {
		panic(err)
	}

	eng := engine.NewGengine()
	err = eng.Execute(ruleBuilder, true)
	if err != nil {
		panic(err)
	}
}

//pool 模式
func TestRuleInitTest_pool(t *testing.T) {

	//init和execute是异步的
	//init 放置于对象初始化方法中
	//execute放置于service中

	//init,apis放置一些与请求无关的状态函数
	apis := make(map[string]interface{})
	apis["println"] = fmt.Println
	apis["container"] = &Container{}
	apis["getRequest"] = getRequest
	pool, e1 := engine.NewGenginePool(1, 2, 1, ruleInitTest, apis)
	if e1 != nil {
		panic(e1)
	}

	//这里注入与请求相关的API，可以很好的隔离状态
	data := make(map[string]interface{})
	contextInt := make(map[string]int)
	contextInt["a"] = 1
	contextString := make(map[string]string)
	data["contextInt"] = &contextInt
	data["contextString"] = &contextString

	//execute 需要用什么模式，就调用什么方法
	e2, _ := pool.ExecuteRulesWithMultiInputWithSpecifiedEM(data)
	if e2 != nil {
		panic(e2)
	}

}
