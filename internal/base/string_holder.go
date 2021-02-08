package base

type StringHolder interface {
	AcceptString(str string) error
}
