package base

import (
	"errors"
	"github.com/bilibili/gengine/context"
	"reflect"
)

type Arg struct {
	Constant       *Constant
	Variable       string
	FunctionCall   *FunctionCall
	MethodCall     *MethodCall
	ThreeLevelCall *ThreeLevelCall
	MapVar         *MapVar
	Expression     *Expression
}

func (a *Arg) Evaluate(dc *context.DataContext, Vars map[string]reflect.Value) (reflect.Value, error) {
	if len(a.Variable) > 0 {
		return dc.GetValue(Vars, a.Variable)
	}

	if a.Constant != nil {
		return a.Constant.Evaluate(dc, Vars)
	}

	if a.FunctionCall != nil {
		return a.FunctionCall.Evaluate(dc, Vars)
	}

	if a.MethodCall != nil {
		return a.MethodCall.Evaluate(dc, Vars)
	}

	if a.ThreeLevelCall != nil {
		return a.ThreeLevelCall.Evaluate(dc, Vars)
	}

	if a.MapVar != nil {
		return a.MapVar.Evaluate(dc, Vars)
	}

	if a.Expression != nil {
		return a.Expression.Evaluate(dc, Vars)
	}

	return reflect.ValueOf(nil), errors.New("argHolder holder has more values than want！")
}
