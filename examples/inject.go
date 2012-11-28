package main

import (
	"fmt"
	. "github.com/metakeule/grack"
	. "github.com/metakeule/grack/base"
)

func hello(outer Contexter, err_outer error) {
	SetOut(outer, "Hello")
	io := NewIO()
	SetOut(io, Out(outer))
	Inject(io, attr_rack, func(inner Contexter, err_inner error) {
		SetOut(outer, OutString(inner))
	})
	Next(outer)
}

func strange(c Contexter, err error) {
	SetOut(c, OutString(c)+" strange")
	Next(c)
}

func lovely(c Contexter, err error) {
	SetOut(c, OutString(c)+" and lovely")
	Next(c)
}

func world(c Contexter, err error) {
	SetOut(c, OutString(c)+" world!")
}

func print(c Contexter, err error) {
	fmt.Println(OutString(c))
}

var rack = NewRack()
var attr_rack = NewRack()

func init() {
	rack.Push(hello)
	rack.Push(world)
	rack.SetResponder(print)

	attr_rack.Push(strange)
	attr_rack.Push(lovely)
}

func main() {
	Run(rack, NewIO())
}
