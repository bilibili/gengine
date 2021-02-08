package base

import (
	"errors"
	"github.com/bilibili/gengine/context"
	"reflect"
)

type IfStmt struct {
	Expression     *Expression
	StatementList  *Statements
	ElseIfStmtList []*ElseIfStmt
	ElseStmt       *ElseStmt
}

func (i *IfStmt) Evaluate(dc *context.DataContext, Vars map[string]reflect.Value) (reflect.Value, error, bool) {

	it, err := i.Expression.Evaluate(dc, Vars)
	if err != nil {
		return reflect.ValueOf(nil), err, false
	}

	if it.Bool() {
		if i.StatementList == nil {
			return reflect.ValueOf(nil), nil, false
		} else {
			return i.StatementList.Evaluate(dc, Vars)
		}

	} else {

		if i.ElseIfStmtList != nil {
			for _, elseIfStmt := range i.ElseIfStmtList {
				v, err := elseIfStmt.Expression.Evaluate(dc, Vars)
				if err != nil {
					return reflect.ValueOf(nil), err, false
				}

				if v.Bool() {
					return elseIfStmt.StatementList.Evaluate(dc, Vars)
				}
			}
		}

		if i.ElseStmt != nil {
			return i.ElseStmt.Evaluate(dc, Vars)
		} else {
			return reflect.ValueOf(nil), nil, false
		}
	}
}

func (i *IfStmt) AcceptExpression(expr *Expression) error {
	if i.Expression == nil {
		i.Expression = expr
		return nil
	}
	return errors.New("IfStmt Expression set twice ")
}

func (i *IfStmt) AcceptStatements(stmts *Statements) error {
	if i.StatementList == nil {
		i.StatementList = stmts
		return nil
	}
	return errors.New("ifStmt's statements set twice ")
}
