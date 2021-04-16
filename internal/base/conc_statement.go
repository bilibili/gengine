package base

import (
	"errors"
	"fmt"
	"github.com/bilibili/gengine/context"
	"reflect"
	"sync"
)

type ConcStatement struct {
	Assignments     []*Assignment
	FunctionCalls   []*FunctionCall
	MethodCalls     []*MethodCall
	ThreeLevelCalls []*ThreeLevelCall
}

func (cs *ConcStatement) AcceptAssignment(assignment *Assignment) error {
	cs.Assignments = append(cs.Assignments, assignment)
	return nil
}

func (cs *ConcStatement) AcceptFunctionCall(funcCall *FunctionCall) error {
	cs.FunctionCalls = append(cs.FunctionCalls, funcCall)
	return nil
}

func (cs *ConcStatement) AcceptMethodCall(methodCall *MethodCall) error {
	cs.MethodCalls = append(cs.MethodCalls, methodCall)
	return nil
}

func (cs *ConcStatement) AcceptThreeLevelCall(threeLevelCall *ThreeLevelCall) error {
	cs.ThreeLevelCalls = append(cs.ThreeLevelCalls, threeLevelCall)
	return nil
}

func (cs *ConcStatement) Evaluate(dc *context.DataContext, Vars map[string]reflect.Value) (reflect.Value, error) {

	aLen := len(cs.Assignments)
	fLen := len(cs.FunctionCalls)
	mLen := len(cs.MethodCalls)
	tLen := len(cs.ThreeLevelCalls)
	l := aLen + fLen + mLen + tLen
	if l <= 0 {
		return reflect.ValueOf(nil), nil
	} else {

		var errLock sync.Mutex
		var eMsg []string

		var wg sync.WaitGroup
		wg.Add(l)
		go func() {
			for _, assign := range cs.Assignments {
				assignment := assign
				go func() {
					_, e := assignment.Evaluate(dc, Vars)
					if e != nil {
						errLock.Lock()
						eMsg = append(eMsg, fmt.Sprintf("%+v", e))
						errLock.Unlock()
					}
					wg.Done()
				}()
			}
		}()

		go func() {
			for _, fu := range cs.FunctionCalls {
				fun := fu
				go func() {
					_, e := fun.Evaluate(dc, Vars)
					if e != nil {
						errLock.Lock()
						eMsg = append(eMsg, fmt.Sprintf("%+v", e))
						errLock.Unlock()
					}
					wg.Done()
				}()
			}
		}()
		go func() {
			for _, me := range cs.MethodCalls {
				meth := me
				go func() {
					_, e := meth.Evaluate(dc, Vars)
					if e != nil {
						errLock.Lock()
						eMsg = append(eMsg, fmt.Sprintf("%+v", e))
						errLock.Unlock()
					}
					wg.Done()
				}()
			}
		}()

		go func() {
			for _, c := range cs.ThreeLevelCalls {
				tlc := c
				go func() {
					_, e := tlc.Evaluate(dc, Vars)
					if e != nil {
						errLock.Lock()
						eMsg = append(eMsg, fmt.Sprintf("%+v", e))
						errLock.Unlock()
					}
					wg.Done()
				}()
			}
		}()

		wg.Wait()

		if len(eMsg) > 0 {
			return reflect.ValueOf(nil), errors.New(fmt.Sprintf("%+v", eMsg))
		}
	}
	return reflect.ValueOf(nil), nil
}
