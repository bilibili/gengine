package icore

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"
)

type Animal struct {
	Name string
}

func copyRuleBuilder(src, dst interface{}) error {
	buf := new(bytes.Buffer)
	if e := gob.NewEncoder(buf).Encode(src); e != nil {
		return e
	}
	return gob.NewDecoder(buf).Decode(dst)
}

func Clone(a, b interface{}) error {

	buff := new(bytes.Buffer)
	e := gob.NewEncoder(buff).Encode(a)
	e = gob.NewDecoder(buff).Decode(b)
	return e
}

func Test_copy(t *testing.T) {

	animal := &Animal{Name: "cat"}
	var B Animal

	e := Clone(animal, &B)
	if e != nil {
		panic(e)
	}

	println(B.Name)
	B.Name = "dog"
	println(animal.Name)
	println(B.Name)
}

func Test_1_copy(t *testing.T) {
	//gw := gp.freeGengines[0]
	//copy(gp.freeGengines, gp.freeGengines[:1])
	//gp.freeGengines = gp.freeGengines[:numFree-1]

	a1 := []int16{1, 2}
	x := a1[0]
	//copy(a1, a1[:1])
	a1 = a1[1:]
	println(fmt.Sprintf("%d,%+v", x, a1))
	//var s []int
	//println(len(s))
	//s = append(s, 1)
	//println(len(s))

}
