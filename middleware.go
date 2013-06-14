package grack

import (
	"net/http"
)

func MiddlewareFunc(fn func(Racker)) *middlewareFunc {
	name := GetFunctionName(fn)
	return &middlewareFunc{fn, name}
}

type middlewareFunc struct {
	fn   func(Racker)
	name string
}

func (ø middlewareFunc) Call(r Racker) { ø.fn(r) }
func (ø middlewareFunc) Name() string  { return ø.name }

type Middleware interface {
	Call(Racker)
	Name() string
}

func MiddlewareHandler(h http.Handler, name string) *middlewareFunc {
	fn := func(r Racker) {
		h.ServeHTTP(r, r.Request())
	}
	return &middlewareFunc{fn, name}
}
