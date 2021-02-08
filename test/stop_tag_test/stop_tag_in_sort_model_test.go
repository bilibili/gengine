package stop_tag_test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
)

const sort_rule = `
rule "whitelist" "use in white list pass" salience 10
BEGIN
    if InWhitelist(User.UserNo) {
		println("User.UserNo =", User.UserNo, "in white list, passed, will not to execute rule 'blacklist'.")
		stag.StopTag = true        // if true, it will not to execute  the last rule "blacklist", if stag.StopTag= false, it will continue to execute last rule "blacklist"
	} else {
		println("User.UserNo =", User.UserNo, "not in white list, continue to execute last rules...")
	}
END
rule "blacklist" "use in black list deny" salience 5
BEGIN
	if InBlacklist(User.UserNo) {
		println(User.UserNo, "in black list, denied...")
	}else{
		println("continue..2..")
	}
END
`

func InWhitelist(uid int) bool {
	if uid > 100 {
		return true
	}
	return false
}

func InBlacklist(uid int) bool {
	if uid <= 100 {
		return true
	}
	return false
}

type User struct {
	UserNo int
}

func Test_stop_tag_in_sort_model(t *testing.T) {

	dataContext := context.NewDataContext()
	dataContext.Add("InWhitelist", InWhitelist)
	dataContext.Add("InBlacklist", InBlacklist)
	dataContext.Add("println", fmt.Println)

	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	e1 := ruleBuilder.BuildRuleFromString(sort_rule)
	if e1 != nil {
		panic(e1)
	}

	//must default false
	stag := &engine.Stag{StopTag: false}
	dataContext.Add("stag", stag)

	user := &User{UserNo: 10} // change this to test different conditions
	dataContext.Add("User", user)

	eng := engine.NewGengine()
	e2 := eng.ExecuteWithStopTagDirect(ruleBuilder, true, stag)
	if e2 != nil {
		panic(e2)
	}
}
