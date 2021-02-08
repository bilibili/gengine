package base

type ArgsHolder interface {
	AcceptArgs(funcArg *Args) error
}
