package main

import (
	"fmt"
	. "github.com/metakeule/grack"
	. "github.com/metakeule/grack/http"
	"net/http"
)

// middleware that makes use of AppContext
var pre = func(c Ctx, err error) {
	c.(*AppContext).Layout = "<div style=\"background-color:green;color:white;padding:10px;\">%s</div>"
	Next(c)
}

// normal middleware
var path = func(c Ctx, err error) {
	Add(c, "Path: "+Path(c))
	Next(c)
}

// middleware that makes use of AppContext
var post = func(c Ctx, err error) {
	SetOut(c, fmt.Sprintf(c.(*AppContext).Layout, Out(c)))
}

// responder that makes use of AppRack
func responder(c Ctx, err error) {
	HTML(c)
	hc := c.(HTTPContexter).HTTPContext_()
	hc.Out = "<html><body><h1>App: " + c.Ctx().Rack.(*AppRack).App + "</h1>" + hc.Out + "</body></html>"
	Flush(hc)
}

type AppContext struct {
	*HTTPContext        // important
	Layout       string // specific property
}

type AppRack struct {
	*HTTPRack        // important
	App       string // specific property
}

// important: overwrite this to use the AppContext
func (h *AppRack) NewContext() (c HTTPContexter) {
	return &AppContext{HTTPContext: NewContext()}
}

var rack = &AppRack{HTTPRack: NewRack(), App: "green layouter"}

func init() {
	rack.Push(pre)
	rack.Push(path)
	rack.Push(post)
	rack.SetResponder(responder)
}

func main() {
	fmt.Println("look at http://localhost:8080/strange")
	http.ListenAndServe(":8080", NewServer(rack))
}
