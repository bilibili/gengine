package base

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/internal/iter"
)

// 	forRane循环语法 forRange key=mapInfo{}
type ForRangeStmt struct {
	SourceCode
	StatementList *Statements
	keyName       string // key的名字
	name          string // 切片或map变量的名字
}

func (forRangeStmt *ForRangeStmt) Evaluate(dc *context.DataContext, Vars map[string]reflect.Value) (reflect.Value, error, bool) {
	var err error
	value, e := dc.GetValue(Vars, forRangeStmt.name)
	if e != nil {
		return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code: %s, %+v",
			forRangeStmt.LineNum, forRangeStmt.Column, forRangeStmt.Code, e)), false
	}

	// 判断是否可以迭代
	iterer, err := iter.NewInter(value)
	if err != nil {
		return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code: %s, %+v",
			forRangeStmt.LineNum, forRangeStmt.Column, forRangeStmt.Code, e)), false
	}
	for iterer.Next() {
		key := iterer.Key()
		err = dc.SetValue(Vars, forRangeStmt.keyName, key)
		if err != nil {
			return reflect.ValueOf(nil), err, false
		}
		if forRangeStmt.StatementList != nil {
			vStatement, errStatement, bStatement := forRangeStmt.StatementList.Evaluate(dc, Vars)
			if errStatement != nil && errStatement != BREAKFLAG && errStatement != CONTINUEFLAG {
				return reflect.ValueOf(nil), errStatement, false
			}
			// 如果是continue就继续循环
			if errStatement == CONTINUEFLAG {
				continue
			}

			// 如果是break就跳出循环
			if errStatement == BREAKFLAG {
				break
			}
			// 如果是return返回，直接终止流程
			if bStatement {
				return vStatement, errStatement, bStatement
			}
		} else {
			break
		}
	}
	return reflect.ValueOf(nil), nil, false
}

func (forRangeStmt *ForRangeStmt) AcceptVariable(name string) error {
	if len(forRangeStmt.keyName) == 0 {
		forRangeStmt.keyName = name
		return nil
	}

	if len(forRangeStmt.name) == 0 {
		forRangeStmt.name = name
		return nil
	}
	return errors.New("forRang Varkey set three times! ")
}

func (forRangeStmt *ForRangeStmt) AcceptStatements(stmts *Statements) error {
	if forRangeStmt.StatementList == nil {
		forRangeStmt.StatementList = stmts
		return nil
	}
	return errors.New("forRang set twice! ")
}
