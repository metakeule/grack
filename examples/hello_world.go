package main

import (
	"fmt"
	. "github.com/metakeule/grack"
	. "github.com/metakeule/grack/base"
)

func hello(c Contexter, err error) {
	SetIn(c, "Hello")
	Next(c)
}

func world(c Contexter, err error) {
	SetOut(c, InString(c)+" world!")
}

func print(c Contexter, err error) {
	fmt.Println(Out(c))
}

var rack = NewRack()

func init() {
	rack.Push(hello)
	rack.Push(world)
	rack.SetResponder(print)
}

func main() {
	Run(rack, NewIO())
}
