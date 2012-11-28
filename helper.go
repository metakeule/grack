package grack

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

type debugCall struct {
	Function string
	File     string
	Line     uint
}

func DebugCall() *debugCall {
	pc, file, line, _ := runtime.Caller(3)
	n := runtime.FuncForPC(pc).Name()
	a := strings.Split(n, "/")
	name := a[len(a)-1]
	return &debugCall{Function: name, File: file, Line: uint(line)}
}

func NameOf(i interface{}) string {
	return reflect.TypeOf(i).Elem().Name()
}

func RackName(c Contexter) string {
	return NameOf(c.Ctx().Rack)
}

func RackCreated(c Contexter) *debugCall {
	return c.Ctx().Rack.Rack_().Created
}

func DebugRack(c Contexter) (r map[string]string) {
	r = map[string]string{}
	ra := c.Ctx().Rack
	creation := ra.Rack_().Created
	r["Name"] = NameOf(ra)
	r["CreatedByFunction"] = creation.Function
	r["CreatedInFile"] = creation.File
	r["CreatedInLine"] = fmt.Sprintf("%d", creation.Line)
	// r["Inspect"] = fmt.Sprintf("%s", ra)
	r["Len"] = fmt.Sprintf("%d", Len(ra))
	return
}

func ContextName(c Contexter) string {
	return NameOf(c)
}

func Reset(c Contexter) {
	d := c.Ctx()
	d.Position = uint(0)
	d.Stopped = false
}

// call the middleware at the current position
func Call(c Contexter) {
	d := c.Ctx()
	if d.Stopped == false {
		r := d.Rack.Rack_()
		p := d.Position
		// fmt.Printf("pos: %d, len: %d", p, r.Len())
		if l := Len(r); l <= p {
			Return(c)
			return
		}
		//r.Get(p)(c, nil)
		Get(r, p).Call(c, nil)
	}
}

// return immediately to the outputter
func Return(c Contexter) {
	d := c.Ctx()
	if d.Stopped == false {
		d.Stopped = true
		res := d.Rack.Rack_().Responder
		if res != nil {
			res.Call(c, nil)
		}
	}
}

// set an error and report it immediately to the outputter
func Error(c Contexter, e error) {
	d := c.Ctx()
	if d.Stopped == false {
		d.Stopped = true
		res := d.Rack.Rack_().Responder
		if res != nil {
			res.Call(c, e)
		}
	}
}

func CallAndReturn(c Contexter) {
	Call(c)
	if c.Ctx().Stopped == false {
		Return(c)
	}
}

// increments the position of context if its rack has
// enough entries. returns true on success, otherwise false
func incrPosition(c Contexter) (ok bool) {
	d := c.Ctx()
	p := d.Position + 1
	r := d.Rack.Rack_()
	ok = false
	if Len(r) > p {
		d.Position = p
		ok = true
	}
	return
}

// call the next middleware and return to the outputter
// if there is no middleware left
func Next(c Contexter) {
	if ok := incrPosition(c); ok != true {
		// fmt.Println("can not increment")
		Return(c)
		return
	}
	Call(c)
}

func Run(r Racker, c Contexter) {
	c.Ctx().Rack = r
	Reset(c)
	CallAndReturn(c)
}

// goto a position p within the middleware stack, handle with care
func GoTo(c Contexter, p uint) {
	c.Ctx().Position = p
	Call(c)
}

// skip n middlewares
func Skip(c Contexter, n uint) {
	GoTo(c, c.Ctx().Position+n+1)
}

// delegates to a racker. the context then is gone and handled by the racker
func Delegate(c Contexter, r Racker) {
	c.Ctx().Position = uint(0)
	c.Ctx().Rack = r
	CallAndReturn(c)
}

// injects a racker with a context, runs it and returns the result to the supplied responder
func Inject(c Contexter, r Racker, responder interface{}) {
	rr := CloneRack(r)
	c.Ctx().Rack = rr
	c.Ctx().Position = uint(0)
	c.Ctx().Stopped = false
	rr.SetResponder(responder)
	CallAndReturn(c)
}

// // Response as an error.
// func (res MCResponse) Error() string {
// 	return fmt.Sprintf("MCResponse status=%v, msg: %s",
// 		res.Status, string(res.Body))
// }
