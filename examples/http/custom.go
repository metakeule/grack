package main

import (
	"fmt"
	. "github.com/metakeule/grack"
	. "github.com/metakeule/grack/http"
	"net/http"
)

// nice helpers that make use of NamedRack and LayoutContext
func Name(c Ctx) string         { return c.Rack().(*NamedRack).Name }
func SetLayout(c Ctx, l string) { c.(*LayoutContext).Layout = l }
func Layout(c Ctx) string       { return c.(*LayoutContext).Layout }

// middleware that makes use of LayoutContext (sets the layout)
var pre = func(c Ctx, err error) {
	SetLayout(c, "<div style=\"background-color:green;color:white;padding:10px;\">%s</div>")
	Next(c)
}

// normal middleware
var path = func(c Ctx, err error) {
	Add(c, "Path: "+Path(c))
	Next(c)
}

// middleware that makes use of LayoutContext (embeds the context into the layout)
var post = func(c Ctx, err error) { SetOut(c, fmt.Sprintf(Layout(c), Out(c))) }

// responder that makes use of NamedRack
func responder(c Ctx, err error) {
	HTML(c)
	SetOut(c, "<html><body><h1>Name: '"+Name(c)+"'</h1>"+Out(c)+"</body></html>")
	Flush(c)
}

type LayoutContext struct {
	*HTTPContext        // important
	Layout       string // additional property
}

type NamedRack struct {
	*HTTPRack        // important
	Name      string // additional property
}

// important: overwrite this to use the LayoutContext
func (h *NamedRack) NewContext() (c HTTPContexter) {
	return &LayoutContext{HTTPContext: NewHTTPContext()}
}

var rack = &NamedRack{HTTPRack: NewHTTPRack(), Name: "greened"}

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
