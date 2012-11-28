package grack_http

import (
	"fmt"
	. "github.com/metakeule/grack"
	R "github.com/metakeule/grack/base"
	"log"
)

/*
	To make your own HTTPRacker, simply "inherit" like this:

		import (
			. "github.com/metakeule/grack"
			. "github.com/metakeule/grack/http"
		)

		type MyRack struct {
			*HTTPRack
			// something other
		}

		func main {
			my := &MyRack{HTTPRack: NewRack()}
		}

	If you want to access properties of the HTTPRack, use HTTPRack_():

		func Rack(my *MyRack) Racker {
			return my.HTTPRack_().Rack
		}

	All methods are of cause inherited to the toplevel, so that

		func AddLogger(my *MyRack, logger func(HTTPContexter, error)) {
			my.Push(logger)
		}

	If you want to force the middleware to accept a certain context, you will have to overwrite
	some methods, see the file

		http_rack_context.go

	for more information
*/

type HTTPRack struct {
	*Rack
}

type HTTPRacker interface {
	Racker
	HTTPRack_() *HTTPRack
	NewContext() HTTPContexter
}

func (h *HTTPRack) HTTPRack_() *HTTPRack {
	return h
}

func (h *HTTPRack) Clone() Racker {
	r := h.Rack.Clone()
	return &HTTPRack{Rack: r.(*Rack)}
}

func HTTPResponder(c Ctx, err error) {
	hc := c.HTTPContext_()
	if err != nil {
		hc.Status = 500
		s := fmt.Sprintf("Error: %#v", err)
		log.Fatal(s)
		hc.Out = s
	}
	Flush(hc)
}

func NewRack() *HTTPRack {
	rck := &HTTPRack{Rack: R.NewRack()}
	rck.SetResponder(HTTPResponder)
	return rck
}
