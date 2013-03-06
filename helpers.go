package grack

import (
	// ŧ "fmt"
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
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
