package base

type VariableHolder interface {
	AcceptVariable(name string) error
}
