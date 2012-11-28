grack - racks for go
====================

grack is inspired by Christian Neukirchens rack (for Ruby).

It offers a general way to stack middleware and is not focussed on
web stacks.

It has the basic infrastructure for your own racks,
middlewares and routers.

simple example
--------------

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


simple web example
------------------

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


more examples
-------------

examples for routing, injection, delegation, error handling and custom racks and contexts could be found in the examples directory.

Thanks
------

* Christian Neukirchen for rack
* the go authors for go