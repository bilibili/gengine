package base

type ExpressionHolder interface {
	AcceptExpression(expression *Expression) error
}
