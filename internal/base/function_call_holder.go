package base

type FunctionCallHolder interface {
	AcceptFunctionCall(funcCall *FunctionCall) error
}
