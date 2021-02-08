package base

import (
	"github.com/bilibili/gengine/context"
	"reflect"
)

type Statements struct {
	StatementList   []*Statement
	ReturnStatement *ReturnStatement
}

func (s *Statements) Evaluate(dc *context.DataContext, Vars map[string]reflect.Value) (reflect.Value, error, bool) {
	for _, statement := range s.StatementList {
		v, err, b := statement.Evaluate(dc, Vars)
		if err != nil {
			return reflect.ValueOf(nil), err, false
		}

		if b {
			//important: meet return，not continue to execute
			return v, nil, b
		}
	}
	if s.ReturnStatement != nil {
		return s.ReturnStatement.Evaluate(dc, Vars)
	}
	return reflect.ValueOf(nil), nil, false
}
