package main

import (
	"fmt"
	. "github.com/metakeule/grack"
	. "github.com/metakeule/grack/http"
	"net/http"
)

func hello(c Ctx, err error) {
	SetOut(c, "Hello")
	Next(c)
}

func world(c Ctx, err error) {
	SetOut(c, Out(c)+" world!\nYour Path is: "+Path(c))
}

var rack = NewRack()

func init() {
	rack.Push(hello)
	rack.Push(world)
}

func main() {
	fmt.Println("look at http://localhost:8080/strange")
	http.ListenAndServe(":8080", NewServer(rack))
}
