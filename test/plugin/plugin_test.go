package main

import (
	"fmt"
	"github.com/bilibili/gengine/builder"
	"github.com/bilibili/gengine/context"
	"github.com/bilibili/gengine/engine"
	"os"
	"path"
	"path/filepath"
	"plugin"
	"testing"
)


func Test_pligin(t *testing.T) {

	dir, err := os.Getwd()
	if err!=nil {
		panic(err)
	}

	// load module 插件您也可以使用go http.Request从远程下载到本地,在加载做到动态的执行不同的功能
	// 1. open the so file to load the symbols
	plug, err := plugin.Open(dir + "/plugin_M_m.so")
	if err != nil {
		panic(err)
	}
	println("plugin opened")

	// 2. look up a symbol (an exported function or variable)
	// in this case, variable Greeter
	m, err := plug.Lookup("M") //大写
	if err != nil {
		panic(err)
	}

	// 3. Assert that loaded symbol is of a desired type
	man, ok := m.(Man)
	if !ok {
		fmt.Println("unexpected type from module symbol")
		os.Exit(1)
	}

	// 4. use the module
	if err := man.SaveLive(); err != nil {
		println("use plugin man failed, ", err)
	}

}

func Test_plugin_with_gengine(t *testing.T)  {

	dir, err := os.Getwd()
	if err!=nil {
		panic(err)
	}

	dc := context.NewDataContext()
	//3.load plugin into apiName, exportApi
	_, _, e := dc.PluginLoader( dir + "/plugin_M_m.so")
	if e != nil {
		panic(e)
	}

	dc.Add("println", fmt.Println)
	ruleBuilder := builder.NewRuleBuilder(dc)
	err = ruleBuilder.BuildRuleFromString(`
	rule "1"
	begin
	 
	//this method is defined in plugin
	err = m.SaveLive()

	if isNil(err) {
	   println("err is nil")
	}
	end
	`)

	if err != nil {
		panic(err)
	}
	gengine := engine.NewGengine()
	err = gengine.Execute(ruleBuilder, false)

	if err!=nil {
		panic(err)
	}
}


func Test_plugin_with_pool(t *testing.T)  {

	rule :=`
	rule "1"
	begin
	 
	//this method is defined in plugin
	err = m.SaveLive()

	if isNil(err) {
	   println("err is nil")
	}
	end`

	apis := make(map[string]interface{})
	apis["println"] = fmt.Println
	pool, e := engine.NewGenginePool(1, 2, 1, rule, apis)
	if e != nil {
		panic(e)
	}

	dir, err := os.Getwd()
	if err!=nil {
		panic(err)
	}

	e = pool.PluginLoader( dir + "/plugin_M_m.so")
	if e != nil {
		panic(e)
	}
	data := make(map[string]interface{})
	e, _ = pool.Execute(data, true)
	if e != nil {
		panic(e)
	}

	//twice execute
	e, _ = pool.Execute(data, true)
	if e != nil {
		panic(e)
	}
}

func Test_file(t *testing.T) {

	///Users/renyunyi/go_project/gengine/test/plugin_t/
	files :="plugin_superman.so"
	dir, file := filepath.Split(files)
	println(dir, file, filepath.Base(files), path.Ext(files))
	s, err := os.Getwd()
	if err!=nil {
		panic(err)
	}

	println(s)
}
