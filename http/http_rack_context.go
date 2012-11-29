package grack_http

import . "github.com/metakeule/grack"

/*
	If you create an own HTTPRacker and want to enforce the type of the middleware
	you will have to overwrite some methods (due to the lack of generics in go).

	Let's say your context struct is MyContext and your rack MyRack.

		import (
			h "github.com/metakeule/grack"
		)

		type MyContext struct {
				h.HTTPContexter
			  *h.HTTPContext
			  // other things
		}

		type MyRack struct {
				h.HTTPRacker
			  *h.HTTPRack
			  // other things
		}

	Then you might basically make a copy of this file and do the following:

	- change
			import . "github.com/metakeule/grack"
		to
			import h "github.com/metakeule/grack"

	- change
			type Ctx HTTPContexter
		to
			type Ctx *MyContext

	- modify (r *HTTPRack) NewContext() (c HTTPContexter) to
			func (r *MyRack) NewContext() (c h.HTTPContexter) {
				// do further initialization here
				return &MyContext{HTTPContext: h.NewContext()}
			}

	- replace *HTTPRack with *MyRack for all the other methods in this file
*/
type Ctx HTTPContexter

type middlewareCaller struct {
	MiddlewareCaller
	Middleware func(Ctx, error)
}

func (c *middlewareCaller) Call(ctx Contexter, err error) {
	c.Middleware(ctx.(Ctx), err)
}

func (r *HTTPRack) NewContext() (c HTTPContexter) {
	return NewHTTPContext()
}

func (r *HTTPRack) Push(i interface{}) {
	mw := i.(func(Ctx, error))
	r.Rack.Stack = append(r.Rack.Stack, &middlewareCaller{Middleware: mw})
}

func (r *HTTPRack) SetResponder(i interface{}) {
	mw := i.(func(Ctx, error))
	r.Rack.Responder = &middlewareCaller{Middleware: mw}
}
