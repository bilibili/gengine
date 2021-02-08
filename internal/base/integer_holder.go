package base

type IntegerHolder interface {
	AcceptInteger(val int64) error
}
