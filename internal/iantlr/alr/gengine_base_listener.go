// Code generated from /Users/renyunyi/go_project/gengine/internal/iantlr/gengine.g4 by ANTLR 4.9. DO NOT EDIT.

package parser // gengine

import "github.com/antlr/antlr4/runtime/Go/antlr"

// BasegengineListener is a complete listener for a parse tree produced by gengineParser.
type BasegengineListener struct{}

var _ gengineListener = &BasegengineListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BasegengineListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BasegengineListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BasegengineListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BasegengineListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterPrimary is called when production primary is entered.
func (s *BasegengineListener) EnterPrimary(ctx *PrimaryContext) {}

// ExitPrimary is called when production primary is exited.
func (s *BasegengineListener) ExitPrimary(ctx *PrimaryContext) {}

// EnterRuleEntity is called when production ruleEntity is entered.
func (s *BasegengineListener) EnterRuleEntity(ctx *RuleEntityContext) {}

// ExitRuleEntity is called when production ruleEntity is exited.
func (s *BasegengineListener) ExitRuleEntity(ctx *RuleEntityContext) {}

// EnterRuleName is called when production ruleName is entered.
func (s *BasegengineListener) EnterRuleName(ctx *RuleNameContext) {}

// ExitRuleName is called when production ruleName is exited.
func (s *BasegengineListener) ExitRuleName(ctx *RuleNameContext) {}

// EnterRuleDescription is called when production ruleDescription is entered.
func (s *BasegengineListener) EnterRuleDescription(ctx *RuleDescriptionContext) {}

// ExitRuleDescription is called when production ruleDescription is exited.
func (s *BasegengineListener) ExitRuleDescription(ctx *RuleDescriptionContext) {}

// EnterSalience is called when production salience is entered.
func (s *BasegengineListener) EnterSalience(ctx *SalienceContext) {}

// ExitSalience is called when production salience is exited.
func (s *BasegengineListener) ExitSalience(ctx *SalienceContext) {}

// EnterRuleContent is called when production ruleContent is entered.
func (s *BasegengineListener) EnterRuleContent(ctx *RuleContentContext) {}

// ExitRuleContent is called when production ruleContent is exited.
func (s *BasegengineListener) ExitRuleContent(ctx *RuleContentContext) {}

// EnterStatements is called when production statements is entered.
func (s *BasegengineListener) EnterStatements(ctx *StatementsContext) {}

// ExitStatements is called when production statements is exited.
func (s *BasegengineListener) ExitStatements(ctx *StatementsContext) {}

// EnterStatement is called when production statement is entered.
func (s *BasegengineListener) EnterStatement(ctx *StatementContext) {}

// ExitStatement is called when production statement is exited.
func (s *BasegengineListener) ExitStatement(ctx *StatementContext) {}

// EnterConcStatement is called when production concStatement is entered.
func (s *BasegengineListener) EnterConcStatement(ctx *ConcStatementContext) {}

// ExitConcStatement is called when production concStatement is exited.
func (s *BasegengineListener) ExitConcStatement(ctx *ConcStatementContext) {}

// EnterExpression is called when production expression is entered.
func (s *BasegengineListener) EnterExpression(ctx *ExpressionContext) {}

// ExitExpression is called when production expression is exited.
func (s *BasegengineListener) ExitExpression(ctx *ExpressionContext) {}

// EnterMathExpression is called when production mathExpression is entered.
func (s *BasegengineListener) EnterMathExpression(ctx *MathExpressionContext) {}

// ExitMathExpression is called when production mathExpression is exited.
func (s *BasegengineListener) ExitMathExpression(ctx *MathExpressionContext) {}

// EnterExpressionAtom is called when production expressionAtom is entered.
func (s *BasegengineListener) EnterExpressionAtom(ctx *ExpressionAtomContext) {}

// ExitExpressionAtom is called when production expressionAtom is exited.
func (s *BasegengineListener) ExitExpressionAtom(ctx *ExpressionAtomContext) {}

// EnterAssignment is called when production assignment is entered.
func (s *BasegengineListener) EnterAssignment(ctx *AssignmentContext) {}

// ExitAssignment is called when production assignment is exited.
func (s *BasegengineListener) ExitAssignment(ctx *AssignmentContext) {}

// EnterReturnStmt is called when production returnStmt is entered.
func (s *BasegengineListener) EnterReturnStmt(ctx *ReturnStmtContext) {}

// ExitReturnStmt is called when production returnStmt is exited.
func (s *BasegengineListener) ExitReturnStmt(ctx *ReturnStmtContext) {}

// EnterIfStmt is called when production ifStmt is entered.
func (s *BasegengineListener) EnterIfStmt(ctx *IfStmtContext) {}

// ExitIfStmt is called when production ifStmt is exited.
func (s *BasegengineListener) ExitIfStmt(ctx *IfStmtContext) {}

// EnterElseIfStmt is called when production elseIfStmt is entered.
func (s *BasegengineListener) EnterElseIfStmt(ctx *ElseIfStmtContext) {}

// ExitElseIfStmt is called when production elseIfStmt is exited.
func (s *BasegengineListener) ExitElseIfStmt(ctx *ElseIfStmtContext) {}

// EnterElseStmt is called when production elseStmt is entered.
func (s *BasegengineListener) EnterElseStmt(ctx *ElseStmtContext) {}

// ExitElseStmt is called when production elseStmt is exited.
func (s *BasegengineListener) ExitElseStmt(ctx *ElseStmtContext) {}

// EnterConstant is called when production constant is entered.
func (s *BasegengineListener) EnterConstant(ctx *ConstantContext) {}

// ExitConstant is called when production constant is exited.
func (s *BasegengineListener) ExitConstant(ctx *ConstantContext) {}

// EnterFunctionArgs is called when production functionArgs is entered.
func (s *BasegengineListener) EnterFunctionArgs(ctx *FunctionArgsContext) {}

// ExitFunctionArgs is called when production functionArgs is exited.
func (s *BasegengineListener) ExitFunctionArgs(ctx *FunctionArgsContext) {}

// EnterInteger is called when production integer is entered.
func (s *BasegengineListener) EnterInteger(ctx *IntegerContext) {}

// ExitInteger is called when production integer is exited.
func (s *BasegengineListener) ExitInteger(ctx *IntegerContext) {}

// EnterRealLiteral is called when production realLiteral is entered.
func (s *BasegengineListener) EnterRealLiteral(ctx *RealLiteralContext) {}

// ExitRealLiteral is called when production realLiteral is exited.
func (s *BasegengineListener) ExitRealLiteral(ctx *RealLiteralContext) {}

// EnterStringLiteral is called when production stringLiteral is entered.
func (s *BasegengineListener) EnterStringLiteral(ctx *StringLiteralContext) {}

// ExitStringLiteral is called when production stringLiteral is exited.
func (s *BasegengineListener) ExitStringLiteral(ctx *StringLiteralContext) {}

// EnterBooleanLiteral is called when production booleanLiteral is entered.
func (s *BasegengineListener) EnterBooleanLiteral(ctx *BooleanLiteralContext) {}

// ExitBooleanLiteral is called when production booleanLiteral is exited.
func (s *BasegengineListener) ExitBooleanLiteral(ctx *BooleanLiteralContext) {}

// EnterFunctionCall is called when production functionCall is entered.
func (s *BasegengineListener) EnterFunctionCall(ctx *FunctionCallContext) {}

// ExitFunctionCall is called when production functionCall is exited.
func (s *BasegengineListener) ExitFunctionCall(ctx *FunctionCallContext) {}

// EnterMethodCall is called when production methodCall is entered.
func (s *BasegengineListener) EnterMethodCall(ctx *MethodCallContext) {}

// ExitMethodCall is called when production methodCall is exited.
func (s *BasegengineListener) ExitMethodCall(ctx *MethodCallContext) {}

// EnterVariable is called when production variable is entered.
func (s *BasegengineListener) EnterVariable(ctx *VariableContext) {}

// ExitVariable is called when production variable is exited.
func (s *BasegengineListener) ExitVariable(ctx *VariableContext) {}

// EnterMathPmOperator is called when production mathPmOperator is entered.
func (s *BasegengineListener) EnterMathPmOperator(ctx *MathPmOperatorContext) {}

// ExitMathPmOperator is called when production mathPmOperator is exited.
func (s *BasegengineListener) ExitMathPmOperator(ctx *MathPmOperatorContext) {}

// EnterMathMdOperator is called when production mathMdOperator is entered.
func (s *BasegengineListener) EnterMathMdOperator(ctx *MathMdOperatorContext) {}

// ExitMathMdOperator is called when production mathMdOperator is exited.
func (s *BasegengineListener) ExitMathMdOperator(ctx *MathMdOperatorContext) {}

// EnterComparisonOperator is called when production comparisonOperator is entered.
func (s *BasegengineListener) EnterComparisonOperator(ctx *ComparisonOperatorContext) {}

// ExitComparisonOperator is called when production comparisonOperator is exited.
func (s *BasegengineListener) ExitComparisonOperator(ctx *ComparisonOperatorContext) {}

// EnterLogicalOperator is called when production logicalOperator is entered.
func (s *BasegengineListener) EnterLogicalOperator(ctx *LogicalOperatorContext) {}

// ExitLogicalOperator is called when production logicalOperator is exited.
func (s *BasegengineListener) ExitLogicalOperator(ctx *LogicalOperatorContext) {}

// EnterAssignOperator is called when production assignOperator is entered.
func (s *BasegengineListener) EnterAssignOperator(ctx *AssignOperatorContext) {}

// ExitAssignOperator is called when production assignOperator is exited.
func (s *BasegengineListener) ExitAssignOperator(ctx *AssignOperatorContext) {}

// EnterNotOperator is called when production notOperator is entered.
func (s *BasegengineListener) EnterNotOperator(ctx *NotOperatorContext) {}

// ExitNotOperator is called when production notOperator is exited.
func (s *BasegengineListener) ExitNotOperator(ctx *NotOperatorContext) {}

// EnterMapVar is called when production mapVar is entered.
func (s *BasegengineListener) EnterMapVar(ctx *MapVarContext) {}

// ExitMapVar is called when production mapVar is exited.
func (s *BasegengineListener) ExitMapVar(ctx *MapVarContext) {}

// EnterAtName is called when production atName is entered.
func (s *BasegengineListener) EnterAtName(ctx *AtNameContext) {}

// ExitAtName is called when production atName is exited.
func (s *BasegengineListener) ExitAtName(ctx *AtNameContext) {}

// EnterAtId is called when production atId is entered.
func (s *BasegengineListener) EnterAtId(ctx *AtIdContext) {}

// ExitAtId is called when production atId is exited.
func (s *BasegengineListener) ExitAtId(ctx *AtIdContext) {}

// EnterAtDesc is called when production atDesc is entered.
func (s *BasegengineListener) EnterAtDesc(ctx *AtDescContext) {}

// ExitAtDesc is called when production atDesc is exited.
func (s *BasegengineListener) ExitAtDesc(ctx *AtDescContext) {}

// EnterAtSal is called when production atSal is entered.
func (s *BasegengineListener) EnterAtSal(ctx *AtSalContext) {}

// ExitAtSal is called when production atSal is exited.
func (s *BasegengineListener) ExitAtSal(ctx *AtSalContext) {}
