package grack

// import "fmt"

//type Middleware func(Contexter, error)
type MiddlewareCaller interface {
	Call(Contexter, error)
}

type middlewareCaller struct {
	MiddlewareCaller
	Middleware func(Contexter, error)
	Context    Contexter
}

func (c *middlewareCaller) Call(ctx Contexter, err error) {
	c.Middleware(ctx, err)
}

func NewMwCaller(mw func(Contexter, error)) *middlewareCaller {
	return &middlewareCaller{Middleware: mw}
}

/*
	A simple and general implementation of a rack
*/
type Rack struct {
	Racker
	Created   *debugCall
	Stack     []MiddlewareCaller
	Responder MiddlewareCaller
}

func (r *Rack) Push(i interface{}) {
	mw := NewMwCaller(i.(func(Contexter, error)))
	r.Stack = append(r.Stack, mw)
}

func (r *Rack) SetResponder(i interface{}) {
	mw := NewMwCaller(i.(func(Contexter, error)))
	r.Responder = mw
}

func (r *Rack) Clone() Racker {
	return &Rack{Stack: r.Stack, Responder: r.Responder}
}

func (r *Rack) Rack_() *Rack {
	return r
}

// returns the middleware at position p
func Get(r Racker, p uint) MiddlewareCaller {
	return r.Rack_().Stack[p]
}

// the length of the rack (number of middlewares)
func Len(r Racker) uint {
	return uint(len(r.Rack_().Stack))
}

// adds middleware to the end of the stack
func Push(r Racker, i interface{}) {
	r.Push(i)
}

// set the middleware to respond with
func SetResponder(r Racker, i interface{}) {
	r.SetResponder(i)
}

// adds middleware to the end of the stack
func CloneRack(r Racker) Racker {
	return r.Clone()
}

// Racker represents a rack of middlewares
type Racker interface {
	Rack_() *Rack
	// we need to be able to overwrite them,
	// even if they are predefined in *Rack
	Push(interface{})
	SetResponder(interface{})
	Clone() Racker
}

// a contexter represents the context through the stack run of a racker
// all processing will be done from the outside via setting and getting of things
type Contexter interface {
	Ctx() *ContextData
}

type ContextData struct {
	Rack     Racker
	Position uint
	Stopped  bool
}
