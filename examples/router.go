package main

import (
	"errors"
	"fmt"
	. "github.com/metakeule/grack"
	. "github.com/metakeule/grack/base"
	. "github.com/metakeule/grack/routing"
)

func Panic(c Contexter, s string) {
	Error(c, errors.New(s))
}

type MyContext struct {
	*IO
	hasLayout bool
	layout    string
}

func SetLayout(c Contexter, s string) {
	d := c.(*MyContext)
	d.hasLayout = true
	d.layout = s
}

func Layout(c Contexter) (s string) {
	d := c.(*MyContext)
	if d.hasLayout {
		s = d.layout
		return
	}
	s = ""
	return
}

func NewContext(in string) Contexter {
	c := &MyContext{IO: NewIO().(*IO)}
	SetIn(c, in)
	return c
}

func mw1(c Contexter, err error) {
	fmt.Println("mw1 called")
	in := InString(c)
	if in == "panic" {
		Panic(c, "in the streets of london")
	}
	SetOut(c, "Hello "+in)
	if in == "router" {
		JumpToRouter(c)
	}
	Next(c)
}

func mw1a(c Contexter, err error) {
	fmt.Println("mw1a called")
	in := InString(c)
	if in == "skip" {
		Skip(c, 2)
	}
	Next(c)
}

func mw2(c Contexter, err error) {
	fmt.Println("mw2 called")
	SetOut(c, OutString(c)+", welcome!")
	Next(c)
}

func router(c Contexter, err error) {
	fmt.Println("router called")
	SetLayout(c, "here the layout")
	Next(c)
}

func response(c Contexter, err error) {
	fmt.Printf("Layout: %#v\n", Layout(c))
	if err != nil {
		fmt.Printf("unhandled panic: %s\n", err)
		return
	}
	fmt.Println(Out(c))
}

// var MainRack = New(NewRack())
var MainRack = NewRoutingRack(NewRack())

func init() {
	Push(MainRack, mw1)
	Push(MainRack, mw1a)
	Router(MainRack, router)
	Push(MainRack, mw2)
	SetResponder(MainRack, response)
}

func main() {
	c := NewContext("World")
	Run(MainRack, c)
	fmt.Println("-----------")
	c = NewContext("panic")
	Run(MainRack, c)
	fmt.Println("-----------")
	c = NewContext("router")
	Run(MainRack, c)
	fmt.Println("-----------")
	c = NewContext("skip")
	Run(MainRack, c)
}
