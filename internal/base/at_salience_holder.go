package base

type AtSalienceHolder interface {
	AcceptSalience(val int64) error
}
