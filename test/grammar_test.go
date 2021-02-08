package test

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/internal/base"
	parser "github.com/bilibili/gengine/internal/iantlr/alr"
	"github.com/bilibili/gengine/internal/iparser"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"testing"
)

const (
	rule3 = `
rule "name test" "i can" salience 0
BEGIN
		if 7 == User.GetNum(7){
			User.Age = User.GetNum(89767) + 10000000
			User.Print("6666")
		}else{
			User.Name = "yyyy"
		}
END`
)

func Test_grammar(t *testing.T) {
	in := antlr.NewInputStream(rule3)
	lexer := parser.NewgengineLexer(in)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	knowledgeBase := base.NewKnowledgeContext()
	listener := iparser.NewGengineParserListener(knowledgeBase)

	psr := parser.NewgengineParser(stream)
	psr.BuildParseTrees = true
	antlr.ParseTreeWalkerDefault.Walk(listener, psr.Primary())
	if len(listener.ParseErrors) > 0 {
		fmt.Println(listener.ParseErrors)
	}
}

func Test_base_msg(t *testing.T) {

	dataContext := context.NewDataContext()
	//init rule engine
	ruleBuilder := builder.NewRuleBuilder(dataContext)

	//读取规则
	err := ruleBuilder.BuildRuleFromString(rule3)
	if err != nil {
		panic(err)
	}

	rule := ruleBuilder.Kc.RuleEntities
	for k, v := range rule {
		fmt.Println(k)
		fmt.Println(v.RuleName)
		fmt.Println(v.RuleDescription)
		fmt.Println(v.Salience)
		fmt.Println(v.RuleContent)
	}
}

/**
测试语法错误
*/
func Test_err(t *testing.T) {

	in := antlr.NewInputStream(rule3)
	lexer := parser.NewgengineLexer(in)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	knowledgeBase := base.NewKnowledgeContext()
	listener := iparser.NewGengineParserListener(knowledgeBase)

	psr := parser.NewgengineParser(stream)
	psr.BuildParseTrees = true

	//test -- here
	errListener := iparser.NewGengineErrorListener()
	psr.AddErrorListener(errListener)
	//test -- here

	antlr.ParseTreeWalkerDefault.Walk(listener, psr.Primary())
	fmt.Println(errListener.GrammarErrors)
}
