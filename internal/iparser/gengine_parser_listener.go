package iparser

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/bilibili/gengine/internal/base"
	parser "github.com/bilibili/gengine/internal/iantlr/alr"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/golang-collections/collections/stack"
)

func NewGengineParserListener(ctx *base.KnowledgeContext) *GengineParserListener {
	return &GengineParserListener{
		Stack:            stack.New(),
		ParseErrors:      make([]string, 0),
		KnowledgeContext: ctx,
	}
}

type GengineParserListener struct {
	parser.BasegengineListener
	ParseErrors []string

	KnowledgeContext *base.KnowledgeContext
	Stack            *stack.Stack
	ruleName         string
	ruleDescription  string
	salience         int64
}

func (g *GengineParserListener) AddError(e error) {
	g.ParseErrors = append(g.ParseErrors, e.Error())
}

func (g *GengineParserListener) VisitTerminal(node antlr.TerminalNode) {}

func (g *GengineParserListener) VisitErrorNode(node antlr.ErrorNode) {
	g.AddError(errors.New(fmt.Sprintf("cannot recognize token : %s", node.GetText())))
}

func (g *GengineParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

func (g *GengineParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

func (g *GengineParserListener) EnterPrimary(ctx *parser.PrimaryContext) {}

func (g *GengineParserListener) ExitPrimary(ctx *parser.PrimaryContext) {}

func (g *GengineParserListener) EnterRuleEntity(ctx *parser.RuleEntityContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	entity := &base.RuleEntity{
		Salience: 0,
	}
	//init
	g.ruleName = ""
	g.ruleDescription = ""
	g.salience = 0
	g.Stack.Push(entity)
}

func (g *GengineParserListener) ExitRuleEntity(ctx *parser.RuleEntityContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	entity := g.Stack.Pop().(*base.RuleEntity)
	if _, ok := g.KnowledgeContext.RuleEntities[entity.RuleName]; ok {
		g.AddError(errors.New(fmt.Sprintf("already existed entity's name \"%s\"", entity.RuleName)))
		return
	}
	g.KnowledgeContext.RuleEntities[entity.RuleName] = entity
}

func (g *GengineParserListener) EnterRuleName(ctx *parser.RuleNameContext) {}

func (g *GengineParserListener) ExitRuleName(ctx *parser.RuleNameContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	text := ctx.GetText()
	ruleName := strings.Trim(text, "\"")
	if len(ruleName) == 0 {
		g.AddError(errors.New("rule name is \"\""))
		return
	}
	entity := g.Stack.Peek().(*base.RuleEntity)
	g.ruleName = ruleName
	entity.RuleName = ruleName
}

func (g *GengineParserListener) EnterSalience(ctx *parser.SalienceContext) {
}

func (g *GengineParserListener) ExitSalience(ctx *parser.SalienceContext) {

	if len(g.ParseErrors) > 0 {
		return
	}
	text := ctx.GetText()
	lower := strings.ToLower(text)
	is := strings.ReplaceAll(lower, "salience", "")
	i, err := strconv.ParseInt(is, 10, 64)
	if err != nil {
		g.AddError(errors.New(fmt.Sprintf("salience is not int, salience = \"%s\"", text)))
		return
	}
	g.salience = i
}

func (g *GengineParserListener) EnterRuleDescription(ctx *parser.RuleDescriptionContext) {}

func (g *GengineParserListener) ExitRuleDescription(ctx *parser.RuleDescriptionContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	text := ctx.GetText()
	ruleDescription := strings.Trim(text, "\"")
	entity := g.Stack.Peek().(*base.RuleEntity)
	entity.RuleDescription = ruleDescription
	g.ruleDescription = ruleDescription
}

func (g *GengineParserListener) EnterRuleContent(ctx *parser.RuleContentContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	ruleContent := &base.RuleContent{}
	g.Stack.Push(ruleContent)
}

func (g *GengineParserListener) ExitRuleContent(ctx *parser.RuleContentContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	ruleContent := g.Stack.Pop().(*base.RuleContent)
	entity := g.Stack.Peek().(*base.RuleEntity)
	entity.RuleContent = ruleContent
}

func (g *GengineParserListener) EnterConcStatement(ctx *parser.ConcStatementContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	concStatement := &base.ConcStatement{}
	g.Stack.Push(concStatement)
}

func (g *GengineParserListener) ExitConcStatement(ctx *parser.ConcStatementContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	concStatement := g.Stack.Pop().(*base.ConcStatement)
	statement := g.Stack.Peek().(*base.Statement)
	statement.ConcStatement = concStatement
}

func (g *GengineParserListener) EnterAssignment(ctx *parser.AssignmentContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	assignment := &base.Assignment{}
	g.Stack.Push(assignment)
}

func (g *GengineParserListener) ExitAssignment(ctx *parser.AssignmentContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	expr := g.Stack.Pop().(*base.Assignment)

	expr.Code = ctx.GetText()
	expr.LineNum = ctx.GetStart().GetLine()
	expr.Column = ctx.GetStart().GetColumn()
	expr.LineStop = ctx.GetStop().GetColumn()

	holder := g.Stack.Peek().(base.AssignmentHolder)
	err := holder.AcceptAssignment(expr)
	if err != nil {
		g.AddError(err)
	}
}

func (g *GengineParserListener) EnterMathExpression(ctx *parser.MathExpressionContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	me := &base.MathExpression{}
	g.Stack.Push(me)

}

func (g *GengineParserListener) ExitMathExpression(ctx *parser.MathExpressionContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	expr := g.Stack.Pop().(*base.MathExpression)

	expr.Code = ctx.GetText()
	expr.LineNum = ctx.GetStart().GetLine()
	expr.Column = ctx.GetStart().GetColumn()
	expr.LineStop = ctx.GetStop().GetColumn()

	holder := g.Stack.Peek().(base.MathExpressionHolder)
	err := holder.AcceptMathExpression(expr)
	if err != nil {
		g.AddError(err)
	}
}

func (g *GengineParserListener) EnterExpression(ctx *parser.ExpressionContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	expression := &base.Expression{}
	g.Stack.Push(expression)
}

func (g *GengineParserListener) ExitExpression(ctx *parser.ExpressionContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	expr := g.Stack.Pop().(*base.Expression)

	expr.Code = ctx.GetText()
	expr.LineNum = ctx.GetStart().GetLine()
	expr.Column = ctx.GetStart().GetColumn()
	expr.LineStop = ctx.GetStop().GetColumn()

	holder := g.Stack.Peek().(base.ExpressionHolder)
	err := holder.AcceptExpression(expr)
	if err != nil {
		g.AddError(err)
	}
}

func (g *GengineParserListener) EnterExpressionAtom(ctx *parser.ExpressionAtomContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	exprAtom := &base.ExpressionAtom{}
	g.Stack.Push(exprAtom)
}

func (g *GengineParserListener) ExitExpressionAtom(ctx *parser.ExpressionAtomContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	expr := g.Stack.Pop().(*base.ExpressionAtom)

	expr.Code = ctx.GetText()
	expr.LineNum = ctx.GetStart().GetLine()
	expr.Column = ctx.GetStart().GetColumn()
	expr.LineStop = ctx.GetStop().GetColumn()

	holder := g.Stack.Peek().(base.ExpressionAtomHolder)
	err := holder.AcceptExpressionAtom(expr)
	if err != nil {
		g.AddError(err)
	}
}

func (g *GengineParserListener) EnterMethodCall(ctx *parser.MethodCallContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	methodCall := &base.MethodCall{
		MethodName: ctx.DOTTEDNAME().GetText(),
	}
	g.Stack.Push(methodCall)
}

func (g *GengineParserListener) ExitMethodCall(ctx *parser.MethodCallContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	expr := g.Stack.Pop().(*base.MethodCall)

	expr.Code = ctx.GetText()
	expr.LineNum = ctx.GetStart().GetLine()
	expr.Column = ctx.GetStart().GetColumn()
	expr.LineStop = ctx.GetStop().GetColumn()

	holder := g.Stack.Peek().(base.MethodCallHolder)
	err := holder.AcceptMethodCall(expr)
	if err != nil {
		g.AddError(err)
	}
}

// EnterThreeLevelCall is called when production threeLevelCall is entered.
func (g *GengineParserListener) EnterThreeLevelCall(ctx *parser.ThreeLevelCallContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	threeLevelCall := &base.ThreeLevelCall{
		ThreeLevel: ctx.DOUBLEDOTTEDNAME().GetText(),
	}
	g.Stack.Push(threeLevelCall)
}

// ExitThreeLevelCall is called when production threeLevelCall is exited.
func (g *GengineParserListener) ExitThreeLevelCall(ctx *parser.ThreeLevelCallContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	expr := g.Stack.Pop().(*base.ThreeLevelCall)
	expr.Code = ctx.GetText()
	expr.LineNum = ctx.GetStart().GetLine()
	expr.Column = ctx.GetStart().GetColumn()
	expr.LineStop = ctx.GetStop().GetColumn()

	holder := g.Stack.Peek().(base.ThreeLevelCallHolder)
	err := holder.AcceptThreeLevelCall(expr)
	if err != nil {
		g.AddError(err)
	}
}

func (g *GengineParserListener) EnterFunctionCall(ctx *parser.FunctionCallContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	funcCall := &base.FunctionCall{
		FunctionName: ctx.SIMPLENAME().GetText(),
	}
	g.Stack.Push(funcCall)
}

func (g *GengineParserListener) ExitFunctionCall(ctx *parser.FunctionCallContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	expr := g.Stack.Pop().(*base.FunctionCall)

	expr.Code = ctx.GetText()
	expr.LineNum = ctx.GetStart().GetLine()
	expr.Column = ctx.GetStart().GetColumn()
	expr.LineStop = ctx.GetStop().GetColumn()
	holder := g.Stack.Peek().(base.FunctionCallHolder)
	err := holder.AcceptFunctionCall(expr)
	if err != nil {
		g.AddError(err)
	}
}

func (g *GengineParserListener) EnterFunctionArgs(ctx *parser.FunctionArgsContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	funcArg := &base.Args{
		ArgList: make([]*base.Arg, 0),
	}
	g.Stack.Push(funcArg)
}

func (g *GengineParserListener) ExitFunctionArgs(ctx *parser.FunctionArgsContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	expr := g.Stack.Pop().(*base.Args)
	argHolder := g.Stack.Peek().(base.ArgsHolder)
	err := argHolder.AcceptArgs(expr)
	if err != nil {
		g.AddError(err)
	}
}

func (g *GengineParserListener) EnterLogicalOperator(ctx *parser.LogicalOperatorContext) {}

func (g *GengineParserListener) ExitLogicalOperator(ctx *parser.LogicalOperatorContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	expr := g.Stack.Peek().(*base.Expression)
	// && and ||
	expr.LogicalOperator = ctx.GetText()
}

func (g *GengineParserListener) EnterNotOperator(ctx *parser.NotOperatorContext) {}

func (g *GengineParserListener) ExitNotOperator(ctx *parser.NotOperatorContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	expr := g.Stack.Peek().(*base.Expression)
	// !
	expr.NotOperator = ctx.GetText()
}

func (g *GengineParserListener) EnterVariable(ctx *parser.VariableContext) {}

func (g *GengineParserListener) ExitVariable(ctx *parser.VariableContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	varName := ctx.GetText()
	holder := g.Stack.Peek().(base.VariableHolder)
	err := holder.AcceptVariable(varName)
	if err != nil {
		g.AddError(err)
	}
}

func (g *GengineParserListener) EnterMathPmOperator(ctx *parser.MathPmOperatorContext) {}

func (g *GengineParserListener) ExitMathPmOperator(ctx *parser.MathPmOperatorContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	expr := g.Stack.Peek().(*base.MathExpression)
	// + -
	expr.MathPmOperator = ctx.GetText()
}

func (g *GengineParserListener) EnterMathMdOperator(ctx *parser.MathMdOperatorContext) {}

func (g *GengineParserListener) ExitMathMdOperator(ctx *parser.MathMdOperatorContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	expr := g.Stack.Peek().(*base.MathExpression)
	// * /
	expr.MathMdOperator = ctx.GetText()
}

func (g *GengineParserListener) EnterComparisonOperator(ctx *parser.ComparisonOperatorContext) {}

func (g *GengineParserListener) ExitComparisonOperator(ctx *parser.ComparisonOperatorContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	expr := g.Stack.Peek().(*base.Expression)
	// ==  !=  <  > <= >=
	expr.ComparisonOperator = ctx.GetText()
}

func (g *GengineParserListener) EnterConstant(ctx *parser.ConstantContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	cons := &base.Constant{}
	g.Stack.Push(cons)
}

func (g *GengineParserListener) ExitConstant(ctx *parser.ConstantContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	cons := g.Stack.Pop().(*base.Constant)
	holder := g.Stack.Peek().(base.ConstantHolder)
	err := holder.AcceptConstant(cons)
	if err != nil {
		g.AddError(err)
	}
}

func (g *GengineParserListener) EnterStringLiteral(ctx *parser.StringLiteralContext) {}

func (g *GengineParserListener) ExitStringLiteral(ctx *parser.StringLiteralContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	holder := g.Stack.Peek().(base.StringHolder)
	text := ctx.GetText()
	txt := strings.Trim(text, "\"")
	if reflect.TypeOf(holder).String() == "*base.MapVar" {
		if txt == "" {
			g.AddError(errors.New("MAP key should not be null string"))
		}
	}
	err := holder.AcceptString(txt)
	if err != nil {
		g.AddError(err)
	}
}

func (g *GengineParserListener) EnterBooleanLiteral(ctx *parser.BooleanLiteralContext) {}

func (g *GengineParserListener) ExitBooleanLiteral(ctx *parser.BooleanLiteralContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	cons := g.Stack.Peek().(*base.Constant)
	b, e := strconv.ParseBool(ctx.GetText())
	if e != nil {
		g.AddError(e)
		return
	}
	cons.ConstantValue = reflect.ValueOf(b)
}

func (g *GengineParserListener) EnterRealLiteral(ctx *parser.RealLiteralContext) {}

func (g *GengineParserListener) ExitRealLiteral(ctx *parser.RealLiteralContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	cons := g.Stack.Peek().(*base.Constant)
	flo, err := strconv.ParseFloat(ctx.GetText(), 64)
	if err != nil {
		g.AddError(errors.New(fmt.Sprintf("string to float conversion error. String is not real type '%s'", ctx.GetText())))
		return
	}
	cons.ConstantValue = reflect.ValueOf(flo)
}

func (g *GengineParserListener) EnterIfStmt(ctx *parser.IfStmtContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	ifStmt := &base.IfStmt{}
	g.Stack.Push(ifStmt)
}

func (g *GengineParserListener) ExitIfStmt(ctx *parser.IfStmtContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	ifStmt := g.Stack.Pop().(*base.IfStmt)
	statement := g.Stack.Peek().(*base.Statement)
	statement.IfStmt = ifStmt
}

// 语法树进入for节点
func (g *GengineParserListener) EnterForStmt(ctx *parser.ForStmtContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	forStmt := &base.ForStmt{}
	g.Stack.Push(forStmt)
}

// 语法树退出for节点
func (g *GengineParserListener) ExitForStmt(ctx *parser.ForStmtContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	forStmt := g.Stack.Pop().(*base.ForStmt)
	statement := g.Stack.Peek().(*base.Statement)
	statement.ForStmt = forStmt
}

// 语法树进入forRange节点
func (g *GengineParserListener) EnterForRangeStmt(ctx *parser.ForRangeStmtContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	forRangeStmt := &base.ForRangeStmt{}
	g.Stack.Push(forRangeStmt)
}

// 语法树退出forRange节点
func (g *GengineParserListener) ExitForRangeStmt(ctx *parser.ForRangeStmtContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	forRangeStmt := g.Stack.Pop().(*base.ForRangeStmt)
	statement := g.Stack.Peek().(*base.Statement)
	statement.ForRangeStmt = forRangeStmt
}

// 语法树进入brak节点
func (g *GengineParserListener) EnterBreakStmt(ctx *parser.BreakStmtContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	breakStmt := &base.BreakStmt{}
	g.Stack.Push(breakStmt)
}

// 语法树退出brak节点
func (g *GengineParserListener) ExitBreakStmt(ctx *parser.BreakStmtContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	breakStmt := g.Stack.Pop().(*base.BreakStmt)
	statement := g.Stack.Peek().(*base.Statement)
	statement.BreakStmt = breakStmt
}

// 语法树进入continue节点
func (g *GengineParserListener) EnterContinueStmt(ctx *parser.ContinueStmtContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	continueStmt := &base.ContinueStmt{}
	g.Stack.Push(continueStmt)
}

// 语法树退出continue节点
func (g *GengineParserListener) ExitContinueStmt(ctx *parser.ContinueStmtContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	continueStmt := g.Stack.Pop().(*base.ContinueStmt)
	statement := g.Stack.Peek().(*base.Statement)
	statement.ContinueStmt = continueStmt
}

func (g *GengineParserListener) EnterStatement(ctx *parser.StatementContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	statement := &base.Statement{}
	g.Stack.Push(statement)
}

func (g *GengineParserListener) ExitStatement(ctx *parser.StatementContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	statement := g.Stack.Pop().(*base.Statement)
	statements := g.Stack.Peek().(*base.Statements)
	statements.StatementList = append(statements.StatementList, statement)
}

func (g *GengineParserListener) EnterStatements(ctx *parser.StatementsContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	statements := &base.Statements{
		StatementList: make([]*base.Statement, 0),
	}
	g.Stack.Push(statements)
}

func (g *GengineParserListener) ExitStatements(ctx *parser.StatementsContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	statements := g.Stack.Pop().(*base.Statements)
	holder := g.Stack.Peek().(base.StatementsHolder)
	err := holder.AcceptStatements(statements)
	if err != nil {
		g.AddError(err)
	}
}

func (g *GengineParserListener) EnterAssignOperator(ctx *parser.AssignOperatorContext) {}

func (g *GengineParserListener) ExitAssignOperator(ctx *parser.AssignOperatorContext) {

	if len(g.ParseErrors) > 0 {
		return
	}
	expr := g.Stack.Peek().(*base.Assignment)
	expr.AssignOperator = ctx.GetText()
}

func (g *GengineParserListener) EnterElseIfStmt(ctx *parser.ElseIfStmtContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	elseIfStmt := &base.ElseIfStmt{}
	g.Stack.Push(elseIfStmt)
}

func (g *GengineParserListener) ExitElseIfStmt(ctx *parser.ElseIfStmtContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	elseIfStmt := g.Stack.Pop().(*base.ElseIfStmt)
	ifStmt := g.Stack.Peek().(*base.IfStmt)
	ifStmt.ElseIfStmtList = append(ifStmt.ElseIfStmtList, elseIfStmt)
}

func (g *GengineParserListener) EnterElseStmt(ctx *parser.ElseStmtContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	elseStmt := &base.ElseStmt{}
	g.Stack.Push(elseStmt)
}

func (g *GengineParserListener) ExitElseStmt(ctx *parser.ElseStmtContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	elseStmt := g.Stack.Pop().(*base.ElseStmt)
	ifStmt := g.Stack.Peek().(*base.IfStmt)
	ifStmt.ElseStmt = elseStmt
}

func (g *GengineParserListener) ExitInteger(ctx *parser.IntegerContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	val, err := strconv.ParseInt(ctx.GetText(), 10, 64)
	if err != nil {
		g.AddError(err)
		return
	}
	holder := g.Stack.Peek().(base.IntegerHolder)
	err = holder.AcceptInteger(val)
	if err != nil {
		g.AddError(err)
	}
}
func (g *GengineParserListener) EnterInteger(ctx *parser.IntegerContext) {}

func (g *GengineParserListener) EnterAtName(ctx *parser.AtNameContext) {}

func (g *GengineParserListener) ExitAtName(ctx *parser.AtNameContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	holder := g.Stack.Peek().(base.AtNameHolder)
	err := holder.AcceptName(strings.ReplaceAll(g.ruleName, "\"", ""))
	if err != nil {
		g.AddError(err)
	}
}
func (g *GengineParserListener) EnterAtDesc(ctx *parser.AtDescContext) {}

func (g *GengineParserListener) ExitAtDesc(ctx *parser.AtDescContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	holder := g.Stack.Peek().(base.AtDescHolder)
	err := holder.AcceptDesc(strings.ReplaceAll(g.ruleDescription, "\"", ""))
	if err != nil {
		g.AddError(err)
	}
}

func (g *GengineParserListener) EnterAtId(ctx *parser.AtIdContext) {}

func (g *GengineParserListener) ExitAtId(ctx *parser.AtIdContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	holder := g.Stack.Peek().(base.AtIdHolder)
	i, e := strconv.ParseInt(strings.Trim(g.ruleName, " "), 10, 64)
	if e != nil {
		err := holder.AcceptId(0)
		if err != nil {
			g.AddError(err)
		}
		return
	}

	err := holder.AcceptId(i)
	if err != nil {
		g.AddError(err)
	}
}

func (g *GengineParserListener) EnterAtSal(ctx *parser.AtSalContext) {}

func (g *GengineParserListener) ExitAtSal(ctx *parser.AtSalContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	holder := g.Stack.Peek().(base.AtSalienceHolder)
	err := holder.AcceptSalience(g.salience)
	if err != nil {
		g.AddError(err)
	}
}

func (g *GengineParserListener) EnterMapVar(ctx *parser.MapVarContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	mapVar := &base.MapVar{}
	g.Stack.Push(mapVar)

}

func (g *GengineParserListener) ExitMapVar(ctx *parser.MapVarContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	mapVar := g.Stack.Pop().(*base.MapVar)
	holder := g.Stack.Peek().(base.MapVarHolder)
	err := holder.AcceptMapVar(mapVar)
	if err != nil {
		g.AddError(err)
	}
}

func (g *GengineParserListener) EnterReturnStmt(c *parser.ReturnStmtContext) {
	if len(g.ParseErrors) > 0 {
		return
	}
	rs := &base.ReturnStatement{}
	g.Stack.Push(rs)
}

func (g *GengineParserListener) ExitReturnStmt(c *parser.ReturnStmtContext) {
	if len(g.ParseErrors) > 0 {
		return
	}

	rs := g.Stack.Pop().(*base.ReturnStatement)
	stats := g.Stack.Peek().(*base.Statements)
	stats.ReturnStatement = rs
}
