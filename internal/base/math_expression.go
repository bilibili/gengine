package base

import (
	"errors"
	"fmt"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/internal/core"
	"reflect"
)

type MathExpression struct {
	SourceCode
	MathExpressionLeft  *MathExpression
	MathPmOperator      string
	MathMdOperator      string
	MathExpressionRight *MathExpression
	ExpressionAtom      *ExpressionAtom
}

func (e *MathExpression) AcceptMathExpression(atom *MathExpression) error {
	if e.MathExpressionLeft == nil {
		e.MathExpressionLeft = atom
		return nil
	}
	if e.MathExpressionRight == nil {
		e.MathExpressionRight = atom
		return nil
	}
	return errors.New("expressionAtom set twice")
}

func (e *MathExpression) AcceptExpressionAtom(atom *ExpressionAtom) error {
	if e.ExpressionAtom == nil {
		e.ExpressionAtom = atom
		return nil
	}
	return errors.New("ExpressionAtom already set twice ")
}

func (e *MathExpression) Evaluate(dc *context.DataContext, Vars map[string]reflect.Value) (reflect.Value, error) {

	//priority to calculate single value
	if e.ExpressionAtom != nil {
		return e.ExpressionAtom.Evaluate(dc, Vars)
	}

	// check the right whether is nil
	if e.MathExpressionRight == nil {
		return e.MathExpressionLeft.Evaluate(dc, Vars)
	}

	lv, err := e.MathExpressionLeft.Evaluate(dc, Vars)
	if err != nil {
		return reflect.ValueOf(nil), err
	}
	rv, err := e.MathExpressionRight.Evaluate(dc, Vars)
	if err != nil {
		return reflect.ValueOf(nil), err
	}

	if e.MathPmOperator == "+" {
		add, err := core.Add(lv, rv)
		if err != nil {
			return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code: %s, %+v", e.LineNum, e.Column, e.Code, err))
		}
		return reflect.ValueOf(add), nil
	}

	if e.MathPmOperator == "-" {
		sub, err := core.Sub(lv, rv)
		if err != nil {
			return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code: %s, %+v", e.LineNum, e.Column, e.Code, err))
		}
		return reflect.ValueOf(sub), nil
	}

	if e.MathMdOperator == "*" {
		mul, err := core.Mul(lv, rv)
		if err != nil {
			return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code: %s, %+v", e.LineNum, e.Column, e.Code, err))
		}
		return reflect.ValueOf(mul), nil
	}

	if e.MathMdOperator == "/" {
		div, err := core.Div(lv, rv)
		if err != nil {
			return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code: %s, %+v", e.LineNum, e.Column, e.Code, err))
		}
		return reflect.ValueOf(div), nil
	}
	return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code: %s, MathExpression calculate evaluate error", e.LineNum, e.Column, e.Code))
}
