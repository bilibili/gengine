package pool

import (
	"fmt"
	"github.com/bilibili/gengine/engine"
	"testing"
)

const pool_sorted_name_rule string = `	
	rule "1"
	begin
		println("----1-------")
	end
	
	rule "2"
	begin
		println("----2-------")
	end
	
	rule "3"
	begin
		println("----3-------")
	end
	
	rule "4"
	begin
		println("----4-------")
	end
`

func Test_pool_sorted_name(t *testing.T) {
	apis := make(map[string]interface{})
	apis["println"] = fmt.Println
	pool, err := engine.NewGenginePool(1, 2, 1, pool_sorted_name_rule, apis)
	if err != nil {
		panic(err)
	}

	data := make(map[string]interface{})
	e1, _ := pool.ExecuteSelectedRulesWithControlAsGivenSortedName(data, true, []string{"1", "3", "4", "2"})
	if e1 != nil {
		panic(e1)
	}

	println("==================1==================")

	sTag := &engine.Stag{StopTag: false}
	e2, _ := pool.ExecuteSelectedRulesWithControlAndStopTagAsGivenSortedName(data, true, sTag, []string{"2", "1", "4", "3"})
	if e2 != nil {
		panic(e2)
	}

	println("==================2==================")
}
