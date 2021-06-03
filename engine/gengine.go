package engine

import (
	"errors"
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/internal/base"
	"sort"
	"sync"

	"github.com/google/martian/log"
)

type Gengine struct {
	lock         sync.Mutex
	returnResult map[string]interface{}
}

func NewGengine() *Gengine {
	return &Gengine{}
}

type Stag struct {
	StopTag bool
}

func (g *Gengine) addResult(name string, returnResult interface{}) {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.returnResult[name] = returnResult
}

func (g *Gengine) GetRulesResultMap() (map[string]interface{}, error) {
	return g.returnResult, nil
}

/**
sort execute model

when b is true it means when there are many rules， if one rule execute error，continue to execute rules after the occur error rule
*/
func (g *Gengine) Execute(rb *builder.RuleBuilder, b bool) error {

	//check rb
	if rb == nil {
		return errors.New("ruleBuilder is nil")
	}

	g.returnResult = make(map[string]interface{})

	if len(rb.Kc.SortRules) == 0 {
		return errors.New("no rule has been injected into engine! ")
	}

	var eMsg []string
	for _, r := range rb.Kc.SortRules {
		v, err, bx := r.Execute(rb.Dc)
		if bx {
			g.addResult(r.RuleName, v)
		}

		if err != nil {
			if b {
				eMsg = append(eMsg, fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", r.RuleName, err))
			} else {
				return errors.New(fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", r.RuleName, err))
			}
		}
	}

	if len(eMsg) > 0 {
		return errors.New(fmt.Sprintf("%+v", eMsg))
	}
	return nil
}

/**
sort execute model

when b is true it means when there are many rules， if one rule execute error，continue to execute rules after the occur error rule;
if stopTag become true,it will not continue to execute

sTag is a struct given by user, and user can use it  to control rules execute behavior in rules, it can improve performance

it used in this scene:
where some high priority rules execute finished, you don't want to execute to the last rules, you can use sTag to control it out of gengine
*/
func (g *Gengine) ExecuteWithStopTagDirect(rb *builder.RuleBuilder, b bool, sTag *Stag) error {

	//check rb
	if rb == nil {
		return errors.New("ruleBuilder is nil")
	}

	g.returnResult = make(map[string]interface{})

	if len(rb.Kc.SortRules) == 0 {
		return errors.New("no rule has been injected into engine! ")
	}

	var eMsg []string
	for _, r := range rb.Kc.SortRules {
		v, err, bx := r.Execute(rb.Dc)
		if bx {
			g.addResult(r.RuleName, v)
		}
		if err != nil {
			if b {
				eMsg = append(eMsg, fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", r.RuleName, err))
			} else {
				return errors.New(fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", r.RuleName, err))
			}
		}

		if sTag.StopTag {
			break
		}
	}

	if len(eMsg) > 0 {
		return errors.New(fmt.Sprintf("%+v", eMsg))
	}

	return nil
}

/*
 concurrent execute model
 in this mode, it will not consider the priority  and not consider err control
*/
func (g *Gengine) ExecuteConcurrent(rb *builder.RuleBuilder) error {

	//check rb
	if rb == nil {
		return errors.New("ruleBuilder is nil")
	}

	g.returnResult = make(map[string]interface{})

	if len(rb.Kc.RuleEntities) == 0 {
		return errors.New("no rule has been injected into engine! ")
	}

	var errLock sync.Mutex
	var eMsg []string

	var wg sync.WaitGroup
	wg.Add(len(rb.Kc.RuleEntities))
	for _, r := range rb.Kc.RuleEntities {
		rr := r
		go func() {
			v, e, bx := rr.Execute(rb.Dc)
			if bx {
				g.addResult(rr.RuleName, v)
			}
			if e != nil {
				errLock.Lock()
				eMsg = append(eMsg, fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", rr.RuleName, e))
				errLock.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if len(eMsg) > 0 {
		return errors.New(fmt.Sprintf("%+v", eMsg))
	}
	return nil
}

/*
 mix model to execute rules

 in this mode, it will not consider the priority，and it also concurrently to execute rules
 first to execute the most high priority rule，then concurrently to execute last rules without consider the priority
*/
func (g *Gengine) ExecuteMixModel(rb *builder.RuleBuilder) error {

	//check rb
	if rb == nil {
		return errors.New("ruleBuilder is nil")
	}

	g.returnResult = make(map[string]interface{})

	if len(rb.Kc.SortRules) == 0 {
		return errors.New("no rule has been injected into engine! ")
	}

	rules := rb.Kc.SortRules
	v, e, bx := rules[0].Execute(rb.Dc)
	if bx {
		g.addResult(rules[0].RuleName, v)
	}

	if e != nil {
		return errors.New(fmt.Sprintf("the most high priority rule: \"%s\"  executed, error:\n %+v", rules[0].RuleName, e))
	}

	var errLock sync.Mutex
	var eMsg []string

	if (len(rules) - 1) >= 1 {
		var wg sync.WaitGroup
		wg.Add(len(rules) - 1)
		for _, r := range rules[1:] {
			rr := r
			go func() {
				v, e, bx := rr.Execute(rb.Dc)
				if bx {
					g.addResult(rr.RuleName, v)
				}
				if e != nil {
					errLock.Lock()
					eMsg = append(eMsg, fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", rr.RuleName, e))
					errLock.Unlock()
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}

	if len(eMsg) > 0 {
		return errors.New(fmt.Sprintf("%+v", eMsg))
	}
	return nil
}

/**
 mix execute model

base type :golang translate value
not base type: golang translate pointer

if stopTag become true,it will not continue to execute
stopTag is a name given by user, and user can use it  to control rules execute behavior in rules, it can improve performance

it used in this scene:
where the first rule execute finished, you don't want to execute to the last rules, you can use sTag to control it out of gengine

*/
func (g *Gengine) ExecuteMixModelWithStopTagDirect(rb *builder.RuleBuilder, sTag *Stag) error {

	//check rb
	if rb == nil {
		return errors.New("ruleBuilder is nil")
	}

	g.returnResult = make(map[string]interface{})

	if len(rb.Kc.SortRules) == 0 {
		return errors.New("no rule has been injected into engine! ")
	}

	rules := rb.Kc.SortRules
	v, e, bx := rules[0].Execute(rb.Dc)
	if bx {
		g.addResult(rules[0].RuleName, v)
	}
	if e != nil {
		return errors.New(fmt.Sprintf("the most high priority rule: \"%s\"  executed, error:\n %+v", rules[0].RuleName, e))
	}

	var errLock sync.Mutex
	var eMsg []string

	if !sTag.StopTag {
		if (len(rules) - 1) >= 1 {
			var wg sync.WaitGroup
			wg.Add(len(rules) - 1)
			for _, r := range rules[1:] {
				rr := r
				go func() {
					v, e, bx := rr.Execute(rb.Dc)
					if bx {
						g.addResult(rr.RuleName, v)
					}
					if e != nil {
						errLock.Lock()
						eMsg = append(eMsg, fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", rr.RuleName, e))
						errLock.Unlock()
					}
					wg.Done()
				}()
			}
			wg.Wait()
		}
	}

	if len(eMsg) > 0 {
		return errors.New(fmt.Sprintf("%+v", eMsg))
	}
	return nil
}

/**
user can choose specified name rules to run with sort, and it will continue to execute the last rules,even if there rule execute error
*/
func (g *Gengine) ExecuteSelectedRules(rb *builder.RuleBuilder, names []string) error {

	//check rb
	if rb == nil {
		return errors.New("ruleBuilder is nil")
	}

	g.returnResult = make(map[string]interface{})

	if len(rb.Kc.RuleEntities) == 0 {
		return errors.New("no rule has been injected into engine! ")
	}

	var rules []*base.RuleEntity
	for _, name := range names {
		if ruleEntity, ok := rb.Kc.RuleEntities[name]; ok {
			rr := ruleEntity
			rules = append(rules, rr)
		} else {
			log.Errorf("no such rule named: \"%s\"", name)
		}
	}

	if len(rules) < 1 {
		return errors.New(fmt.Sprintf("no rules have been selected, names=%+v", names))
	}

	if len(rules) >= 2 {
		sort.SliceStable(rules, func(i, j int) bool {
			return rules[i].Salience > rules[j].Salience
		})
	}

	var eMsg []string
	for _, rule := range rules {
		rr := rule
		v, e, bx := rr.Execute(rb.Dc)
		if bx {
			g.addResult(rr.RuleName, v)
		}
		if e != nil {
			eMsg = append(eMsg, fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", rr.RuleName, e))
		}
	}

	if len(eMsg) > 0 {
		return errors.New(fmt.Sprintf("%+v", eMsg))
	}
	return nil
}

/**
user can choose specified name rules to run with sort
b bool:control whether continue to execute last rules ,when a rule execute error; if b == true ,the func is same to ExecuteSelectedRules
*/
func (g *Gengine) ExecuteSelectedRulesWithControl(rb *builder.RuleBuilder, b bool, names []string) error {

	//check rb
	if rb == nil {
		return errors.New("ruleBuilder is nil")
	}

	g.returnResult = make(map[string]interface{})

	if len(rb.Kc.SortRules) == 0 {
		return errors.New("no rule has been injected into engine! ")
	}

	var rules []*base.RuleEntity
	for _, name := range names {
		if ruleEntity, ok := rb.Kc.RuleEntities[name]; ok {
			rr := ruleEntity
			rules = append(rules, rr)
		} else {
			log.Errorf("no such rule named: \"%s\"", name)
		}
	}

	if len(rules) < 1 {
		return errors.New(fmt.Sprintf("no rule has been selected, names=%+v", names))
	}

	if len(rules) >= 2 {
		sort.SliceStable(rules, func(i, j int) bool {
			return rules[i].Salience > rules[j].Salience
		})
	}

	var eMsg []string
	for _, rule := range rules {
		rr := rule
		v, e, bx := rr.Execute(rb.Dc)
		if bx {
			g.addResult(rr.RuleName, v)
		}
		if e != nil {
			if b {
				eMsg = append(eMsg, fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", rr.RuleName, e))
			} else {
				return errors.New(fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", rr.RuleName, e))
			}
		}
	}

	if len(eMsg) > 0 {
		return errors.New(fmt.Sprintf("%+v", eMsg))
	}
	return nil
}

/**
user can choose specified name rules to run with given sorted name
b bool:control whether continue to execute last rules ,when a rule execute error; if b == true ,the func is same to ExecuteSelectedRules

gengine won't sort the rules by the salience, and the executed order will  based on the user given sorted names
*/
func (g *Gengine) ExecuteSelectedRulesWithControlAsGivenSortedName(rb *builder.RuleBuilder, b bool, sortedNames []string) error {

	//check rb
	if rb == nil {
		return errors.New("ruleBuilder is nil")
	}

	g.returnResult = make(map[string]interface{})

	if len(rb.Kc.SortRules) == 0 {
		return errors.New("no rule has been injected into engine! ")
	}

	var rules []*base.RuleEntity
	for _, name := range sortedNames {
		if ruleEntity, ok := rb.Kc.RuleEntities[name]; ok {
			rr := ruleEntity
			rules = append(rules, rr)
		} else {
			log.Errorf("no such rule named: \"%s\"", name)
		}
	}

	if len(rules) < 1 {
		return errors.New(fmt.Sprintf("no rule has been selected, names=%+v", sortedNames))
	}

	var eMsg []string
	for _, rule := range rules {
		rr := rule
		v, e, bx := rr.Execute(rb.Dc)
		if bx {
			g.addResult(rr.RuleName, v)
		}
		if e != nil {
			if b {
				eMsg = append(eMsg, fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", rr.RuleName, e))
			} else {
				return errors.New(fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", rr.RuleName, e))
			}
		}
	}

	if len(eMsg) > 0 {
		return errors.New(fmt.Sprintf("%+v", eMsg))
	}
	return nil
}

/**
user can choose specified name rules to run with sort
b bool:control whether continue to execute last rules ,when a rule execute error; if b == true ,the func is same to ExecuteSelectedRules
*/
func (g *Gengine) ExecuteSelectedRulesWithControlAndStopTag(rb *builder.RuleBuilder, b bool, sTag *Stag, names []string) error {

	//check rb
	if rb == nil {
		return errors.New("ruleBuilder is nil")
	}

	g.returnResult = make(map[string]interface{})

	if len(rb.Kc.SortRules) == 0 {
		return errors.New("no rule has been injected into engine! ")
	}

	var rules []*base.RuleEntity
	for _, name := range names {
		if ruleEntity, ok := rb.Kc.RuleEntities[name]; ok {
			rr := ruleEntity
			rules = append(rules, rr)
		} else {
			log.Errorf("no such rule named: \"%s\"", name)
		}
	}

	if len(rules) < 1 {
		return errors.New(fmt.Sprintf("no rule has been selected, names=%+v", names))
	}

	if len(rules) >= 2 {
		sort.SliceStable(rules, func(i, j int) bool {
			return rules[i].Salience > rules[j].Salience
		})
	}

	var eMsg []string
	for _, rule := range rules {
		rr := rule
		v, e, bx := rr.Execute(rb.Dc)
		if bx {
			g.addResult(rr.RuleName, v)
		}
		if e != nil {
			if b {
				eMsg = append(eMsg, fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", rr.RuleName, e))
			} else {
				return errors.New(fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", rr.RuleName, e))
			}
		}

		if sTag.StopTag {
			break
		}
	}

	if len(eMsg) > 0 {
		return errors.New(fmt.Sprintf("%+v", eMsg))
	}
	return nil
}

/**
user can choose specified name rules to run with given sorted name
b bool:control whether continue to execute last rules ,when a rule execute error; if b == true ,the func is same to ExecuteSelectedRules

gengine won't sort the rules by the salience, and the executed order will  based on the user given sorted names
*/
func (g *Gengine) ExecuteSelectedRulesWithControlAndStopTagAsGivenSortedName(rb *builder.RuleBuilder, b bool, sTag *Stag, sortedNames []string) error {

	//check rb
	if rb == nil {
		return errors.New("ruleBuilder is nil")
	}

	g.returnResult = make(map[string]interface{})

	if len(rb.Kc.SortRules) == 0 {
		return errors.New("no rule has been injected into engine! ")
	}

	var rules []*base.RuleEntity
	for _, name := range sortedNames {
		if ruleEntity, ok := rb.Kc.RuleEntities[name]; ok {
			rr := ruleEntity
			rules = append(rules, rr)
		} else {
			log.Errorf("no such rule named: \"%s\"", name)
		}
	}

	if len(rules) < 1 {
		return errors.New(fmt.Sprintf("no rule has been selected, names=%+v", sortedNames))
	}

	var eMsg []string
	for _, rule := range rules {
		rr := rule
		v, e, bx := rr.Execute(rb.Dc)
		if bx {
			g.addResult(rr.RuleName, v)
		}
		if e != nil {
			if b {
				eMsg = append(eMsg, fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", rr.RuleName, e))
			} else {
				return errors.New(fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", rr.RuleName, e))
			}
		}

		if sTag.StopTag {
			break
		}
	}

	if len(eMsg) > 0 {
		return errors.New(fmt.Sprintf("%+v", eMsg))
	}
	return nil
}

/**
user can choose specified name rules to concurrent run
*/
func (g *Gengine) ExecuteSelectedRulesConcurrent(rb *builder.RuleBuilder, names []string) error {

	//check rb
	if rb == nil {
		return errors.New("ruleBuilder is nil")
	}

	g.returnResult = make(map[string]interface{})

	if len(rb.Kc.RuleEntities) == 0 {
		return errors.New("no rule has been injected into engine! ")
	}

	var rules []*base.RuleEntity
	for _, name := range names {
		if ruleEntity, ok := rb.Kc.RuleEntities[name]; ok {
			rr := ruleEntity
			rules = append(rules, rr)
		} else {
			log.Errorf("no such rule named: \"%s\"", name)
		}
	}

	if len(rules) == 0 {
		return errors.New(fmt.Sprintf("no rule has been selected, names=%+v", names))
	}

	if len(rules) == 1 {
		v, e, bx := rules[0].Execute(rb.Dc)
		if bx {
			g.addResult(rules[0].RuleName, v)
		}
		if e != nil {
			return errors.New(fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", rules[0].RuleName, e))
		}
		return nil
	}

	var errLock sync.Mutex
	var eMsg []string

	// len(rule) >= 2
	var wg sync.WaitGroup
	wg.Add(len(rules))
	for _, r := range rules {
		rr := r
		go func() {
			v, e, bx := rr.Execute(rb.Dc)
			if bx {
				g.addResult(rr.RuleName, v)
			}
			if e != nil {
				errLock.Lock()
				eMsg = append(eMsg, fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", rr.RuleName, e))
				errLock.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if len(eMsg) > 0 {
		return errors.New(fmt.Sprintf("%+v", eMsg))
	}
	return nil
}

/**
user can choose specified name rules to run with mix model
*/
func (g *Gengine) ExecuteSelectedRulesMixModel(rb *builder.RuleBuilder, names []string) error {

	//check rb
	if rb == nil {
		return errors.New("ruleBuilder is nil")
	}

	g.returnResult = make(map[string]interface{})

	if len(rb.Kc.RuleEntities) == 0 {
		return errors.New("no rule has been injected into engine! ")
	}

	var rules []*base.RuleEntity
	for _, name := range names {
		if ruleEntity, ok := rb.Kc.RuleEntities[name]; ok {
			rr := ruleEntity
			rules = append(rules, rr)
		} else {
			log.Errorf("no such rule named: \"%s\"", name)
		}
	}

	if len(rules) == 0 {
		return errors.New(fmt.Sprintf("no rule has been selected, names=%+v", names))
	}

	if len(rules) == 1 {
		v, e, bx := rules[0].Execute(rb.Dc)
		if bx {
			g.addResult(rules[0].RuleName, v)
		}
		if e != nil {
			return errors.New(fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", rules[0].RuleName, e))
		}
		return nil
	}

	sort.SliceStable(rules, func(i, j int) bool {
		return rules[i].Salience > rules[j].Salience
	})

	if len(rules) == 2 {
		for _, r := range rules {
			v, err, bx := r.Execute(rb.Dc)
			if bx {
				g.addResult(r.RuleName, v)
			}
			if err != nil {
				return errors.New(fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", r.RuleName, err))
			}
		}
		return nil
	}

	// rLen >= 3
	v, e, bx := rules[0].Execute(rb.Dc)
	if bx {
		g.addResult(rules[0].RuleName, v)
	}
	if e != nil {
		return errors.New(fmt.Sprintf("the most high priority rule: \"%s\"  executed, error:\n %+v", rules[0].RuleName, e))
	}

	var errLock sync.Mutex
	var eMsg []string

	var wg sync.WaitGroup
	wg.Add(len(rules) - 1)
	for _, r := range rules[1:] {
		rr := r
		go func() {
			v, e, bx := rr.Execute(rb.Dc)
			if bx {
				g.addResult(rr.RuleName, v)
			}
			if e != nil {
				errLock.Lock()
				eMsg = append(eMsg, fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", rr.RuleName, e))
				errLock.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if len(eMsg) > 0 {
		return errors.New(fmt.Sprintf("%+v", eMsg))
	}
	return nil
}

//inverse mix model
func (g *Gengine) ExecuteInverseMixModel(rb *builder.RuleBuilder) error {
	//check rb
	if rb == nil {
		return errors.New("ruleBuilder is nil")
	}

	g.returnResult = make(map[string]interface{})

	rules := rb.Kc.SortRules
	length := len(rules)
	if length == 0 {
		return errors.New("no rule has been injected into engine! ")
	}

	if length <= 2 {
		for _, r := range rules {
			v, e, bx := r.Execute(rb.Dc)
			if bx {
				g.addResult(r.RuleName, v)
			}
			if e != nil {
				return errors.New(fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", r.RuleName, e))
			}
		}
		return nil
	}

	var errLock sync.Mutex
	var eMsg []string

	var wg sync.WaitGroup
	wg.Add(length - 1)
	for _, r := range rules[:length-1] {
		rr := r
		go func() {
			v, e, bx := rr.Execute(rb.Dc)
			if bx {
				g.addResult(rr.RuleName, v)
			}
			if e != nil {
				errLock.Lock()
				eMsg = append(eMsg, fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", rr.RuleName, e))
				errLock.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if len(eMsg) > 0 {
		return errors.New(fmt.Sprintf("%+v", eMsg))
	}

	v, e, bx := rules[length-1].Execute(rb.Dc)
	if bx {
		g.addResult(rules[length-1].RuleName, v)
	}
	return e
}

//inverse mix model with user selected
func (g *Gengine) ExecuteSelectedRulesInverseMixModel(rb *builder.RuleBuilder, names []string) error {
	//check rb
	if rb == nil {
		return errors.New("ruleBuilder is nil")
	}

	g.returnResult = make(map[string]interface{})

	var rules []*base.RuleEntity
	//choose user need!
	for _, name := range names {
		if re, ok := rb.Kc.RuleEntities[name]; ok {
			rules = append(rules, re)
		} else {
			log.Errorf("no such rule named: \"%s\"", name)
		}
	}

	length := len(rules)
	if length == 0 {
		return errors.New("no rule has been selected to execute! ")
	}

	//resort
	sort.SliceStable(rules, func(i, j int) bool {
		return rules[i].Salience > rules[j].Salience
	})

	if length <= 2 {
		for _, r := range rules {
			v, e, bx := r.Execute(rb.Dc)
			if bx {
				g.addResult(r.RuleName, v)
			}
			if e != nil {
				return errors.New(fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", r.RuleName, e))
			}
		}
		return nil
	}

	var errLock sync.Mutex
	var eMsg []string

	var wg sync.WaitGroup
	wg.Add(length - 1)
	for _, r := range rules[:length-1] {
		rr := r
		go func() {
			v, e, bx := rr.Execute(rb.Dc)
			if bx {
				g.addResult(rr.RuleName, v)
			}
			if e != nil {
				errLock.Lock()
				eMsg = append(eMsg, fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", rr.RuleName, e))
				errLock.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if len(eMsg) > 0 {
		return errors.New(fmt.Sprintf("%+v", eMsg))
	}

	v, e, bx := rules[length-1].Execute(rb.Dc)
	if bx {
		g.addResult(rules[length-1].RuleName, v)
	}
	return e
}

// 1.first n piece rules to sort execute based on priority
// 2.bool b means: when in sort execute stage,if a rule execute error whether continue to execute the last all rules,
//   if b == true, means continue, if false, means stop and return
// 3.then m piece rules to concurrent execute based without priority
func (g *Gengine) ExecuteNSortMConcurrent(nSort, mConcurrent int, rb *builder.RuleBuilder, b bool) error {

	//check rb
	if rb == nil {
		return errors.New("ruleBuilder is nil")
	}

	g.returnResult = make(map[string]interface{})

	//strictly params check
	if nSort <= 0 {
		return errors.New(fmt.Sprintf("params should be bigger than 0, nSort=%d", nSort))
	}

	if mConcurrent <= 0 {
		return errors.New(fmt.Sprintf("params should be bigger than 0, mConcurrent=%d", nSort))
	}

	if nSort+mConcurrent > len(rb.Kc.SortRules) {
		return errors.New(fmt.Sprintf("not enough rules to complete N-M execute model, nSort+mConcurrent = %d, while rules.len=%d", nSort+mConcurrent, len(rb.Kc.SortRules)))
	}

	var errLock sync.Mutex
	var eMsg []string

	//nSort
	nRules := rb.Kc.SortRules[:nSort]
	for _, rule := range nRules {
		v, e, bx := rule.Execute(rb.Dc)
		if bx {
			g.addResult(rule.RuleName, v)
		}
		if b {
			if e != nil {
				eMsg = append(eMsg, fmt.Sprintf("%+v", e))
			}
		} else {
			return e
		}
	}

	//mConcurrent
	mRules := rb.Kc.SortRules[nSort:][:mConcurrent]
	var wg sync.WaitGroup
	wg.Add(mConcurrent)
	for _, r := range mRules {
		rr := r
		go func() {
			v, e, bx := rr.Execute(rb.Dc)
			if bx {
				g.addResult(rr.RuleName, v)
			}
			if e != nil {
				errLock.Lock()
				eMsg = append(eMsg, fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", rr.RuleName, e))
				errLock.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if len(eMsg) > 0 {
		return errors.New(fmt.Sprintf("%+v", eMsg))
	}

	return nil
}

// 1. first n piece rules to concurrent execute based without priority
// 2. bool b means: after concurrent execute stage,if a rule execute error whether continue to execute the last all rules,
//    if b == true, means continue, if false, means stop and return
// 3. then m piece rules to sort execute based on priority
func (g *Gengine) ExecuteNConcurrentMSort(nConcurrent, mSort int, rb *builder.RuleBuilder, b bool) error {
	//check rb
	if rb == nil {
		return errors.New("ruleBuilder is nil")
	}

	g.returnResult = make(map[string]interface{})

	//strictly params check
	if nConcurrent <= 0 {
		return errors.New(fmt.Sprintf("params should be bigger than 0, nConcurrent=%d", nConcurrent))
	}

	if mSort <= 0 {
		return errors.New(fmt.Sprintf("params should be bigger than 0, mSort=%d", mSort))
	}

	if nConcurrent+mSort > len(rb.Kc.SortRules) {
		return errors.New(fmt.Sprintf("not enough rules to complete N-M execute model, nConcurrent+mSort = %d, while rules.len=%d", nConcurrent+mSort, len(rb.Kc.SortRules)))
	}

	var errLock sync.Mutex
	var eMsg []string

	//nConcurrent
	nRules := rb.Kc.SortRules[:nConcurrent]
	var wg sync.WaitGroup
	wg.Add(nConcurrent)
	for _, r := range nRules {
		rr := r
		go func() {
			v, e, bx := rr.Execute(rb.Dc)
			if bx {
				g.addResult(rr.RuleName, v)
			}
			if e != nil {
				errLock.Lock()
				eMsg = append(eMsg, fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", rr.RuleName, e))
				errLock.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if !b {
		if len(eMsg) > 0 {
			return errors.New(fmt.Sprintf("%+v", eMsg))
		}
	}

	//mSort
	mRules := rb.Kc.SortRules[nConcurrent:][:mSort]
	for _, rule := range mRules {
		v, e, bx := rule.Execute(rb.Dc)
		if bx {
			g.addResult(rule.RuleName, v)
		}
		if b {
			if e != nil {
				eMsg = append(eMsg, fmt.Sprintf("%+v", e))
			}
		} else {
			return e
		}
	}

	if len(eMsg) > 0 {
		return errors.New(fmt.Sprintf("%+v", eMsg))
	}

	return nil
}

// 1. first n piece rules to concurrent execute based without priority
// 2. bool b means: if the first stage executed error, whether continue to execute the next concurrent stage
//    if b == true,   means continue, if false, means stop and return
// 3. then m piece rules to concurrent execute based without priority
func (g *Gengine) ExecuteNConcurrentMConcurrent(nConcurrent, mConcurrent int, rb *builder.RuleBuilder, b bool) error {

	//check rb
	if rb == nil {
		return errors.New("ruleBuilder is nil")
	}

	g.returnResult = make(map[string]interface{})

	//strictly params check
	if nConcurrent <= 0 {
		return errors.New(fmt.Sprintf("params should be bigger than 0, nConcurrent=%d", nConcurrent))
	}

	if mConcurrent <= 0 {
		return errors.New(fmt.Sprintf("params should be bigger than 0, mConcurrent=%d", mConcurrent))
	}

	if nConcurrent+mConcurrent > len(rb.Kc.SortRules) {
		return errors.New(fmt.Sprintf("not enough rules to complete N-M execute model, nConcurrent+mConcurrent = %d, while rules.len=%d", nConcurrent+mConcurrent, len(rb.Kc.SortRules)))
	}

	var errLock sync.Mutex
	var eMsg []string

	//nConcurrent
	nRules := rb.Kc.SortRules[:nConcurrent]
	var nwg sync.WaitGroup
	nwg.Add(nConcurrent)
	for _, r := range nRules {
		rr := r
		go func() {
			v, e, bx := rr.Execute(rb.Dc)
			if bx {
				g.addResult(rr.RuleName, v)
			}
			if e != nil {
				errLock.Lock()
				eMsg = append(eMsg, fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", rr.RuleName, e))
				errLock.Unlock()
			}
			nwg.Done()
		}()
	}
	nwg.Wait()

	if !b {
		if len(eMsg) > 0 {
			return errors.New(fmt.Sprintf("%+v", eMsg))
		}
	}

	//mConcurrent
	mRules := rb.Kc.SortRules[nConcurrent:][:mConcurrent]
	var mwg sync.WaitGroup
	mwg.Add(mConcurrent)
	for _, r := range mRules {
		rr := r
		go func() {
			v, e, bx := rr.Execute(rb.Dc)
			if bx {
				g.addResult(rr.RuleName, v)
			}
			if e != nil {
				errLock.Lock()
				eMsg = append(eMsg, fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", rr.RuleName, e))
				errLock.Unlock()
			}
			mwg.Done()
		}()
	}
	mwg.Wait()

	if len(eMsg) > 0 {
		return errors.New(fmt.Sprintf("%+v", eMsg))
	}

	return nil
}

// 0.based on selected rules
// 1.first n piece rules to sort execute based on priority
// 2.bool b means: when in sort execute stage,if a rule execute error whether continue to execute the last all rules,
//   if b == true, means continue, if false, means stop and return
// 3.then m piece rules to concurrent execute based without priority
func (g *Gengine) ExecuteSelectedNSortMConcurrent(nSort, mConcurrent int, rb *builder.RuleBuilder, b bool, names []string) error {

	//check rb
	if rb == nil {
		return errors.New("ruleBuilder is nil")
	}

	g.returnResult = make(map[string]interface{})

	//strictly params check
	if nSort <= 0 {
		return errors.New(fmt.Sprintf("params should be bigger than 0, nSort=%d", nSort))
	}

	if mConcurrent <= 0 {
		return errors.New(fmt.Sprintf("params should be bigger than 0, mConcurrent=%d", nSort))
	}

	if nSort+mConcurrent != len(names) {
		return errors.New(fmt.Sprintf("selected rules' len should equals the nSort+mConcurrent, selected rules' len=%d, nSort+mConcurrent=%d", len(names), nSort+mConcurrent))
	}

	if nSort+mConcurrent > len(rb.Kc.SortRules) {
		return errors.New(fmt.Sprintf("not enough selected rules to complete N-M execute model, nSort+mConcurrent = %d, while rules.len=%d", nSort+mConcurrent, len(rb.Kc.SortRules)))
	}

	//selected based on names
	var rules []*base.RuleEntity
	for _, v := range names {
		if rule, ok := rb.Kc.RuleEntities[v]; ok {
			rules = append(rules, rule)
		} else {
			return errors.New(fmt.Sprintf("not exist rule:%s", rule.RuleName))
		}
	}

	//resort
	sort.SliceStable(rules, func(i, j int) bool {
		return rules[i].Salience > rules[j].Salience
	})

	var errLock sync.Mutex
	var eMsg []string

	//nSort
	nRules := rules[:nSort]
	for _, rule := range nRules {
		v, e, bx := rule.Execute(rb.Dc)
		if bx {
			g.addResult(rule.RuleName, v)
		}
		if b {
			if e != nil {
				eMsg = append(eMsg, fmt.Sprintf("%+v", e))
			}
		} else {
			return e
		}
	}

	//mConcurrent
	mRules := rules[nSort:][:mConcurrent]
	var wg sync.WaitGroup
	wg.Add(mConcurrent)
	for _, r := range mRules {
		rr := r
		go func() {
			v, e, bx := rr.Execute(rb.Dc)
			if bx {
				g.addResult(rr.RuleName, v)
			}
			if e != nil {
				errLock.Lock()
				eMsg = append(eMsg, fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", rr.RuleName, e))
				errLock.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if len(eMsg) > 0 {
		return errors.New(fmt.Sprintf("%+v", eMsg))
	}

	return nil
}

// 0. based on selected rules
// 1. first n piece rules to concurrent execute based without priority
// 2. bool b means: after concurrent execute stage,if a rule execute error whether continue to execute the last all rules,
//    if b == true, means continue, if false, means stop and return
// 3. then m piece rules to sort execute based on priority
func (g *Gengine) ExecuteSelectedNConcurrentMSort(nConcurrent, mSort int, rb *builder.RuleBuilder, b bool, names []string) error {

	//check rb
	if rb == nil {
		return errors.New("ruleBuilder is nil")
	}

	g.returnResult = make(map[string]interface{})

	//strictly params check
	if nConcurrent <= 0 {
		return errors.New(fmt.Sprintf("params should be bigger than 0, nConcurrent=%d", nConcurrent))
	}

	if mSort <= 0 {
		return errors.New(fmt.Sprintf("params should be bigger than 0, mSort=%d", mSort))
	}

	if nConcurrent+mSort != len(names) {
		return errors.New(fmt.Sprintf("selected rules' len should equals the nConcurrent+mSort, selected rules' len=%d, nConcurrent+mSort=%d", len(names), nConcurrent+mSort))
	}

	if nConcurrent+mSort > len(rb.Kc.SortRules) {
		return errors.New(fmt.Sprintf("not enough selected rules to complete N-M execute model, nConcurrent+mSort = %d, while rules.len=%d", nConcurrent+mSort, len(rb.Kc.SortRules)))
	}

	//selected based on names
	var rules []*base.RuleEntity
	for _, v := range names {
		if rule, ok := rb.Kc.RuleEntities[v]; ok {
			rules = append(rules, rule)
		} else {
			return errors.New(fmt.Sprintf("not exist rule:%s", rule.RuleName))
		}
	}

	//resort
	sort.SliceStable(rules, func(i, j int) bool {
		return rules[i].Salience > rules[j].Salience
	})

	var errLock sync.Mutex
	var eMsg []string

	//nConcurrent
	nRules := rules[:nConcurrent]
	var wg sync.WaitGroup
	wg.Add(nConcurrent)
	for _, r := range nRules {
		rr := r
		go func() {
			v, e, bx := rr.Execute(rb.Dc)
			if bx {
				g.addResult(rr.RuleName, v)
			}
			if e != nil {
				errLock.Lock()
				eMsg = append(eMsg, fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", rr.RuleName, e))
				errLock.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if !b {
		if len(eMsg) > 0 {
			return errors.New(fmt.Sprintf("%+v", eMsg))
		}
	}

	//mSort
	mRules := rules[nConcurrent:][:mSort]
	for _, rule := range mRules {
		v, e, bx := rule.Execute(rb.Dc)
		if bx {
			g.addResult(rule.RuleName, v)
		}
		if b {
			if e != nil {
				eMsg = append(eMsg, fmt.Sprintf("%+v", e))
			}
		} else {
			return e
		}
	}

	if len(eMsg) > 0 {
		return errors.New(fmt.Sprintf("%+v", eMsg))
	}

	return nil
}

// based on selected rules
// 1. first n piece rules to concurrent execute based without priority
// 2. bool b means: if the first stage executed error, whether continue to execute the next concurrent stage
//    if b == true,   means continue, if false, means stop and return
// 3. then m piece rules to concurrent execute based without priority
func (g *Gengine) ExecuteSelectedNConcurrentMConcurrent(nConcurrent, mConcurrent int, rb *builder.RuleBuilder, b bool, names []string) error {

	//check rb
	if rb == nil {
		return errors.New("ruleBuilder is nil")
	}

	g.returnResult = make(map[string]interface{})

	//strictly params check
	if nConcurrent <= 0 {
		return errors.New(fmt.Sprintf("params should be bigger than 0, nConcurrent=%d", nConcurrent))
	}

	if mConcurrent <= 0 {
		return errors.New(fmt.Sprintf("params should be bigger than 0, mConcurrent=%d", mConcurrent))
	}

	if nConcurrent+mConcurrent != len(names) {
		return errors.New(fmt.Sprintf("selected rules' len should equals the nConcurrent+mConcurrent, selected rules' len=%d, nConcurrent+mConcurrent=%d", len(names), nConcurrent+mConcurrent))
	}

	if nConcurrent+mConcurrent > len(rb.Kc.SortRules) {
		return errors.New(fmt.Sprintf("not enough selected rules to complete N-M execute model, nConcurrent+mConcurrent = %d, while rules.len=%d", nConcurrent+mConcurrent, len(rb.Kc.SortRules)))
	}

	//selected based on names
	var rules []*base.RuleEntity
	for _, v := range names {
		if rule, ok := rb.Kc.RuleEntities[v]; ok {
			rules = append(rules, rule)
		} else {
			return errors.New(fmt.Sprintf("not exist rule:%s", rule.RuleName))
		}
	}

	//resort
	sort.SliceStable(rules, func(i, j int) bool {
		return rules[i].Salience > rules[j].Salience
	})

	var errLock sync.Mutex
	var eMsg []string

	//nConcurrent
	nRules := rules[:nConcurrent]
	var nwg sync.WaitGroup
	nwg.Add(nConcurrent)
	for _, r := range nRules {
		rr := r
		go func() {
			v, e, bx := rr.Execute(rb.Dc)
			if bx {
				g.addResult(rr.RuleName, v)
			}
			if e != nil {
				errLock.Lock()
				eMsg = append(eMsg, fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", rr.RuleName, e))
				errLock.Unlock()
			}
			nwg.Done()
		}()
	}
	nwg.Wait()

	if !b {
		if len(eMsg) > 0 {
			return errors.New(fmt.Sprintf("%+v", eMsg))
		}
	}

	//mConcurrent
	mRules := rules[nConcurrent:][:mConcurrent]
	var mwg sync.WaitGroup
	mwg.Add(mConcurrent)
	for _, r := range mRules {
		rr := r
		go func() {
			v, e, bx := rr.Execute(rb.Dc)
			if bx {
				g.addResult(rr.RuleName, v)
			}
			if e != nil {
				errLock.Lock()
				eMsg = append(eMsg, fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", rr.RuleName, e))
				errLock.Unlock()
			}
			mwg.Done()
		}()
	}
	mwg.Wait()

	if len(eMsg) > 0 {
		return errors.New(fmt.Sprintf("%+v", eMsg))
	}

	return nil
}

//DAG model
func (g *Gengine) ExecuteDAGModel(rb *builder.RuleBuilder, dag [][]string) error {

	//check rb
	if rb == nil {
		return errors.New("ruleBuilder is nil")
	}

	//check params
	if len(dag) == 0 {
		return nil
	}

	var errLock sync.Mutex
	var eMsg []string

	//row
	for i := 0; i < len(dag); i++ {
		//col
		var rules []*base.RuleEntity
		for j := 0; j < len(dag[i]); j++ {
			//filter the rules which do not exist.
			if rule, ok := rb.Kc.RuleEntities[dag[i][j]]; ok {
				rules = append(rules, rule)
			}
		}

		//并发执行
		if len(rules) > 0 {
			var mwg sync.WaitGroup
			mwg.Add(len(rules))
			for _, r := range rules {
				rr := r
				go func() {
					v, e, bx := rr.Execute(rb.Dc)
					if bx {
						g.addResult(rr.RuleName, v)
					}
					if e != nil {
						errLock.Lock()
						eMsg = append(eMsg, fmt.Sprintf("rule: \"%s\" executed, error:\n %+v ", rr.RuleName, e))
						errLock.Unlock()
					}
					mwg.Done()
				}()
			}
			mwg.Wait()
		}
		if len(eMsg) > 0 {
			return errors.New(fmt.Sprintf("%+v", eMsg))
		}

	}
	return nil
}
