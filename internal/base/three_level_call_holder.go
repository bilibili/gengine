package base

type ThreeLevelCallHolder interface {
	AcceptThreeLevelCall(threeLevelCall *ThreeLevelCall) error
}
