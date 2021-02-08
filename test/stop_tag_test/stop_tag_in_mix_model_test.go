package stop_tag

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"testing"
)

const mix_rule = `
rule "the most import rule" "use in white list pass" salience 1000
BEGIN
    if InWhitelist(User.UserNo) {
		println("User.UserNo =",  User.UserNo, " in white list, passed, will not to execute last rules")
		stag.StopTag = true      //if true, it will not to execute  the last rules, if stag.StopTag= false, it will continue to execute last rules
	} else {
		println("User.UserNo =",  User.UserNo, "not in white list, continue to execute last rules...")
	}
END

rule "the sub import rule 1" "1" salience 5
BEGIN
	println("the sub import rule 1")
END

rule "the sub import rule 2" "2" salience 4
BEGIN
	println("the sub import rule 2")
END

rule "the sub import rule 3" "3" salience 6
BEGIN
	println("the sub import rule 3")
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

func Test_stop_tag_in_mix_model(t *testing.T) {

	dataContext := context.NewDataContext()
	dataContext.Add("InWhitelist", InWhitelist)
	dataContext.Add("InBlacklist", InBlacklist)
	dataContext.Add("println", fmt.Println)

	ruleBuilder := builder.NewRuleBuilder(dataContext)

	e1 := ruleBuilder.BuildRuleFromString(mix_rule)
	if e1 != nil {
		panic(e1)
	}

	//must default false
	stag := &engine.Stag{StopTag: false}
	dataContext.Add("stag", stag)

	user := &User{UserNo: 1000} // change this to test different conditions
	dataContext.Add("User", user)

	eng := engine.NewGengine()
	e1 = eng.ExecuteMixModelWithStopTagDirect(ruleBuilder, stag)
	if e1 != nil {
		panic(e1)
	}
}
