package test

import (
	"fmt"
	"github.com/bilibili/gengine/engine"
	"testing"
)

/*type User struct {
	Name string
	Age int
}*/

func Test_list_data_check_by_rules(t *testing.T) {

	// there are many rules you add
	var rules = `
rule "1" 
begin
// you want to do 
end

rule "2" 
begin
// you want to do 
end

// rules "n"...
`
	//apis use to inject non-state function or service without worry about thread-safety
	apis := make(map[string]interface{})
	apis["println"] = fmt.Println
	//just init once !
	pool, e := engine.NewGenginePool(10, 20, 1, rules, apis)
	if e != nil {
		//init pool failed
		panic(e)
	}

	// there is your a list of data(users)
	users := make([]*User, 100)

	for _, user := range users {
		tmpApisOrData := make(map[string]interface{})
		tmpApisOrData["user"] = user
		//execute rules
		//all methods the gengine pool supply is thread-safety
		e, _ := pool.Execute(tmpApisOrData, true)
		if e != nil {
			//log or return the err
			//errors.New(fmt.Sprintf("execute rules err:%+v", e))
		}
		// err== nil to do you want to do
	}

	// or concurrent to execute
	/*
		for _, user := range users  {
			u := user
			go func() {
				tmpApisOrData := make(map[string]interface{})
				tmpApisOrData["user"] = u
				//execute rules
				//all methods the gengine pool supply is thread-safety
				e, _ := pool.Execute(tmpApisOrData, true)
				if e != nil {
					//log or return the err
					//errors.New(fmt.Sprintf("execute rules err:%+v", e))
				}
				// err== nil to do you want to do
			}()
		}
	*/
}
