package base

type MethodCallHolder interface {
	AcceptMethodCall(methodCall *MethodCall) error
}
