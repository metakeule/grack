package grack_http

import (
	. "github.com/metakeule/grack"
	"net/http"
)

/*
	To make your own HTTPContexter, simply "inherit" like this:

		import (
			. "github.com/metakeule/grack"
			. "github.com/metakeule/grack/http"
		)

		type MyContext struct {
			*HTTPContext
			// other things
		}

		func main {
			my := &MyContext{HTTPContext: NewContext()}
		}

	If you want to access properties of the HTTPContext, use HTTPContext_():

		func SetOKStatus(my *MyContext) {
			my.HTTPContext_().Status = 200
		}

	All methods are of cause inherited to the toplevel, so that

		func Pos(my *MyContext) uint {
			return my.Ctx().Position
		}

*/

type HTTPContext struct {
	*ContextData
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Status         int
	Out            string
}

type HTTPContexter interface {
	Contexter
	HTTPContext_() *HTTPContext
	Rack() Racker
}

func (c *HTTPContext) HTTPContext_() *HTTPContext {
	return c
}

func (c *HTTPContext) Ctx() *ContextData {
	return c.ContextData
}

func (c *HTTPContext) Rack() (r Racker) {
	return c.ContextData.Rack
}

func NewContext() *HTTPContext {
	return &HTTPContext{ContextData: &ContextData{}}
}
