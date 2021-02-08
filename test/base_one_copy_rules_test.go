package test

import (
	"testing"
	"time"
)

type oneCopy struct {
}

//此测试为了证明在不改变struct自身的数据情况下，同时传入多个参数到结构体所指的方法上去执行时，这些参数是不会相互覆盖的
func Test_base_one_copy_rules(t *testing.T) {

	m1 := make(map[string]string)
	m1["a1"] = "b1"
	m1["a2"] = "b2"
	m1["a3"] = "b3"

	m2 := make(map[string]string)
	m2["a11"] = "b11"
	m2["a22"] = "b22"
	m2["a33"] = "b33"

	m3 := make(map[string]string)
	m3["a1x"] = "b1x"
	m3["a2x"] = "b2x"
	m3["a3x"] = "b3x"

	op := &oneCopy{}

	go func() {
		for i := 0; i < 100000; i++ {
			op.change(&m1)
		}
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			op.change(&m2)
		}
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			op.change(&m3)
		}
	}()

	time.Sleep(2 * time.Second)

	for k, v := range m1 {
		println("m1-->", k, v)
	}

	for k, v := range m2 {
		println("m2-->", k, v)
	}

	for k, v := range m3 {
		println("m3-->", k, v)
	}

}

func (op *oneCopy) change(m *map[string]string) {
	(*m)["a"] = "b"
	(*m)["a1"] = "b1000"
}
