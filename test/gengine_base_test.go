package test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
	"time"
)

type User struct {
	Name string
	Age  int64
	Male bool
}

func (u *User) GetNum(i int64) int64 {
	return i
}

func (u *User) Print(s string) {
	fmt.Println(s)
}

func (u *User) Say() {
	fmt.Println("hello world")
}

const (
	base_rule = `
rule "测试" "测试描述"  salience 0 
begin
		// 重命名函数 测试; @name represent the rule name "测试"
		Sout(@name)
		// 普通函数 测试
		Hello()
		//结构提方法 测试
		User.Say()
		// if
		if !(7 == User.GetNum(7)) || !(7 > 8)  {
			//自定义变量 和 加法 测试
			variable = "hello" + (" world" + "zeze")
			// 加法 与 内建函数 测试 ; @name is just a string  
			User.Name = "hhh" + strconv.FormatInt(10, 10) + "@name"
			//结构体属性、方法调用 和 除法 测试
			User.Age = User.GetNum(8976) / 1000+ 3*(1+1) 
			//布尔值设置 测试
			User.Male = false
			//规则内自定义变量调用 测试
			User.Print(variable)
			//float测试	也支持科学计数法		
			f = 9.56			
			PrintReal(f)
			//嵌套if-else测试
			if false	{
				Sout("嵌套if测试")
			}else{
				Sout("嵌套else测试")
			}
		}else{ //else
			//字符串设置 测试
			User.Name = "yyyy"
		}
		
		if true {
			Sout("if true ")
		}

		if true{}else{}
end`
)

func Hello() {
	fmt.Println("hello")
}

func PrintReal(real float64) {
	fmt.Println(real)
}

func exe(user *User) {
	dataContext := context.NewDataContext()
	//inject struct
	strconv := &StrconvWrapper{}
	dataContext.Add("strconv", strconv)
	dataContext.Add("User", user)
	//rename and inject
	dataContext.Add("Sout", fmt.Println)
	//直接注入函数
	dataContext.Add("Hello", Hello)
	dataContext.Add("PrintReal", PrintReal)

	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	//读取规则
	start1 := time.Now().UnixNano()
	err := ruleBuilder.BuildRuleFromString(base_rule)
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
	println(fmt.Sprintf("user.Age=%d,Name=%s,Male=%t", user.Age, user.Name, user.Male))

}

func Test_Base(t *testing.T) {
	user := &User{
		Name: "Calo",
		Age:  0,
		Male: true,
	}
	exe(user)
}
