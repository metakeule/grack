package grack_base

import (
	"fmt"
	. "github.com/metakeule/grack"
)

func NewRack() *Rack {
	return &Rack{Stack: []MiddlewareCaller{}, Created: DebugCall()}
}

func DebugContext(c Contexter) (contextDebug map[string]string, rackDebug map[string]string) {
	contextDebug = map[string]string{}
	rackDebug = DebugRack(c)
	cc := c.Ctx()
	contextDebug["Name"] = ContextName(c)
	contextDebug["Inspect"] = fmt.Sprintf("%v", c)
	contextDebug["Position"] = fmt.Sprintf("%d", cc.Position)
	contextDebug["Stopped"] = fmt.Sprintf("%v", cc.Stopped)
	return
}

/*
	A simple and general implementation of an IO context
*/
func NewIO() IOer {
	return &IO{ContextData: &ContextData{}}
}

type IOer interface {
	Contexter
	Io() *IO
}

type IO struct {
	*ContextData
	In  interface{}
	Out interface{}
}

func (c *IO) Ctx() *ContextData {
	return c.ContextData
}

func (c *IO) Io() *IO {
	return c
}

/*
	Helper for the IO Context
*/

// get the input of a Contexter
func In(c Contexter) interface{} {
	return c.(IOer).Io().In
}

// get the output of a Contexter
func Out(c Contexter) interface{} {
	return c.(IOer).Io().Out
}

// get the input as String
func InString(c Contexter) string {
	i := In(c)
	if i == nil {
		return ""
	}
	return i.(string)
}

// get the output as String
func OutString(c Contexter) string {
	o := Out(c)
	if o == nil {
		return ""
	}
	return o.(string)
}

// set the output of a Contexter
func SetOut(c Contexter, val interface{}) {
	c.(IOer).Io().Out = val
}

// set the input of a Contexter
func SetIn(c Contexter, val interface{}) {
	c.(IOer).Io().In = val
}
