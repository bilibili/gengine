package base

import (
	"errors"
	"fmt"
	"github.com/bilibili/gengine/context"
	"reflect"
	"runtime"
	"strings"
)

type ThreeLevelCall struct {
	SourceCode
	ThreeLevel string
	MethodArgs *Args
}

func (tlc *ThreeLevelCall) AcceptArgs(funcArg *Args) error {
	if tlc.MethodArgs == nil {
		tlc.MethodArgs = funcArg
		return nil
	}
	return errors.New("ThreeLevelCall set twice! ")
}

func (tlc *ThreeLevelCall) Evaluate(dc *context.DataContext, Vars map[string]reflect.Value) (mr reflect.Value, err error) {

	defer func() {
		if e := recover(); e != nil {
			size := 1 << 10 * 10
			buf := make([]byte, size)
			rs := runtime.Stack(buf, false)
			if rs > size {
				rs = size
			}
			buf = buf[:rs]
			eMsg := fmt.Sprintf("line %d, column %d, code: %s, %+v \n%s", tlc.LineNum, tlc.Column, tlc.Code, e, string(buf))
			eMsg = strings.ReplaceAll(eMsg, "panic", "error")
			err = errors.New(eMsg)
		}
	}()

	var argumentValues []reflect.Value
	if tlc.MethodArgs == nil {
		argumentValues = make([]reflect.Value, 0)
	} else {
		av, err := tlc.MethodArgs.Evaluate(dc, Vars)
		if err != nil {
			return reflect.ValueOf(nil), err
		}
		argumentValues = av
	}

	mr, err = dc.ExecThreeLevel(Vars, tlc.ThreeLevel, argumentValues)
	if err != nil {
		return reflect.ValueOf(nil), errors.New(fmt.Sprintf("line %d, column %d, code: %s, %+v", tlc.LineNum, tlc.Column, tlc.Code, err))
	}
	return
}
