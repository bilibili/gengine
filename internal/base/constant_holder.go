package base

type ConstantHolder interface {
	AcceptConstant(cons *Constant) error
}
