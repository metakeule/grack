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

type FakeResponseWriter struct {
	original   http.ResponseWriter
	Out        []byte
	StatusCode int
}

func NewFakeResponseWriter(original http.ResponseWriter) *FakeResponseWriter {
	return &FakeResponseWriter{original, []byte{}, 0}
}

func (ø *FakeResponseWriter) Header() http.Header {
	return ø.original.Header()
}

func (ø *FakeResponseWriter) Write(b []byte) (written int, err error) {
	if ø.StatusCode == 0 {
		ø.StatusCode = 200
	}
	written = 0
	err = nil
	for i := 0; i < len(b); i++ {
		written += 1
		ø.Out = append(ø.Out, b[i])
	}
	return
}

func (ø *FakeResponseWriter) WriteHeader(code int) {
	ø.StatusCode = code
}

type ClonedHTTPContext struct {
	HTTPContexter HTTPContexter
	*FakeResponseWriter
}

func (ø *ClonedHTTPContext) Merge(c HTTPContexter) {
	f := ø.FakeResponseWriter
	if f.StatusCode != 0 {
		c.HTTPContext_().Status = ø.FakeResponseWriter.StatusCode
	}

	if len(f.Out) > 0 {
		Add(c, string(f.Out))
	}
}

func (ø *ClonedHTTPContext) HTTPContext_() *HTTPContext {
	return ø.HTTPContexter.HTTPContext_()
}

func (ø *ClonedHTTPContext) Ctx() *ContextData {
	return ø.HTTPContexter.Ctx()
}

func (ø *ClonedHTTPContext) Rack() Racker {
	return ø.HTTPContexter.Rack()
}

func CloneHTTPContext(c HTTPContexter) *ClonedHTTPContext {
	hc := c.HTTPContext_()
	n := &ClonedHTTPContext{
		HTTPContexter:      c.Ctx().Rack.(HTTPRacker).NewContext(),
		FakeResponseWriter: NewFakeResponseWriter(hc.ResponseWriter),
	}
	//n := c.Ctx().Rack.(HTTPRacker).NewContext()
	n.HTTPContext_().Out = hc.Out
	n.HTTPContext_().Request = hc.Request
	n.HTTPContext_().ResponseWriter = n.FakeResponseWriter
	//n.HTTPContext_().ResponseWriter = hc.ResponseWriter
	Reset(n.HTTPContext_())
	return n
}
