package grack

import (
	ħ "net/http"
)

type Racker interface {
	// error handling
	Raise(status int, err error, backtrace ...string)
	Error() (status int, err error, backtrace []string, recovered bool) // recovered indicates if the error has been from a recovered panic
	HasError() bool

	// debugging
	Name() string
	Mode() string // mode (development, production, stage, ...)
	Debug()
	DebugDetails()

	// raw
	ħ.ResponseWriter
	GetResponseWriter() ħ.ResponseWriter
	Request() *ħ.Request

	// stack navigation
	Finish()
	IsFinished() bool
	Next()
	Prev()
	JumpToApp()
	Delegate(target RackerFull)
	Inject(target RackerFull)

	// shortcuts
	SetHeader(k string, v string)
	SetContentType(t string)

	HtmlString(html string)
	TextString(text string)
	Json(v interface{})
	JsonString(jsonStr string)
	Xml(v interface{})
	XmlString(xmlStr string)
	UploadedFiles(uploadDir string) (files map[string]string)

	// params / context
	Set(key string, i interface{})
	Unset(key string)
	Get(key string) (i interface{})
	IsSet(key string) bool
	GetString(key string) string
	GetBytes(key string) []byte
	GetFloat(key string) float32
	GetFloats(key string) []float32
	GetBool(key string) bool
	GetInt(key string) int
	GetStrings(key string) []string
	GetBools(key string) []bool
	GetInts(key string) []int
}

type RackerFull interface {
	Racker

	// setup
	SetName(string)
	SetMode(string)
	SetApp(Middleware)
	Push(mw Middleware)
	PushFunc(func(Racker))
	SetErrorHandler(Middleware)

	// serving
	ServeHTTP(writer ħ.ResponseWriter, request *ħ.Request)
	SetResponseWriter(ħ.ResponseWriter)
	SetRequest(*ħ.Request)
	ResetParams()
	ParseParams()
	Run()

	// stack navigation special
	Skip(n int)
	GoTo(p int)
	Pos() (i int, e error) // current position in middleware stack
	Len() int              // number of middlewares
	HasApp() bool          // do we have an app?

	// debugging special
	DebugRaw() (r map[string]interface{})
	DebugRawDetails() (r map[string]interface{})
	GetMiddleware(i int) Middleware

	// composition & Hijacking
	Call(parent Racker)
	Params() Params
	SetParams(Params)
	Parent() RackerFull
	SetParent(RackerFull)
	CallMiddleware()
	RaiseRecovered(err error, backtrace ...string)
}
