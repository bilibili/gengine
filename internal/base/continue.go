package base

import (
	"errors"
	"reflect"

	"github.com/bilibili/gengine/context"
)

type ContinueStmt struct {
}

var CONTINUEFLAG = errors.New("break")

func (continueStmt *ContinueStmt) Evaluate(dc *context.DataContext, Vars map[string]reflect.Value) (reflect.Value, error, bool) {
	// 直接返回错误给上一层，上一层收到错误信息后判断err,继续循环
	return reflect.ValueOf(nil), CONTINUEFLAG, true
}
