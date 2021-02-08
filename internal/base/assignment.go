package base

import (
	"errors"
	"fmt"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/internal/core"
	"reflect"
	"runtime"
	"strings"
)

// := or =
type Assignment struct {
	SourceCode
	Variable       string
	MapVar         *MapVar
	AssignOperator string
	MathExpression *MathExpression
	Expression     *Expression
}

func (a *Assignment) Evaluate(dc *context.DataContext, Vars map[string]reflect.Value) (value reflect.Value, err error) {

	defer func() {
		if e := recover(); e != nil {
			size := 1 << 10 * 10
			buf := make([]byte, size)
			rs := runtime.Stack(buf, false)
			if rs > size {
				rs = size
			}
			buf = buf[:rs]
			eMsg := fmt.Sprintf("line %d, column %d, code: %s, %+v \n%s", a.LineNum, a.Column, a.Code, e, string(buf))
			eMsg = strings.ReplaceAll(eMsg, "panic", "error")
			err = errors.New(eMsg)
		}
	}()

	var mv reflect.Value

	if a.MathExpression != nil {
		mv, err = a.MathExpression.Evaluate(dc, Vars)
		if err != nil {
			return reflect.ValueOf(nil), err
		}
	}

	if a.Expression != nil {
		mv, err = a.Expression.Evaluate(dc, Vars)
		if err != nil {
			return reflect.ValueOf(nil), err
		}
	}

	var sv reflect.Value

	if a.AssignOperator == "=" || a.AssignOperator == ":=" {
		goto END
	}

	if len(a.Variable) > 0 {
		sv, err = dc.GetValue(Vars, a.Variable)
		if err != nil {
			return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column:%d, code: %s, %+v:", a.LineNum, a.Column, a.Code, err))
		}
	}

	if a.MapVar != nil {
		sv, err = a.MapVar.Evaluate(dc, Vars)
		if err != nil {
			return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column:%d, code: %s, %+v:", a.LineNum, a.Column, a.Code, err))
		}
	}

	if a.AssignOperator == "+=" {
		_mv, err := core.Add(sv, mv)
		if err != nil {
			return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column:%d, code: %s, %+v:", a.LineNum, a.Column, a.Code, err))
		}
		mv = reflect.ValueOf(_mv)
		goto END
	}

	if a.AssignOperator == "-=" {
		_mv, err := core.Sub(sv, mv)
		if err != nil {
			return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column:%d, code: %s, %+v:", a.LineNum, a.Column, a.Code, err))
		}
		mv = reflect.ValueOf(_mv)
		goto END
	}

	if a.AssignOperator == "*=" {
		_mv, err := core.Mul(sv, mv)
		if err != nil {
			return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column:%d, code: %s, %+v:", a.LineNum, a.Column, a.Code, err))
		}
		mv = reflect.ValueOf(_mv)
		goto END
	}

	if a.AssignOperator == "/=" {
		_mv, err := core.Div(sv, mv)
		if err != nil {
			return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column:%d, code: %s, %+v:", a.LineNum, a.Column, a.Code, err))
		}
		mv = reflect.ValueOf(_mv)
		goto END
	}

END:
	if len(a.Variable) > 0 {
		err = dc.SetValue(Vars, a.Variable, mv)
		if err != nil {
			return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code: %s, %+v", a.LineNum, a.Column, a.Code, err))
		}
		return
	}

	if a.MapVar != nil {
		err = dc.SetMapVarValue(Vars, a.MapVar.Name, a.MapVar.Strkey, a.MapVar.Varkey, a.MapVar.Intkey, mv)
		if err != nil {
			return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column:%d, code: %s, %+v:", a.LineNum, a.Column, a.Code, err))
		}
		return
	}

	return
}

func (a *Assignment) AcceptMathExpression(me *MathExpression) error {
	if a.MathExpression == nil {
		a.MathExpression = me
		return nil
	}
	return errors.New("MathExpression already set twice! ")
}

func (a *Assignment) AcceptVariable(name string) error {
	if len(a.Variable) == 0 {
		a.Variable = name
		return nil
	}
	return errors.New("Variable already set twice! ")
}

func (a *Assignment) AcceptMapVar(mapVar *MapVar) error {
	if a.MapVar == nil {
		a.MapVar = mapVar
		return nil
	}
	return errors.New("MapVar already set twice")
}

func (a *Assignment) AcceptExpression(exp *Expression) error {
	if a.Expression == nil {
		a.Expression = exp
		return nil
	}
	return errors.New("Expression already set twice! ")
}
