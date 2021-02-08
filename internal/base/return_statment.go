package base

import (
	"errors"
	"github.com/bilibili/gengine/context"
	"reflect"
)

type ReturnStatement struct {
	Expression *Expression
}

func (rs *ReturnStatement) Evaluate(dc *context.DataContext, Vars map[string]reflect.Value) (reflect.Value, error, bool) {
	if rs.Expression != nil {
		value, e := rs.Expression.Evaluate(dc, Vars)
		return value, e, true
	}
	return reflect.ValueOf(nil), nil, true
}

func (rs *ReturnStatement) AcceptExpression(expr *Expression) error {
	if rs.Expression == nil {
		rs.Expression = expr
		return nil
	}
	return errors.New("Expression already set twice! ")
}
