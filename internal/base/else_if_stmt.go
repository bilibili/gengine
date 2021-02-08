package base

import (
	"errors"
	"github.com/bilibili/gengine/context"
	"reflect"
)

type ElseIfStmt struct {
	Expression    *Expression
	StatementList *Statements
}

func (ef *ElseIfStmt) Evaluate(dc *context.DataContext, Vars map[string]reflect.Value) (reflect.Value, error, bool) {
	it, err := ef.Expression.Evaluate(dc, Vars)
	if err != nil {
		return reflect.ValueOf(nil), err, false
	}

	if it.Bool() {
		if ef.StatementList == nil {
			return reflect.ValueOf(nil), nil, false
		} else {
			return ef.StatementList.Evaluate(dc, Vars)
		}
	} else {
		return reflect.ValueOf(nil), nil, false
	}
}

func (ef *ElseIfStmt) AcceptExpression(expr *Expression) error {

	if ef.Expression == nil {
		ef.Expression = expr
		return nil
	}
	return errors.New("ElseIfStmt's Expression set twice! ")
}

func (ef *ElseIfStmt) AcceptStatements(stmts *Statements) error {
	if ef.StatementList == nil {
		ef.StatementList = stmts
		return nil
	}
	return errors.New("ElseIfStmt's statements set twice! ")
}
