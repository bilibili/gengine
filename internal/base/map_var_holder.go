package base

type MapVarHolder interface {
	AcceptMapVar(mv *MapVar) error
}
