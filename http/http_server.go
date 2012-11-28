package grack_http

import (
	. "github.com/metakeule/grack"
	"net/http"
)

type HTTPServer struct {
	Rack HTTPRacker
}

func NewServer(r HTTPRacker) *HTTPServer {
	return &HTTPServer{Rack: r}
}

// implements http.Handler
func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	rck := s.Rack
	ctx := rck.NewContext()
	innerCtx := ctx.HTTPContext_()
	innerCtx.Ctx().Rack = rck
	innerCtx.Out = ""
	innerCtx.ResponseWriter = writer
	innerCtx.Request = request
	Reset(innerCtx)
	CallAndReturn(ctx)
}

func CloneHTTPContext(c HTTPContexter) HTTPContexter {
	hc := c.HTTPContext_()
	n := c.Ctx().Rack.(HTTPRacker).NewContext()
	n.HTTPContext_().Out = hc.Out
	n.HTTPContext_().Request = hc.Request
	n.HTTPContext_().ResponseWriter = hc.ResponseWriter
	Reset(n.HTTPContext_())
	return n
}
