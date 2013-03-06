package grack_http

import (
	"encoding/json"
	"fmt"
	. "github.com/metakeule/grack"
	"html"
	"net/http"
	"strings"
)

// returns the h.Header from the original http.ResponseWriter
func Header(c HTTPContexter) http.Header {
	h := c.HTTPContext_()
	return h.ResponseWriter.Header()
}

// this func writes the response back to the server
func Flush(c HTTPContexter) {
	h := c.HTTPContext_()
	if h.Status != 0 {
		h.ResponseWriter.WriteHeader(h.Status)
	}
	_, err := h.ResponseWriter.Write([]byte(h.Out))
	if err != nil {
		fmt.Printf("Error: %#v", err)
		//panic(err)
	}
}

func Request(c HTTPContexter) *http.Request {
	h := c.HTTPContext_()
	return h.Request
}

func Path(c HTTPContexter) string {
	return Request(c).URL.Path
}

func SetOut(c HTTPContexter, s string) {
	h := c.HTTPContext_()
	h.Out = s
}

func Out(c HTTPContexter) (s string) {
	h := c.HTTPContext_()
	return h.Out
}

func Add(c HTTPContexter, s string) {
	h := c.HTTPContext_()
	h.Out = h.Out + s
}

// sets the status
func SetStatus(c HTTPContexter, status int) {
	h := c.HTTPContext_()
	h.Status = status
}

// returns the Status
func Status(c HTTPContexter) int {
	h := c.HTTPContext_()
	return h.Status
}

// sets the status
func SetContentType(c HTTPContexter, ctype string) {
	h := Header(c)
	h["Content-Type"] = []string{ctype}
}

// TODO make a type switch and let it accept templates
func HTML(c HTTPContexter, html ...string) {
	SetContentType(c, "text/html")
	if len(html) > 0 {
		SetOut(c, html[0])
	}
}

// TODO make a type switch and let it accept templates
func JSON(c HTTPContexter, i ...interface{}) {
	SetContentType(c, "application/json")
	if len(i) > 0 {
		jsoned, _ := json.Marshal(i[0])
		s := string(jsoned)
		SetOut(c, s)
	}
}

// TODO make a type switch and let it accept templates
func CSS(c HTTPContexter, s string) {
	SetContentType(c, "text/css")
	SetOut(c, s)
}

func DebugRequest(r *http.Request) (debug map[string]string) {
	debug = map[string]string{}
	debug["Method"] = r.Method
	debug["Proto"] = r.Proto
	debug["Host"] = r.Host
	debug["RemoteAddr"] = r.RemoteAddr
	debug["RequestURI"] = r.RequestURI
	headers := r.Header
	for k, v := range headers {
		debug["Header "+k] = strings.Join(v, "<br>")
	}
	//debug["URL"] = fmt.Sprintf("%#v", r.URL)
	// debug["Form"] = fmt.Sprintf("%#v", r.Form)
	// debug["MultipartForm"] = fmt.Sprintf("%#v", r.MultipartForm)
	// debug["Body"] = fmt.Sprintf("%#v", r.Body)
	// debug["ContentLength"] = fmt.Sprintf("%d", r.ContentLength)
	// debug["TransferEncoding"] = fmt.Sprintf("%#v", r.TransferEncoding)
	return
}

func DebugContext(c HTTPContexter) (contextDebug map[string]string, requestDebug map[string]string, rackDebug map[string]string) {
	contextDebug = map[string]string{}
	rackDebug = DebugRack(c)
	requestDebug = DebugRequest(c.HTTPContext_().Request)
	cc := c.Ctx()
	contextDebug["Name"] = ContextName(c)
	contextDebug["Inspect"] = fmt.Sprintf("%s", c)
	//contextDebug["HTTPContext"] = fmt.Sprintf("%#v", c.HTTPContext_())
	contextDebug["Out"] = c.HTTPContext_().Out
	contextDebug["Status"] = fmt.Sprintf("%#d", c.HTTPContext_().Status)
	contextDebug["Position"] = fmt.Sprintf("%d", cc.Position)
	contextDebug["Stopped"] = fmt.Sprintf("%v", cc.Stopped)
	return
}

func MapToHTMLTable(data map[string]string) (table string) {
	table = "<table>"
	for k, v := range data {
		table += "<tr><td class=\"key\" style=\"vertical-align:top;font-weight:bold\">" + k + "</td><td class=\"val\">" + html.EscapeString(v) + "</td></tr>"
	}
	table += "</table>"
	return
}
