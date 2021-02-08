// Code generated from /Users/renyunyi/go_project/gengine/internal/iantlr/gengine.g4 by ANTLR 4.9. DO NOT EDIT.

package parser // gengine

import "github.com/antlr/antlr4/runtime/Go/antlr"

type BasegengineVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BasegengineVisitor) VisitPrimary(ctx *PrimaryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitRuleEntity(ctx *RuleEntityContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitRuleName(ctx *RuleNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitRuleDescription(ctx *RuleDescriptionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitSalience(ctx *SalienceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitRuleContent(ctx *RuleContentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitStatements(ctx *StatementsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitStatement(ctx *StatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitConcStatement(ctx *ConcStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitExpression(ctx *ExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitMathExpression(ctx *MathExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitExpressionAtom(ctx *ExpressionAtomContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitAssignment(ctx *AssignmentContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitReturnStmt(ctx *ReturnStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitIfStmt(ctx *IfStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitElseIfStmt(ctx *ElseIfStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitElseStmt(ctx *ElseStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitConstant(ctx *ConstantContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitFunctionArgs(ctx *FunctionArgsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitInteger(ctx *IntegerContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitRealLiteral(ctx *RealLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitStringLiteral(ctx *StringLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitBooleanLiteral(ctx *BooleanLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitFunctionCall(ctx *FunctionCallContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitMethodCall(ctx *MethodCallContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitVariable(ctx *VariableContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitMathPmOperator(ctx *MathPmOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitMathMdOperator(ctx *MathMdOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitComparisonOperator(ctx *ComparisonOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitLogicalOperator(ctx *LogicalOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitAssignOperator(ctx *AssignOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitNotOperator(ctx *NotOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitMapVar(ctx *MapVarContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitAtName(ctx *AtNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitAtId(ctx *AtIdContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitAtDesc(ctx *AtDescContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BasegengineVisitor) VisitAtSal(ctx *AtSalContext) interface{} {
	return v.VisitChildren(ctx)
}
