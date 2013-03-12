package grack

import (
	// ŧ "fmt"
	ħ "net/http"
	"reflect"
	"runtime"
	"strings"
)

type debugCall struct {
	Function string
	File     string
	Line     uint
}

var mwFuncType = reflect.TypeOf(func(r *Rack) {})

func DebugCall(num int) *debugCall {
	pc, file, line, _ := runtime.Caller(num)
	n := runtime.FuncForPC(pc).Name()
	a := strings.Split(n, "/")
	name := a[len(a)-1]
	return &debugCall{Function: name, File: file, Line: uint(line)}
}

func NameOf(i interface{}) string {
	t := reflect.TypeOf(i)
	if t == mwFuncType {
		return ""
	}
	//ŧ.Printf("type: %v\n", t)
	el := t.Elem()
	return el.Name()
}

var DEBUG = false

// see http://stackoverflow.com/questions/7052693/how-to-get-the-name-of-a-function-in-go
func GetFunctionName(i interface{}) string {
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer())
	if fn != nil {
		return fn.Name()
	}
	return ""
}

func HandleFunc(fn func(w ħ.ResponseWriter, r *ħ.Request)) Middleware {
	return MiddlewareFunc(func(r Racker) { fn(r, r.Request()) })
}

func Handle(handler ħ.Handler) Middleware {
	return MiddlewareFunc(func(r Racker) { handler.ServeHTTP(r, r.Request()) })
}
