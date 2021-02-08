//package必须是main
package main


type Man interface {
	SaveLive() error
}

type SuperMan struct {

}

func (g *SuperMan) SaveLive() error {

	println("execute finished...")
	return nil
}

//go build -buildmode=plugin -o=plugin_M_m.so plugin_superman.go
// exported as symbol named "M",必须大写开头
var M = SuperMan{}