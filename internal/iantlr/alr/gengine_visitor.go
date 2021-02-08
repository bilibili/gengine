// Code generated from /Users/renyunyi/go_project/gengine/internal/iantlr/gengine.g4 by ANTLR 4.9. DO NOT EDIT.

package parser // gengine

import "github.com/antlr/antlr4/runtime/Go/antlr"
// A complete Visitor for a parse tree produced by gengineParser.
type gengineVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by gengineParser#primary.
	VisitPrimary(ctx *PrimaryContext) interface{}

	// Visit a parse tree produced by gengineParser#ruleEntity.
	VisitRuleEntity(ctx *RuleEntityContext) interface{}

	// Visit a parse tree produced by gengineParser#ruleName.
	VisitRuleName(ctx *RuleNameContext) interface{}

	// Visit a parse tree produced by gengineParser#ruleDescription.
	VisitRuleDescription(ctx *RuleDescriptionContext) interface{}

	// Visit a parse tree produced by gengineParser#salience.
	VisitSalience(ctx *SalienceContext) interface{}

	// Visit a parse tree produced by gengineParser#ruleContent.
	VisitRuleContent(ctx *RuleContentContext) interface{}

	// Visit a parse tree produced by gengineParser#statements.
	VisitStatements(ctx *StatementsContext) interface{}

	// Visit a parse tree produced by gengineParser#statement.
	VisitStatement(ctx *StatementContext) interface{}

	// Visit a parse tree produced by gengineParser#concStatement.
	VisitConcStatement(ctx *ConcStatementContext) interface{}

	// Visit a parse tree produced by gengineParser#expression.
	VisitExpression(ctx *ExpressionContext) interface{}

	// Visit a parse tree produced by gengineParser#mathExpression.
	VisitMathExpression(ctx *MathExpressionContext) interface{}

	// Visit a parse tree produced by gengineParser#expressionAtom.
	VisitExpressionAtom(ctx *ExpressionAtomContext) interface{}

	// Visit a parse tree produced by gengineParser#assignment.
	VisitAssignment(ctx *AssignmentContext) interface{}

	// Visit a parse tree produced by gengineParser#returnStmt.
	VisitReturnStmt(ctx *ReturnStmtContext) interface{}

	// Visit a parse tree produced by gengineParser#ifStmt.
	VisitIfStmt(ctx *IfStmtContext) interface{}

	// Visit a parse tree produced by gengineParser#elseIfStmt.
	VisitElseIfStmt(ctx *ElseIfStmtContext) interface{}

	// Visit a parse tree produced by gengineParser#elseStmt.
	VisitElseStmt(ctx *ElseStmtContext) interface{}

	// Visit a parse tree produced by gengineParser#constant.
	VisitConstant(ctx *ConstantContext) interface{}

	// Visit a parse tree produced by gengineParser#functionArgs.
	VisitFunctionArgs(ctx *FunctionArgsContext) interface{}

	// Visit a parse tree produced by gengineParser#integer.
	VisitInteger(ctx *IntegerContext) interface{}

	// Visit a parse tree produced by gengineParser#realLiteral.
	VisitRealLiteral(ctx *RealLiteralContext) interface{}

	// Visit a parse tree produced by gengineParser#stringLiteral.
	VisitStringLiteral(ctx *StringLiteralContext) interface{}

	// Visit a parse tree produced by gengineParser#booleanLiteral.
	VisitBooleanLiteral(ctx *BooleanLiteralContext) interface{}

	// Visit a parse tree produced by gengineParser#functionCall.
	VisitFunctionCall(ctx *FunctionCallContext) interface{}

	// Visit a parse tree produced by gengineParser#methodCall.
	VisitMethodCall(ctx *MethodCallContext) interface{}

	// Visit a parse tree produced by gengineParser#variable.
	VisitVariable(ctx *VariableContext) interface{}

	// Visit a parse tree produced by gengineParser#mathPmOperator.
	VisitMathPmOperator(ctx *MathPmOperatorContext) interface{}

	// Visit a parse tree produced by gengineParser#mathMdOperator.
	VisitMathMdOperator(ctx *MathMdOperatorContext) interface{}

	// Visit a parse tree produced by gengineParser#comparisonOperator.
	VisitComparisonOperator(ctx *ComparisonOperatorContext) interface{}

	// Visit a parse tree produced by gengineParser#logicalOperator.
	VisitLogicalOperator(ctx *LogicalOperatorContext) interface{}

	// Visit a parse tree produced by gengineParser#assignOperator.
	VisitAssignOperator(ctx *AssignOperatorContext) interface{}

	// Visit a parse tree produced by gengineParser#notOperator.
	VisitNotOperator(ctx *NotOperatorContext) interface{}

	// Visit a parse tree produced by gengineParser#mapVar.
	VisitMapVar(ctx *MapVarContext) interface{}

	// Visit a parse tree produced by gengineParser#atName.
	VisitAtName(ctx *AtNameContext) interface{}

	// Visit a parse tree produced by gengineParser#atId.
	VisitAtId(ctx *AtIdContext) interface{}

	// Visit a parse tree produced by gengineParser#atDesc.
	VisitAtDesc(ctx *AtDescContext) interface{}

	// Visit a parse tree produced by gengineParser#atSal.
	VisitAtSal(ctx *AtSalContext) interface{}

}