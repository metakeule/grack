package main

import (
	"fmt"
	"github.com/metakeule/grack"
	"net/http"
)

func hello(r grack.Racker) {
	r.Set("out", "Hello")
	r.Next()
}

func world(r grack.Racker) {
	r.TextString(r.GetString("out") + " world!\nYour Path is: " + r.Request().URL.RequestURI())
}

var rack = grack.NewRack()

func init() {
	rack.PushFunc(hello)
	rack.PushFunc(world)
}

func main() {
	fmt.Println("look at http://localhost:8080/strange")
	http.ListenAndServe(":8080", rack)
}
