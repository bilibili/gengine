package engine

import (
	"errors"
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/internal/base"
	parser "github.com/bilibili/gengine/internal/iantlr/alr"
	"github.com/bilibili/gengine/internal/iparser"
	"github.com/bilibili/gengine/internal/tool"
	"sync"

	"github.com/google/martian/log"
)

const (
	SortModel       = 1
	ConcurrentModel = 2
	MixModel        = 3
	InverseMixModel = 4
)

// when you use NewGenginePool, you just think of it as the connection pool of mysql, the higher QPS you want to support, the more resource you need to give
type GenginePool struct {
	runningLock  sync.Mutex
	freeGengines []*gengineWrapper

	//just for check whether a rule exist
	ruleBuilder *builder.RuleBuilder

	execModel int
	apis      map[string]interface{}

	additionLock     sync.Mutex
	additionGengines []*gengineWrapper
	additionNum      int64

	updateLock sync.Mutex
	clear      bool //whether rules has been cleared ，if true it means there is no rules in gengine

	rbSlice []*builder.RuleBuilder
	//total gengine instance number
	max int64

	getEngineLock sync.RWMutex //just one can get this lock
}

type gengineWrapper struct {
	tag         int64 // one to one between the ruleBuilder slice
	rulebuilder *builder.RuleBuilder
	gengine     *Gengine

	addition bool // when gengine resource is not enough and poollength >  minPool  and  poollength < maxPool, new gengine will be create, and it will be tagged addition=true; when poollength <  minPool it will be tagged addition=false
}

func (gw *gengineWrapper) clearInjected(keys ...string) {
	if gw == nil || gw.rulebuilder == nil || gw.rulebuilder.Dc == nil {
		return
	}
	gw.rulebuilder.Dc.Del(keys...)
}

//poolLen  -> gengine pool length to init
//em       -> rule execute model: 1 sort model; 2 :concurrent model; 3: mix model; 4: inverse model
//rStr  -> rules string
//apiOuter -> which user want to add rule engine to use
// just init once!!!

// best practise：
// when there has cost-time operate in your rule or you want to support high concurrency(> 200000QPS) , please set poolMinLen bigger Appropriately
// when you use NewGenginePool,you just think of it as the connection pool of mysql, the higher QPS you want to support, the more resource you need to give
func NewGenginePool(poolMinLen, poolMaxLen int64, em int, rulesStr string, apiOuter map[string]interface{}) (*GenginePool, error) {

	if !(0 < poolMinLen && poolMinLen < poolMaxLen) {
		return nil, errors.New("pool length must be bigger than 0, and poolMaxLen must bigger than poolMinLen")
	}

	if em != SortModel && em != ConcurrentModel && em != MixModel && em != InverseMixModel {
		return nil, errors.New(fmt.Sprintf("exec model must be SORT_MODEL(1) or CONCOURRENT_MODEL(2) or MIX_MODEL(3) or INVERSE_MIX_MODEL(4), now it is %d", em))
	}

	fg := make([]*gengineWrapper, poolMinLen)
	for i := int64(0); i < poolMinLen; i++ {
		fg[i] = &gengineWrapper{
			tag:      i,
			gengine:  NewGengine(),
			addition: false,
		}
	}

	ag := make([]*gengineWrapper, poolMaxLen-poolMinLen)
	for j := int64(0); j < poolMaxLen-poolMinLen; j++ {
		ag[j] = &gengineWrapper{
			tag:      j + poolMinLen,
			gengine:  NewGengine(),
			addition: true,
		}
	}

	srcRb, e := makeRuleBuilder(rulesStr, apiOuter)
	if e != nil {
		return nil, e
	}

	rbs := make([]*builder.RuleBuilder, poolMaxLen)
	for i := 0; i < int(poolMaxLen); i++ {
		dataContext := context.NewDataContext()
		if apiOuter != nil {
			for k, v := range apiOuter {
				dataContext.Add(k, v)
			}
		}
		rb := builder.NewRuleBuilder(dataContext)
		rb.Kc = srcRb.Kc
		rbs[i] = rb
	}

	p := &GenginePool{
		ruleBuilder:      srcRb,
		freeGengines:     fg,
		apis:             apiOuter,
		execModel:        em,
		additionNum:      poolMaxLen - poolMinLen,
		additionGengines: ag,
		clear:            false,
		rbSlice:          rbs,
		max:              poolMaxLen,
	}
	return p, nil
}

//this could ensure make thread safety!
func makeRuleBuilder(ruleStr string, apiOuter map[string]interface{}) (*builder.RuleBuilder, error) {
	dataContext := context.NewDataContext()
	if apiOuter != nil {
		for k, v := range apiOuter {
			dataContext.Add(k, v)
		}
	}

	rb := builder.NewRuleBuilder(dataContext)
	if ruleStr != "" {
		if e := rb.BuildRuleFromString(ruleStr); e != nil {
			return nil, errors.New(fmt.Sprintf("build rule from string err: %+v", e))
		}
	} else {
		return nil, errors.New("the ruleStr is \"\"")
	}
	return rb, nil
}

// if there is no enough gengine source, no request will take a lock
func (gp *GenginePool) getGengine() (*gengineWrapper, error) {

	for {
		gp.getEngineLock.Lock()
		//check if there has enough resource in pool
		numFree := len(gp.freeGengines)
		if numFree > 0 {
			gp.runningLock.Lock()
			gw := gp.freeGengines[0]
			gp.freeGengines = gp.freeGengines[1:]
			gp.runningLock.Unlock()
			gp.getEngineLock.Unlock()
			return gw, nil
		}

		//check if there has addition resource
		numAddition := len(gp.additionGengines)
		if numAddition > 0 {
			gp.additionLock.Lock()
			gw := gp.additionGengines[0]
			gp.additionGengines = gp.additionGengines[1:]
			gp.additionLock.Unlock()
			gp.getEngineLock.Unlock()
			return gw, nil
		}

		gp.getEngineLock.Unlock()
	}
}

// async return gengine resource to pool,and update the rules
func (gp *GenginePool) putGengineLocked(gw *gengineWrapper) {
	//addition resource
	go func() {
		if gw.addition {
			gp.additionLock.Lock()
			gp.additionGengines = append(gp.additionGengines, gw)
			gp.additionLock.Unlock()
		} else {
			gp.runningLock.Lock()
			gp.freeGengines = append(gp.freeGengines, gw)
			gp.runningLock.Unlock()
		}
	}()
}

//sync method
//update the all rules in all engine in the pool
//update success: return nil
//update failed: return error
// this is very different from connection pool,
//connection pool just need to init once
//while gengine pool need to update the rules whenever the user want to init
func (gp *GenginePool) UpdatePooledRules(ruleStr string) error {
	//check the rules
	gp.updateLock.Lock()
	defer gp.updateLock.Unlock()

	rbi, e := makeRuleBuilder(ruleStr, gp.apis)
	if e != nil {
		return e
	}

	if len(rbi.Kc.RuleEntities) == 0 {
		return errors.New(fmt.Sprintf("if you want to clear all rules, use method \"pool.ClearPoolRules()\""))
	}

	gp.ruleBuilder = rbi
	for i := 0; i < int(gp.max); i++ {
		gp.rbSlice[i].Kc = gp.ruleBuilder.Kc
	}

	gp.clear = false
	return nil
}

func getKc(ruleString string) (*base.KnowledgeContext, error) {

	in := antlr.NewInputStream(ruleString)
	lexer := parser.NewgengineLexer(in)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	kc := base.NewKnowledgeContext()

	listener := iparser.NewGengineParserListener(kc)
	psr := parser.NewgengineParser(stream)
	psr.BuildParseTrees = true

	errListener := iparser.NewGengineErrorListener()
	psr.AddErrorListener(errListener)
	antlr.ParseTreeWalkerDefault.Walk(listener, psr.Primary())

	if len(errListener.GrammarErrors) > 0 {
		return nil, errors.New(fmt.Sprintf("%+v", errListener.GrammarErrors))
	}

	if len(listener.ParseErrors) > 0 {
		return nil, errors.New(fmt.Sprintf("%+v", listener.ParseErrors))
	}

	if len(kc.RuleEntities) == 0 {
		return nil, errors.New(fmt.Sprintf("no rules to update or add."))
	}

	return kc, nil
}

func updateIncremental(kc *base.KnowledgeContext, rb *builder.RuleBuilder) {
	//copy
	newRuleEntities := make(map[string]*base.RuleEntity, len(rb.Kc.RuleEntities))
	for mk, mv := range rb.Kc.RuleEntities {
		newRuleEntities[mk] = mv
	}

	//copy
	newSortRules := make([]*base.RuleEntity, len(rb.Kc.SortRules))
	for sk, sv := range rb.Kc.SortRules {
		newSortRules[sk] = sv
	}

	//kc store the new rules
	for k, v := range kc.RuleEntities {

		if vm, ok := newRuleEntities[k]; ok {
			//repalce update
			//search
			index := rb.Kc.SortRulesIndexMap[v.RuleName]
			if v.Salience == vm.Salience {
				//replace
				newSortRules[index] = v
			} else {
				newSortRules := append(newSortRules[:index], newSortRules[index+1:]...)
				low, mid := tool.BinarySearch(newSortRules, v.Salience)

				ire := []*base.RuleEntity{v}
				if mid == 0 {
					newRe := append(ire, newSortRules[low:]...)
					newSortRules = append(newSortRules[:low], newRe...)
				} else {
					newRe := append(ire, newSortRules[mid:]...)
					newSortRules = append(newSortRules[:mid], newRe...)
				}

				//update the sort index
				indexMap := make(map[string]int)
				for k, v := range newSortRules {
					indexMap[v.RuleName] = k
				}
				rb.Kc.SortRulesIndexMap = indexMap
			}

			newRuleEntities[k] = v
		} else {
			//add update
			low, mid := tool.BinarySearch(newSortRules, v.Salience)

			ire := []*base.RuleEntity{v}
			if mid == 0 {
				newRe := append(ire, newSortRules[low:]...)
				newSortRules = append(newSortRules[:low], newRe...)
			} else {
				newRe := append(ire, newSortRules[mid:]...)
				newSortRules = append(newSortRules[:mid], newRe...)
			}

			//update the sort index
			indexMap := make(map[string]int)
			for k, v := range newSortRules {
				indexMap[v.RuleName] = k
			}
			rb.Kc.SortRulesIndexMap = indexMap

			newRuleEntities[k] = v
		}
	}

	rb.Kc.RuleEntities = newRuleEntities
	rb.Kc.SortRules = newSortRules
}

//sync method
//incremental update the rules in all engine in the pool
//incremental update success: return nil
//incremental update failed: return error
// if a rule already exists, this method will use the new rule to replace the old one
// if a rule doesn't exist, this method will add the new rule to the existed rules list
//see: func (builder *RuleBuilder)BuildRuleWithIncremental(ruleString string) in rule_builder.go
func (gp *GenginePool) UpdatePooledRulesIncremental(ruleStr string) error {
	gp.updateLock.Lock()
	defer gp.updateLock.Unlock()

	//compile
	kci, e := getKc(ruleStr)
	if e != nil {
		return e
	}

	//update main
	updateIncremental(kci, gp.ruleBuilder)

	//update instance
	for i := 0; i < int(gp.max); i++ {
		gp.rbSlice[i].Kc = gp.ruleBuilder.Kc
	}

	gp.clear = false
	return nil
}

//clear all rules in engine in pool
func (gp *GenginePool) ClearPoolRules() {
	gp.updateLock.Lock()
	defer gp.updateLock.Unlock()
	gp.ruleBuilder = nil
	gp.clear = true
	for i := 0; i < int(gp.max); i++ {
		gp.rbSlice[i].Kc.ClearRules()
	}
}

//remove rules
func (gp *GenginePool) RemoveRules(ruleNames []string) error {
	gp.updateLock.Lock()
	defer gp.updateLock.Unlock()

	e := gp.ruleBuilder.RemoveRules(ruleNames)
	if e != nil {
		return e
	}

	for _, rb := range gp.rbSlice {
		_ = rb.RemoveRules(ruleNames)
	}
	return nil
}

//plugin_exportName_apiName.so
// _ is a separator
//plugin is prefix
//exportName is user export in plugin file
//apiName is plugin used in gengine
func (gp *GenginePool) PluginLoader(absolutePathOfSO string) error {
	gp.updateLock.Lock()
	defer gp.updateLock.Unlock()

	apiName, exportApi, e := gp.ruleBuilder.Dc.PluginLoader(absolutePathOfSO)
	if e != nil {
		return e
	}

	for _, rb := range gp.rbSlice {
		rb.Dc.Add(apiName, exportApi)
	}
	return nil
}

/*
1 sort model
2 concurrent model
3 mix model
4 inverse mix model
*/
func (gp *GenginePool) SetExecModel(execModel int) error {
	gp.updateLock.Lock()
	defer gp.updateLock.Unlock()
	if execModel != SortModel && execModel != ConcurrentModel && execModel != MixModel && execModel != InverseMixModel {
		return errors.New(fmt.Sprintf("exec model must be SORT_MODEL(1) or CONCOURRENT_MODEL(2) or MIX_MODEL(3) or INVERSE_MIX_MODEL(4), now it is %d", execModel))
	} else {
		gp.execModel = execModel
	}
	return nil
}

//get the execute model the user set
func (gp *GenginePool) GetExecModel() int {
	return gp.execModel
}

//check the rule whether exist
func (gp *GenginePool) IsExist(ruleNames []string) []bool {
	gp.updateLock.Lock()
	defer gp.updateLock.Unlock()

	if len(ruleNames) == 0 {
		return make([]bool, 0)
	}

	exist := make([]bool, 0)
	if gp.clear || gp.ruleBuilder == nil {
		for i := 0; i < len(ruleNames); i++ {
			exist = append(exist, false)
		}
		return exist
	}

	for _, name := range ruleNames {
		_, ok := gp.ruleBuilder.Kc.RuleEntities[name]
		exist = append(exist, ok)
	}

	return exist
}

//get the rule's salience
func (gp *GenginePool) GetRuleSalience(ruleName string) (int64, error) {
	gp.updateLock.Lock()
	defer gp.updateLock.Unlock()

	if gp.clear || gp.ruleBuilder == nil {
		return 0, errors.New("no rules in pool! ")
	}

	if rule, ok := gp.ruleBuilder.Kc.RuleEntities[ruleName]; ok {
		return rule.Salience, nil
	} else {
		return 0, errors.New(fmt.Sprintf("no such rules in pool:\"%s\"", ruleName))
	}
}

//get the rule's description
func (gp *GenginePool) GetRuleDesc(ruleName string) (string, error) {
	gp.updateLock.Lock()
	defer gp.updateLock.Unlock()

	if gp.clear || gp.ruleBuilder == nil {
		return "", errors.New("no rules in pool! ")
	}

	if rule, ok := gp.ruleBuilder.Kc.RuleEntities[ruleName]; ok {
		return rule.RuleDescription, nil
	} else {
		return "", errors.New(fmt.Sprintf("no such rules in pool:\"%s\"", ruleName))
	}
}

// count how many different rules in pool
func (gp *GenginePool) GetRulesNumber() int {
	gp.updateLock.Lock()
	defer gp.updateLock.Unlock()

	if gp.clear || gp.ruleBuilder == nil {
		return 0
	}
	return len(gp.ruleBuilder.Kc.RuleEntities)
}

func (gp *GenginePool) prepare(reqName string, req interface{}, respName string, resp interface{}) (*gengineWrapper, error) {
	//get gengine resource
	gw, e := gp.getGengine()
	if e != nil {
		return nil, e
	}

	gw.rulebuilder = gp.rbSlice[gw.tag]

	if reqName != "" && req != nil {
		gw.rulebuilder.Dc.Add(reqName, req)
	}

	if respName != "" && resp != nil {
		gw.rulebuilder.Dc.Add(respName, resp)
	}
	return gw, nil
}

func (gp *GenginePool) prepareWithMultiInput(data map[string]interface{}) (*gengineWrapper, error) {
	//get gengine resource
	gw, e := gp.getGengine()
	if e != nil {
		return nil, e
	}

	gw.rulebuilder = gp.rbSlice[gw.tag]

	for k, v := range data {
		//user should not inject "" string or nil value
		if k != "" && v != nil {
			gw.rulebuilder.Dc.Add(k, v)
		} else {
			log.Errorf("injected null string key or nil value! ")
		}
	}

	return gw, nil
}

//execute rules as the user set execute model when init or update
//req, it is better to be ptr, or you will not get changed data
//resp, it is better to be ptr, or you will not get changed data
// the return map[string]interface{} collection each rule returned result
func (gp *GenginePool) ExecuteRulesWithSpecifiedEM(reqName string, req interface{}, respName string, resp interface{}) (error, map[string]interface{}) {

	returnResultMap := make(map[string]interface{})
	//rules has bean cleared
	if gp.clear {
		//no data to execute rule
		return nil, returnResultMap
	}

	gw, e := gp.prepare(reqName, req, respName, resp)
	if e != nil {
		return e, returnResultMap
	}
	//release resource
	defer func() {
		gw.rulebuilder.Dc.Del(reqName, respName)
		gp.putGengineLocked(gw)
	}()

	if gp.execModel == SortModel { //sort
		// when some rule execute error ,it will continue to execute last
		e := gw.gengine.Execute(gw.rulebuilder, true)
		returnResultMap, _ = gw.gengine.GetRulesResultMap()
		return e, returnResultMap
	}

	if gp.execModel == ConcurrentModel { //concurrent
		e := gw.gengine.ExecuteConcurrent(gw.rulebuilder)
		returnResultMap, _ = gw.gengine.GetRulesResultMap()
		return e, returnResultMap
	}

	if gp.execModel == MixModel { //mix
		e := gw.gengine.ExecuteMixModel(gw.rulebuilder)
		returnResultMap, _ = gw.gengine.GetRulesResultMap()
		return e, returnResultMap
	}

	if gp.execModel == InverseMixModel { // inverse mix model
		e := gw.gengine.ExecuteInverseMixModel(gw.rulebuilder)
		returnResultMap, _ = gw.gengine.GetRulesResultMap()
		return e, returnResultMap
	}

	return nil, returnResultMap
}

/**
user can input more data to use in engine
it is no difference with ExecuteRules, you just can inject more data use this api

the return map[string]interface{} collection each rule returned result
*/
func (gp *GenginePool) ExecuteRulesWithMultiInputWithSpecifiedEM(data map[string]interface{}) (error, map[string]interface{}) {

	returnResultMap := make(map[string]interface{})
	//rules has bean cleared
	if gp.clear {
		//no data to execute rule
		return nil, returnResultMap
	}

	gw, e := gp.prepareWithMultiInput(data)
	if e != nil {
		return e, returnResultMap
	}
	//release resource
	defer func() {
		gw.clearInjected(getKeys(data)...)
		gp.putGengineLocked(gw)
	}()

	if gp.execModel == SortModel { //sort
		// when some rule execute error ,it will continue to execute last
		e := gw.gengine.Execute(gw.rulebuilder, true)
		returnResultMap, _ = gw.gengine.GetRulesResultMap()
		return e, returnResultMap
	}

	if gp.execModel == ConcurrentModel { //concurrent
		e := gw.gengine.ExecuteConcurrent(gw.rulebuilder)
		returnResultMap, _ = gw.gengine.GetRulesResultMap()
		return e, returnResultMap
	}

	if gp.execModel == MixModel { //mix
		e := gw.gengine.ExecuteMixModel(gw.rulebuilder)
		returnResultMap, _ = gw.gengine.GetRulesResultMap()
		return e, returnResultMap
	}

	if gp.execModel == InverseMixModel { // inverse mix model
		e := gw.gengine.ExecuteInverseMixModel(gw.rulebuilder)
		returnResultMap, _ = gw.gengine.GetRulesResultMap()
		return e, returnResultMap
	}

	return nil, returnResultMap

}

/***
this make user could use exemodel to control the select-exemodel

the return map[string]interface{} collection each rule returned result
*/
func (gp *GenginePool) ExecuteSelectedWithSpecifiedEM(data map[string]interface{}, names []string) (error, map[string]interface{}) {

	returnResultMap := make(map[string]interface{})
	//rules has bean cleared
	if gp.clear {
		//no data to execute rule
		return nil, returnResultMap
	}

	gw, e := gp.prepareWithMultiInput(data)
	if e != nil {
		return e, returnResultMap
	}
	//release resource
	defer func() {
		gw.clearInjected(getKeys(data)...)
		gp.putGengineLocked(gw)
	}()

	if gp.execModel == SortModel {
		e = gw.gengine.ExecuteSelectedRules(gw.rulebuilder, names)
		returnResultMap, _ = gw.gengine.GetRulesResultMap()
		return e, returnResultMap
	}

	if gp.execModel == ConcurrentModel {
		e = gw.gengine.ExecuteSelectedRulesConcurrent(gw.rulebuilder, names)
		returnResultMap, _ = gw.gengine.GetRulesResultMap()
		return e, returnResultMap
	}

	if gp.execModel == MixModel {
		e = gw.gengine.ExecuteSelectedRulesMixModel(gw.rulebuilder, names)
		returnResultMap, _ = gw.gengine.GetRulesResultMap()
		return e, returnResultMap
	}

	if gp.execModel == InverseMixModel {
		e = gw.gengine.ExecuteSelectedRulesInverseMixModel(gw.rulebuilder, names)
		returnResultMap, _ = gw.gengine.GetRulesResultMap()
		return e, returnResultMap
	}

	return nil, returnResultMap
}

// see gengine.go  Execute
func (gp *GenginePool) Execute(data map[string]interface{}, b bool) (error, map[string]interface{}) {
	returnResultMap := make(map[string]interface{})
	//rules has bean cleared
	if gp.clear {
		//no data to execute rule
		return nil, returnResultMap
	}

	gw, e := gp.prepareWithMultiInput(data)
	if e != nil {
		return e, returnResultMap
	}
	//release resource
	defer func() {
		gw.clearInjected(getKeys(data)...)
		gp.putGengineLocked(gw)
	}()

	e = gw.gengine.Execute(gw.rulebuilder, b)
	returnResultMap, _ = gw.gengine.GetRulesResultMap()
	return e, returnResultMap
}

// se gengine.go ExecuteWithStopTagDirect
func (gp *GenginePool) ExecuteWithStopTagDirect(data map[string]interface{}, b bool, sTag *Stag) (error, map[string]interface{}) {

	returnResultMap := make(map[string]interface{})
	//rules has bean cleared
	if gp.clear {
		//no data to execute rule
		return nil, returnResultMap
	}

	gw, e := gp.prepareWithMultiInput(data)
	if e != nil {
		return e, returnResultMap
	}
	//release resource
	defer func() {
		gw.clearInjected(getKeys(data)...)
		gp.putGengineLocked(gw)
	}()

	e = gw.gengine.ExecuteWithStopTagDirect(gw.rulebuilder, b, sTag)
	returnResultMap, _ = gw.gengine.GetRulesResultMap()
	return e, returnResultMap
}

//see gengine.go ExecuteConcurrent
func (gp *GenginePool) ExecuteConcurrent(data map[string]interface{}) (error, map[string]interface{}) {
	returnResultMap := make(map[string]interface{})
	//rules has bean cleared
	if gp.clear {
		//no data to execute rule
		return nil, returnResultMap
	}

	gw, e := gp.prepareWithMultiInput(data)
	if e != nil {
		return e, returnResultMap
	}
	//release resource
	defer func() {
		gw.clearInjected(getKeys(data)...)
		gp.putGengineLocked(gw)
	}()

	e = gw.gengine.ExecuteConcurrent(gw.rulebuilder)
	returnResultMap, _ = gw.gengine.GetRulesResultMap()
	return e, returnResultMap
}

// see gengine.go  ExecuteMixModel
func (gp *GenginePool) ExecuteMixModel(data map[string]interface{}) (error, map[string]interface{}) {
	returnResultMap := make(map[string]interface{})
	//rules has bean cleared
	if gp.clear {
		//no data to execute rule
		return nil, returnResultMap
	}

	gw, e := gp.prepareWithMultiInput(data)
	if e != nil {
		return e, returnResultMap
	}
	//release resource
	defer func() {
		gw.clearInjected(getKeys(data)...)
		gp.putGengineLocked(gw)
	}()

	e = gw.gengine.ExecuteMixModel(gw.rulebuilder)
	returnResultMap, _ = gw.gengine.GetRulesResultMap()
	return e, returnResultMap
}

//see gengine.go ExecuteMixModelWithStopTagDirect
func (gp *GenginePool) ExecuteMixModelWithStopTagDirect(data map[string]interface{}, sTag *Stag) (error, map[string]interface{}) {
	returnResultMap := make(map[string]interface{})
	//rules has bean cleared
	if gp.clear {
		//no data to execute rule
		return nil, returnResultMap
	}

	gw, e := gp.prepareWithMultiInput(data)
	if e != nil {
		return e, returnResultMap
	}
	//release resource
	defer func() {
		gw.clearInjected(getKeys(data)...)
		gp.putGengineLocked(gw)
	}()

	e = gw.gengine.ExecuteMixModelWithStopTagDirect(gw.rulebuilder, sTag)
	returnResultMap, _ = gw.gengine.GetRulesResultMap()
	return e, returnResultMap

}

// see gengine.go  ExecuteSelectedRules
func (gp *GenginePool) ExecuteSelectedRules(data map[string]interface{}, names []string) (error, map[string]interface{}) {
	returnResultMap := make(map[string]interface{})
	//rules has bean cleared
	if gp.clear {
		//no data to execute rule
		return nil, returnResultMap
	}

	gw, e := gp.prepareWithMultiInput(data)
	if e != nil {
		return e, returnResultMap
	}
	//release resource
	defer func() {
		gw.clearInjected(getKeys(data)...)
		gp.putGengineLocked(gw)
	}()

	e = gw.gengine.ExecuteSelectedRules(gw.rulebuilder, names)
	returnResultMap, _ = gw.gengine.GetRulesResultMap()
	return e, returnResultMap
}

//see gengine.go ExecuteSelectedRulesWithControl
func (gp *GenginePool) ExecuteSelectedRulesWithControl(data map[string]interface{}, b bool, names []string) (error, map[string]interface{}) {
	returnResultMap := make(map[string]interface{})
	//rules has bean cleared
	if gp.clear {
		//no data to execute rule
		return nil, returnResultMap
	}

	gw, e := gp.prepareWithMultiInput(data)
	if e != nil {
		return e, returnResultMap
	}
	//release resource
	defer func() {
		gw.clearInjected(getKeys(data)...)
		gp.putGengineLocked(gw)
	}()

	e = gw.gengine.ExecuteSelectedRulesWithControl(gw.rulebuilder, b, names)
	returnResultMap, _ = gw.gengine.GetRulesResultMap()
	return e, returnResultMap
}

//see gengine.go ExecuteSelectedRulesWithControlAsGivenSortedName
func (gp *GenginePool) ExecuteSelectedRulesWithControlAsGivenSortedName(data map[string]interface{}, b bool, sortedNames []string) (error, map[string]interface{}) {
	returnResultMap := make(map[string]interface{})
	//rules has bean cleared
	if gp.clear {
		//no data to execute rule
		return nil, returnResultMap
	}

	gw, e := gp.prepareWithMultiInput(data)
	if e != nil {
		return e, returnResultMap
	}
	//release resource
	defer func() {
		gw.clearInjected(getKeys(data)...)
		gp.putGengineLocked(gw)
	}()

	e = gw.gengine.ExecuteSelectedRulesWithControlAsGivenSortedName(gw.rulebuilder, b, sortedNames)
	returnResultMap, _ = gw.gengine.GetRulesResultMap()
	return e, returnResultMap
}

//see gengine.go ExecuteSelectedRulesWithControlAndStopTag
func (gp *GenginePool) ExecuteSelectedRulesWithControlAndStopTag(data map[string]interface{}, b bool, sTag *Stag, names []string) (error, map[string]interface{}) {
	returnResultMap := make(map[string]interface{})
	//rules has bean cleared
	if gp.clear {
		//no data to execute rule
		return nil, returnResultMap
	}

	gw, e := gp.prepareWithMultiInput(data)
	if e != nil {
		return e, returnResultMap
	}
	//release resource
	defer func() {
		gw.clearInjected(getKeys(data)...)
		gp.putGengineLocked(gw)
	}()

	e = gw.gengine.ExecuteSelectedRulesWithControlAndStopTag(gw.rulebuilder, b, sTag, names)
	returnResultMap, _ = gw.gengine.GetRulesResultMap()
	return e, returnResultMap
}

//see gengine.go ExecuteSelectedRulesWithControlAndStopTagAsGivenSortedName
func (gp *GenginePool) ExecuteSelectedRulesWithControlAndStopTagAsGivenSortedName(data map[string]interface{}, b bool, sTag *Stag, sortedNames []string) (error, map[string]interface{}) {
	returnResultMap := make(map[string]interface{})
	//rules has bean cleared
	if gp.clear {
		//no data to execute rule
		return nil, returnResultMap
	}

	gw, e := gp.prepareWithMultiInput(data)
	if e != nil {
		return e, returnResultMap
	}
	//release resource
	defer func() {
		gw.clearInjected(getKeys(data)...)
		gp.putGengineLocked(gw)
	}()

	e = gw.gengine.ExecuteSelectedRulesWithControlAndStopTagAsGivenSortedName(gw.rulebuilder, b, sTag, sortedNames)
	returnResultMap, _ = gw.gengine.GetRulesResultMap()
	return e, returnResultMap
}

//see gengine.go ExecuteSelectedRulesConcurrent
func (gp *GenginePool) ExecuteSelectedRulesConcurrent(data map[string]interface{}, names []string) (error, map[string]interface{}) {

	returnResultMap := make(map[string]interface{})
	//rules has bean cleared
	if gp.clear {
		//no data to execute rule
		return nil, returnResultMap
	}

	gw, e := gp.prepareWithMultiInput(data)
	if e != nil {
		return e, returnResultMap
	}
	//release resource
	defer func() {
		gw.clearInjected(getKeys(data)...)
		gp.putGengineLocked(gw)
	}()

	e = gw.gengine.ExecuteSelectedRulesConcurrent(gw.rulebuilder, names)
	returnResultMap, _ = gw.gengine.GetRulesResultMap()
	return e, returnResultMap
}

//see gengine.go ExecuteSelectedRulesMixModel
func (gp *GenginePool) ExecuteSelectedRulesMixModel(data map[string]interface{}, names []string) (error, map[string]interface{}) {

	returnResultMap := make(map[string]interface{})
	//rules has bean cleared
	if gp.clear {
		//no data to execute rule
		return nil, returnResultMap
	}

	gw, e := gp.prepareWithMultiInput(data)
	if e != nil {
		return e, returnResultMap
	}
	//release resource
	defer func() {
		gw.clearInjected(getKeys(data)...)
		gp.putGengineLocked(gw)
	}()

	e = gw.gengine.ExecuteSelectedRulesMixModel(gw.rulebuilder, names)
	returnResultMap, _ = gw.gengine.GetRulesResultMap()
	return e, returnResultMap

}

// //see gengine.go ExecuteInverseMixModel
func (gp *GenginePool) ExecuteInverseMixModel(data map[string]interface{}) (error, map[string]interface{}) {
	returnResultMap := make(map[string]interface{})
	//rules has bean cleared
	if gp.clear {
		//no data to execute rule
		return nil, returnResultMap
	}

	gw, e := gp.prepareWithMultiInput(data)
	if e != nil {
		return e, returnResultMap
	}
	//release resource
	defer func() {
		gw.clearInjected(getKeys(data)...)
		gp.putGengineLocked(gw)
	}()

	e = gw.gengine.ExecuteInverseMixModel(gw.rulebuilder)
	returnResultMap, _ = gw.gengine.GetRulesResultMap()
	return e, returnResultMap

}

//see gengine.go ExecuteSelectedRulesInverseMixModel
func (gp *GenginePool) ExecuteSelectedRulesInverseMixModel(data map[string]interface{}, names []string) (error, map[string]interface{}) {

	returnResultMap := make(map[string]interface{})
	//rules has bean cleared
	if gp.clear {
		//no data to execute rule
		return nil, returnResultMap
	}

	gw, e := gp.prepareWithMultiInput(data)
	if e != nil {
		return e, returnResultMap
	}
	//release resource
	defer func() {
		gw.clearInjected(getKeys(data)...)
		gp.putGengineLocked(gw)
	}()

	e = gw.gengine.ExecuteSelectedRulesInverseMixModel(gw.rulebuilder, names)
	returnResultMap, _ = gw.gengine.GetRulesResultMap()
	return e, returnResultMap
}

// see gengine.go  ExecuteNSortMConcurrent
func (gp *GenginePool) ExecuteNSortMConcurrent(nSort, mConcurrent int, b bool, data map[string]interface{}) (error, map[string]interface{}) {

	returnResultMap := make(map[string]interface{})
	//rules has bean cleared
	if gp.clear {
		//no data to execute rule
		return nil, returnResultMap
	}

	gw, e := gp.prepareWithMultiInput(data)
	if e != nil {
		return e, returnResultMap
	}
	//release resource
	defer func() {
		gw.clearInjected(getKeys(data)...)
		gp.putGengineLocked(gw)
	}()

	e = gw.gengine.ExecuteNSortMConcurrent(nSort, mConcurrent, gw.rulebuilder, b)
	returnResultMap, _ = gw.gengine.GetRulesResultMap()
	return e, returnResultMap
}

// see gengine.go ExecuteNConcurrentMSort
func (gp *GenginePool) ExecuteNConcurrentMSort(nSort, mConcurrent int, b bool, data map[string]interface{}) (error, map[string]interface{}) {
	returnResultMap := make(map[string]interface{})
	//rules has bean cleared
	if gp.clear {
		//no data to execute rule
		return nil, returnResultMap
	}

	gw, e := gp.prepareWithMultiInput(data)
	if e != nil {
		return e, returnResultMap
	}
	//release resource
	defer func() {
		gw.clearInjected(getKeys(data)...)
		gp.putGengineLocked(gw)
	}()

	e = gw.gengine.ExecuteNConcurrentMSort(nSort, mConcurrent, gw.rulebuilder, b)
	returnResultMap, _ = gw.gengine.GetRulesResultMap()
	return e, returnResultMap
}

//see gengine.go ExecuteNConcurrentMConcurrent
func (gp *GenginePool) ExecuteNConcurrentMConcurrent(nSort, mConcurrent int, b bool, data map[string]interface{}) (error, map[string]interface{}) {
	returnResultMap := make(map[string]interface{})
	//rules has bean cleared
	if gp.clear {
		//no data to execute rule
		return nil, returnResultMap
	}

	gw, e := gp.prepareWithMultiInput(data)
	if e != nil {
		return e, returnResultMap
	}
	//release resource
	defer func() {
		gw.clearInjected(getKeys(data)...)
		gp.putGengineLocked(gw)
	}()

	e = gw.gengine.ExecuteNConcurrentMConcurrent(nSort, mConcurrent, gw.rulebuilder, b)
	returnResultMap, _ = gw.gengine.GetRulesResultMap()
	return e, returnResultMap

}

// see  gengine.go ExecuteSelectedNSortMConcurrent
func (gp *GenginePool) ExecuteSelectedNSortMConcurrent(nSort, mConcurrent int, b bool, names []string, data map[string]interface{}) (error, map[string]interface{}) {
	returnResultMap := make(map[string]interface{})
	//rules has bean cleared
	if gp.clear {
		//no data to execute rule
		return nil, returnResultMap
	}

	gw, e := gp.prepareWithMultiInput(data)
	if e != nil {
		return e, returnResultMap
	}
	//release resource
	defer func() {
		gw.clearInjected(getKeys(data)...)
		gp.putGengineLocked(gw)
	}()

	e = gw.gengine.ExecuteSelectedNSortMConcurrent(nSort, mConcurrent, gw.rulebuilder, b, names)
	returnResultMap, _ = gw.gengine.GetRulesResultMap()
	return e, returnResultMap
}

// see gengine.go ExecuteSelectedNConcurrentMSort
func (gp *GenginePool) ExecuteSelectedNConcurrentMSort(nSort, mConcurrent int, b bool, names []string, data map[string]interface{}) (error, map[string]interface{}) {

	returnResultMap := make(map[string]interface{})
	//rules has bean cleared
	if gp.clear {
		//no data to execute rule
		return nil, returnResultMap
	}

	gw, e := gp.prepareWithMultiInput(data)
	if e != nil {
		return e, returnResultMap
	}
	//release resource
	defer func() {
		gw.clearInjected(getKeys(data)...)
		gp.putGengineLocked(gw)
	}()

	e = gw.gengine.ExecuteSelectedNConcurrentMSort(nSort, mConcurrent, gw.rulebuilder, b, names)
	returnResultMap, _ = gw.gengine.GetRulesResultMap()
	return e, returnResultMap
}

// see gengine.go ExecuteSelectedNConcurrentMConcurrent
func (gp *GenginePool) ExecuteSelectedNConcurrentMConcurrent(nSort, mConcurrent int, b bool, names []string, data map[string]interface{}) (error, map[string]interface{}) {

	returnResultMap := make(map[string]interface{})
	//rules has bean cleared
	if gp.clear {
		//no data to execute rule
		return nil, returnResultMap
	}

	gw, e := gp.prepareWithMultiInput(data)
	if e != nil {
		return e, returnResultMap
	}
	//release resource
	defer func() {
		gw.clearInjected(getKeys(data)...)
		gp.putGengineLocked(gw)
	}()

	e = gw.gengine.ExecuteSelectedNConcurrentMConcurrent(nSort, mConcurrent, gw.rulebuilder, b, names)
	returnResultMap, _ = gw.gengine.GetRulesResultMap()
	return e, returnResultMap
}

// see gengine.go ExecuteDAGModel
func (gp *GenginePool) ExecuteDAGModel(dag [][]string, data map[string]interface{}) (error, map[string]interface{}) {

	returnResultMap := make(map[string]interface{})
	//rules has bean cleared
	if gp.clear {
		//no data to execute rule
		return nil, returnResultMap
	}

	gw, e := gp.prepareWithMultiInput(data)
	if e != nil {
		return e, returnResultMap
	}
	//release resource
	defer func() {
		gw.clearInjected(getKeys(data)...)
		gp.putGengineLocked(gw)
	}()

	e = gw.gengine.ExecuteDAGModel(gw.rulebuilder, dag)
	returnResultMap, _ = gw.gengine.GetRulesResultMap()
	return e, returnResultMap
}

func getKeys(data map[string]interface{}) []string {
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	return keys
}
