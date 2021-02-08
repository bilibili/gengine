// Code generated from /Users/renyunyi/go_project/gengine/internal/iantlr/gengine.g4 by ANTLR 4.9. DO NOT EDIT.

package parser // gengine

import "github.com/antlr/antlr4/runtime/Go/antlr"

// gengineListener is a complete listener for a parse tree produced by gengineParser.
type gengineListener interface {
	antlr.ParseTreeListener

	// EnterPrimary is called when entering the primary production.
	EnterPrimary(c *PrimaryContext)

	// EnterRuleEntity is called when entering the ruleEntity production.
	EnterRuleEntity(c *RuleEntityContext)

	// EnterRuleName is called when entering the ruleName production.
	EnterRuleName(c *RuleNameContext)

	// EnterRuleDescription is called when entering the ruleDescription production.
	EnterRuleDescription(c *RuleDescriptionContext)

	// EnterSalience is called when entering the salience production.
	EnterSalience(c *SalienceContext)

	// EnterRuleContent is called when entering the ruleContent production.
	EnterRuleContent(c *RuleContentContext)

	// EnterStatements is called when entering the statements production.
	EnterStatements(c *StatementsContext)

	// EnterStatement is called when entering the statement production.
	EnterStatement(c *StatementContext)

	// EnterConcStatement is called when entering the concStatement production.
	EnterConcStatement(c *ConcStatementContext)

	// EnterExpression is called when entering the expression production.
	EnterExpression(c *ExpressionContext)

	// EnterMathExpression is called when entering the mathExpression production.
	EnterMathExpression(c *MathExpressionContext)

	// EnterExpressionAtom is called when entering the expressionAtom production.
	EnterExpressionAtom(c *ExpressionAtomContext)

	// EnterAssignment is called when entering the assignment production.
	EnterAssignment(c *AssignmentContext)

	// EnterReturnStmt is called when entering the returnStmt production.
	EnterReturnStmt(c *ReturnStmtContext)

	// EnterIfStmt is called when entering the ifStmt production.
	EnterIfStmt(c *IfStmtContext)

	// EnterElseIfStmt is called when entering the elseIfStmt production.
	EnterElseIfStmt(c *ElseIfStmtContext)

	// EnterElseStmt is called when entering the elseStmt production.
	EnterElseStmt(c *ElseStmtContext)

	// EnterConstant is called when entering the constant production.
	EnterConstant(c *ConstantContext)

	// EnterFunctionArgs is called when entering the functionArgs production.
	EnterFunctionArgs(c *FunctionArgsContext)

	// EnterInteger is called when entering the integer production.
	EnterInteger(c *IntegerContext)

	// EnterRealLiteral is called when entering the realLiteral production.
	EnterRealLiteral(c *RealLiteralContext)

	// EnterStringLiteral is called when entering the stringLiteral production.
	EnterStringLiteral(c *StringLiteralContext)

	// EnterBooleanLiteral is called when entering the booleanLiteral production.
	EnterBooleanLiteral(c *BooleanLiteralContext)

	// EnterFunctionCall is called when entering the functionCall production.
	EnterFunctionCall(c *FunctionCallContext)

	// EnterMethodCall is called when entering the methodCall production.
	EnterMethodCall(c *MethodCallContext)

	// EnterVariable is called when entering the variable production.
	EnterVariable(c *VariableContext)

	// EnterMathPmOperator is called when entering the mathPmOperator production.
	EnterMathPmOperator(c *MathPmOperatorContext)

	// EnterMathMdOperator is called when entering the mathMdOperator production.
	EnterMathMdOperator(c *MathMdOperatorContext)

	// EnterComparisonOperator is called when entering the comparisonOperator production.
	EnterComparisonOperator(c *ComparisonOperatorContext)

	// EnterLogicalOperator is called when entering the logicalOperator production.
	EnterLogicalOperator(c *LogicalOperatorContext)

	// EnterAssignOperator is called when entering the assignOperator production.
	EnterAssignOperator(c *AssignOperatorContext)

	// EnterNotOperator is called when entering the notOperator production.
	EnterNotOperator(c *NotOperatorContext)

	// EnterMapVar is called when entering the mapVar production.
	EnterMapVar(c *MapVarContext)

	// EnterAtName is called when entering the atName production.
	EnterAtName(c *AtNameContext)

	// EnterAtId is called when entering the atId production.
	EnterAtId(c *AtIdContext)

	// EnterAtDesc is called when entering the atDesc production.
	EnterAtDesc(c *AtDescContext)

	// EnterAtSal is called when entering the atSal production.
	EnterAtSal(c *AtSalContext)

	// ExitPrimary is called when exiting the primary production.
	ExitPrimary(c *PrimaryContext)

	// ExitRuleEntity is called when exiting the ruleEntity production.
	ExitRuleEntity(c *RuleEntityContext)

	// ExitRuleName is called when exiting the ruleName production.
	ExitRuleName(c *RuleNameContext)

	// ExitRuleDescription is called when exiting the ruleDescription production.
	ExitRuleDescription(c *RuleDescriptionContext)

	// ExitSalience is called when exiting the salience production.
	ExitSalience(c *SalienceContext)

	// ExitRuleContent is called when exiting the ruleContent production.
	ExitRuleContent(c *RuleContentContext)

	// ExitStatements is called when exiting the statements production.
	ExitStatements(c *StatementsContext)

	// ExitStatement is called when exiting the statement production.
	ExitStatement(c *StatementContext)

	// ExitConcStatement is called when exiting the concStatement production.
	ExitConcStatement(c *ConcStatementContext)

	// ExitExpression is called when exiting the expression production.
	ExitExpression(c *ExpressionContext)

	// ExitMathExpression is called when exiting the mathExpression production.
	ExitMathExpression(c *MathExpressionContext)

	// ExitExpressionAtom is called when exiting the expressionAtom production.
	ExitExpressionAtom(c *ExpressionAtomContext)

	// ExitAssignment is called when exiting the assignment production.
	ExitAssignment(c *AssignmentContext)

	// ExitReturnStmt is called when exiting the returnStmt production.
	ExitReturnStmt(c *ReturnStmtContext)

	// ExitIfStmt is called when exiting the ifStmt production.
	ExitIfStmt(c *IfStmtContext)

	// ExitElseIfStmt is called when exiting the elseIfStmt production.
	ExitElseIfStmt(c *ElseIfStmtContext)

	// ExitElseStmt is called when exiting the elseStmt production.
	ExitElseStmt(c *ElseStmtContext)

	// ExitConstant is called when exiting the constant production.
	ExitConstant(c *ConstantContext)

	// ExitFunctionArgs is called when exiting the functionArgs production.
	ExitFunctionArgs(c *FunctionArgsContext)

	// ExitInteger is called when exiting the integer production.
	ExitInteger(c *IntegerContext)

	// ExitRealLiteral is called when exiting the realLiteral production.
	ExitRealLiteral(c *RealLiteralContext)

	// ExitStringLiteral is called when exiting the stringLiteral production.
	ExitStringLiteral(c *StringLiteralContext)

	// ExitBooleanLiteral is called when exiting the booleanLiteral production.
	ExitBooleanLiteral(c *BooleanLiteralContext)

	// ExitFunctionCall is called when exiting the functionCall production.
	ExitFunctionCall(c *FunctionCallContext)

	// ExitMethodCall is called when exiting the methodCall production.
	ExitMethodCall(c *MethodCallContext)

	// ExitVariable is called when exiting the variable production.
	ExitVariable(c *VariableContext)

	// ExitMathPmOperator is called when exiting the mathPmOperator production.
	ExitMathPmOperator(c *MathPmOperatorContext)

	// ExitMathMdOperator is called when exiting the mathMdOperator production.
	ExitMathMdOperator(c *MathMdOperatorContext)

	// ExitComparisonOperator is called when exiting the comparisonOperator production.
	ExitComparisonOperator(c *ComparisonOperatorContext)

	// ExitLogicalOperator is called when exiting the logicalOperator production.
	ExitLogicalOperator(c *LogicalOperatorContext)

	// ExitAssignOperator is called when exiting the assignOperator production.
	ExitAssignOperator(c *AssignOperatorContext)

	// ExitNotOperator is called when exiting the notOperator production.
	ExitNotOperator(c *NotOperatorContext)

	// ExitMapVar is called when exiting the mapVar production.
	ExitMapVar(c *MapVarContext)

	// ExitAtName is called when exiting the atName production.
	ExitAtName(c *AtNameContext)

	// ExitAtId is called when exiting the atId production.
	ExitAtId(c *AtIdContext)

	// ExitAtDesc is called when exiting the atDesc production.
	ExitAtDesc(c *AtDescContext)

	// ExitAtSal is called when exiting the atSal production.
	ExitAtSal(c *AtSalContext)
}
