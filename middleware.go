package grack

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
