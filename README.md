grack - racks for go
====================

grack is inspired by Christian Neukirchens rack (for Ruby).

It offers a way to organize middleware.

You may build your own customized racks on the basic implementation.

simple example
--------------

```go
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
```

more examples to come
---------------------

examples for routing, injection, delegation, error handling and custom racks will be found in the examples directory.

Thanks
------

* Christian Neukirchen for rack
* the go authors for go