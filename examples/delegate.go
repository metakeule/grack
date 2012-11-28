package main

import (
	"fmt"
	. "github.com/metakeule/grack"
	. "github.com/metakeule/grack/base"
)

func hello(c Contexter, err error) {
	SetIn(c, "Hello")
	Delegate(c, universe_rack)
}

func universe(c Contexter, err error) {
	SetOut(c, InString(c)+" universe!")
}

func universe_print(c Contexter, err error) {
	fmt.Println("ended up in delegate: " + OutString(c))
}

func print(c Contexter, err error) {
	fmt.Println("should never be here: " + OutString(c))
}

var rack = NewRack()
var universe_rack = NewRack()

func init() {
	rack.Push(hello)
	rack.SetResponder(print)

	universe_rack.Push(universe)
	universe_rack.SetResponder(universe_print)
}

func main() {
	Run(rack, NewIO())
}
