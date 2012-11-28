package main

import (
	"fmt"
	. "github.com/metakeule/grack"
	. "github.com/metakeule/grack/base"
	"io/ioutil"
)

func hello(c Contexter, err error) {
	SetOut(c, "Hello")
	_, e := ioutil.ReadFile("c:////invalid..c://invalid")
	if e != nil {
		Error(c, e) // direct return error to responder, skipping further middleware
	}
}

func world(c Contexter, err error) {
	SetOut(c, OutString(c)+" world!")
}

func print(c Contexter, err error) {
	if err != nil {
		fmt.Printf("ERROR %s\n", err)
	} else {
		fmt.Println(Out(c))
	}
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
