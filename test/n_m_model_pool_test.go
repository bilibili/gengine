package test

import (
	"fmt"
	"github.com/bilibili/gengine/engine"
	"testing"
)

//测试代码框架
func nmpFramework(n, m, em int, names []string) {
	type Data struct {
		Count int
	}

	apis := make(map[string]interface{})
	apis["println"] = fmt.Println
	//无状态api在这里注入
	pool, e := engine.NewGenginePool(2, 4, 1, n_m_model_rules, apis)
	if e != nil {
		panic(e)
	}

	data := make(map[string]interface{})
	req := &Data{Count: 0}
	data["data"] = req

	var e2 error
	var resultMap map[string]interface{}
	if em == 1 {
		//执行时注入与本次请求强相关的依赖数据或api
		e2, resultMap = pool.ExecuteNSortMConcurrent(n, m, true, data)
	}
	if em == 2 {
		e2, resultMap = pool.ExecuteNConcurrentMSort(n, m, true, data)
	}
	if em == 3 {
		e2, resultMap = pool.ExecuteNConcurrentMConcurrent(n, m, true, data)
	}
	if em == 4 {
		e2, resultMap = pool.ExecuteSelectedNSortMConcurrent(n, m, true, names, data)
	}
	if em == 5 {
		e2, resultMap = pool.ExecuteSelectedNConcurrentMSort(n, m, true, names, data)
	}
	if em == 6 {
		e2, resultMap = pool.ExecuteSelectedNConcurrentMConcurrent(n, m, true, names, data)
	}

	if e2 != nil {
		panic(e2)
	}

	//如果有规则有返回值, 可以这样获取
	i := resultMap["77"] //获取规则名为77的返回值
	s := i.(string)
	println("get return -->", s)
}

func Test_pool_n_sort_m_concurrent(t *testing.T) {
	//ExecuteNSortMConcurrent
	nmpFramework(3, 3, 1, []string{})
}

func Test_pool_n_concurrent_m_sort(t *testing.T) {
	//ExecuteNConcurrentMSort
	nmpFramework(3, 3, 2, []string{})
}

func Test_pool_n_concurrent_m_concurrent(t *testing.T) {
	//ExecuteNConcurrentMConcurrent
	nmpFramework(3, 3, 3, []string{})
}

func Test_pool_selected_n_sort_m_concurrent(t *testing.T) {
	//ExecuteSelectedNSortMConcurrent
	nmpFramework(3, 2, 4, []string{"100", "98", "89", "90", "77"})
}

func Test_pool_selected_n_concurrent_m_sort(t *testing.T) {
	//ExecuteSelectedNConcurrentMSort
	nmpFramework(3, 3, 5, []string{"100", "98", "89", "90", "77", "50"})
}

func Test_pool_selected_n_concurrent_m_concurrent(t *testing.T) {
	//ExecuteSelectedNConcurrentMConcurrent
	nmpFramework(3, 3, 6, []string{"100", "98", "89", "90", "77", "50"})
}
