package base

type ExpressionAtomHolder interface {
	AcceptExpressionAtom(atom *ExpressionAtom) error
}
