package base

import (
	"errors"
	"github.com/bilibili/gengine/context"
	"reflect"
)

type ExpressionAtom struct {
	SourceCode
	Variable       string
	Constant       *Constant
	FunctionCall   *FunctionCall
	MethodCall     *MethodCall
	ThreeLevelCall *ThreeLevelCall
	MapVar         *MapVar
}

func (e *ExpressionAtom) AcceptVariable(name string) error {
	if len(e.Variable) == 0 {
		e.Variable = name
		return nil
	}
	return errors.New("Variable already defined! ")
}

func (e *ExpressionAtom) AcceptConstant(cons *Constant) error {
	if e.Constant == nil {
		e.Constant = cons
		return nil
	}
	return errors.New("Constant already defined! ")
}

func (e *ExpressionAtom) AcceptFunctionCall(funcCall *FunctionCall) error {
	if e.FunctionCall == nil {
		e.FunctionCall = funcCall
		return nil
	}
	return errors.New("FunctionCall already defined! ")
}

func (e *ExpressionAtom) AcceptMethodCall(methodCall *MethodCall) error {
	if e.MethodCall == nil {
		e.MethodCall = methodCall
		return nil
	}
	return errors.New("MethodCall already defined! ")
}

func (e *ExpressionAtom) AcceptThreeLevelCall(threeLevelCall *ThreeLevelCall) error {
	if e.ThreeLevelCall == nil {
		e.ThreeLevelCall = threeLevelCall
		return nil
	}
	return errors.New("threeLevelCall already defined! ")
}

func (e *ExpressionAtom) AcceptMapVar(mapVar *MapVar) error {
	if e.MapVar == nil {
		e.MapVar = mapVar
		return nil
	}
	return errors.New("MapVar already defined! ")
}

func (e *ExpressionAtom) Evaluate(dc *context.DataContext, Vars map[string]reflect.Value) (reflect.Value, error) {
	if len(e.Variable) > 0 {
		return dc.GetValue(Vars, e.Variable)
	} else if e.Constant != nil {
		return e.Constant.Evaluate(dc, Vars)
	} else if e.FunctionCall != nil {
		return e.FunctionCall.Evaluate(dc, Vars)
	} else if e.MethodCall != nil {
		return e.MethodCall.Evaluate(dc, Vars)
	} else if e.MapVar != nil {
		return e.MapVar.Evaluate(dc, Vars)
	} else if e.ThreeLevelCall != nil {
		return e.ThreeLevelCall.Evaluate(dc, Vars)
	}
	//todo
	return reflect.ValueOf(nil), errors.New("ExpressionAtom Evaluate error! ")
}
