package base

import (
	"errors"
	"fmt"
	"github.com/bilibili/gengine/context"
	"reflect"
)

type RuleEntity struct {
	RuleName        string
	Salience        int64
	RuleDescription string
	RuleContent     *RuleContent
	Vars            map[string]reflect.Value //belongs to current rule,rule execute finish, it will be clear
}

func (r *RuleEntity) AcceptString(s string) error {
	if r.RuleName == "" {
		r.RuleName = s
		return nil
	}

	if r.RuleDescription == "" {
		r.RuleDescription = s
		return nil
	}
	return errors.New(fmt.Sprintf("value = %s set twice!", s))
}

func (r *RuleEntity) AcceptInteger(val int64) error {
	r.Salience = val
	return nil
}

/*func (r *RuleEntity) Initialize(dc *context.DataContext) {
	r.dataCtx = dc

	if r.RuleContent != nil {
		r.RuleContent.Initialize(dc)
	}
}
*/
func (r *RuleEntity) Execute(dc *context.DataContext) (interface{}, error, bool) {
	r.Vars = make(map[string]reflect.Value)
	defer r.clearMap()
	v, e, b := r.RuleContent.Execute(dc, r.Vars)
	if v == reflect.ValueOf(nil) {
		return nil, e, b
	}
	return v.Interface(), e, b
}

func (r *RuleEntity) clearMap() {
	r.Vars = make(map[string]reflect.Value)
}
