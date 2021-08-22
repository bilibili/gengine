package base

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/bilibili/gengine/context"
)

const maxExecuteNum = 10000

// 	for循环语法 for i=0;i<100;i+=1
type ForStmt struct {
	SourceCode
	Expression    *Expression   // 表达式
	StatementList *Statements   // 执行语句列表
	Assignments   []*Assignment // 赋值语句
}

func (forStmt *ForStmt) AcceptAssignment(assignment *Assignment) error {
	forStmt.Assignments = append(forStmt.Assignments, assignment)
	return nil
}

func (forStmt *ForStmt) AcceptExpression(expr *Expression) error {

	if forStmt.Expression == nil {
		forStmt.Expression = expr
		return nil
	}
	return errors.New("ForStmt's Expression set twice! ")
}

func (forStmt *ForStmt) Evaluate(dc *context.DataContext, Vars map[string]reflect.Value) (reflect.Value, error, bool) {
	if len(forStmt.Assignments) < 2 {
		return reflect.ValueOf(nil), errors.New("assignments len failed"), false
	}
	// 首先进行语句的赋值
	_, err := forStmt.Assignments[0].Evaluate(dc, Vars)
	if err != nil {
		return reflect.ValueOf(nil), err, false
	}

	// 设置最大的循环次数，防止因为死循环导致CPU利用率高。当前编排的场景不会有这么多的循环
	iCount := 0
	for {
		iCount++
		if iCount > maxExecuteNum {
			return reflect.ValueOf(nil),
				fmt.Errorf("execute for bigger than maxExecuteNum:%v", maxExecuteNum), false
		}
		it, err := forStmt.Expression.Evaluate(dc, Vars)
		if err != nil {
			return reflect.ValueOf(nil), err, false
		}
		if it.Bool() {
			if forStmt.StatementList != nil {
				vStatement, errStatement, bStatement := forStmt.StatementList.Evaluate(dc, Vars)
				if errStatement != nil && errStatement != BREAKFLAG && errStatement != CONTINUEFLAG {
					return reflect.ValueOf(nil), errStatement, false
				}
				// 如果是continue就继续循环
				if errStatement == CONTINUEFLAG {
					_, err = forStmt.Assignments[1].Evaluate(dc, Vars)
					if err != nil {
						return reflect.ValueOf(nil), err, false
					}
					continue
				}

				// 如果是break就跳出循环
				if errStatement == BREAKFLAG {
					break
				}
				// 如果是return返回，直接终止流程
				if bStatement {
					return vStatement, err, bStatement
				}
			} else {
				break
			}

			_, err = forStmt.Assignments[1].Evaluate(dc, Vars)
			if err != nil {
				return reflect.ValueOf(nil), err, false
			}
		} else {
			break
		}
	}
	return reflect.ValueOf(nil), nil, false
}

func (forStmt *ForStmt) AcceptStatements(stmts *Statements) error {
	if forStmt.StatementList == nil {
		forStmt.StatementList = stmts
		return nil
	}
	return errors.New("forStmt set twice! ")
}
