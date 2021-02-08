package base

type MathExpressionHolder interface {
	AcceptMathExpression(mh *MathExpression) error
}
