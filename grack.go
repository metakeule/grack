package grack

import (
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	ŧ "fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	ħ "net/http"
	"path/filepath"
	"runtime"
	"strings"
)

type paramsMod struct {
	*debugCall
	Typ string
	Val interface{}
}

type Rack struct {
	ħ.ResponseWriter
	CheckResponse bool       // check if xml and json strings are valid
	HiJacker      RackerFull // a hijacker that allows interception of calls
	Buffer        *bytes.Buffer
	readwriter    *bufio.ReadWriter

	middlewares        []Middleware
	app                Middleware
	errorHandler       Middleware
	request            *ħ.Request
	params             Params
	mode               string // development / production / stage mode
	pointer            int
	finished           bool
	error              error
	errorStatus        int
	backtrace          []string
	recovered          bool
	parent             RackerFull
	debugCall          *debugCall
	debugMw            map[int]*debugCall
	debugMwCalls       []int
	name               string
	paramsModification map[string][]*paramsMod
}

func NewRack() (ø *Rack) {
	ø = &Rack{}
	if DEBUG {
		ø.debugCall = DebugCall(2)
		ø.debugMw = map[int]*debugCall{}
	}
	ø.middlewares = []Middleware{}
	return
}

func (ø *Rack) Set(key string, i interface{}) {
	if DEBUG {
		if ø.paramsModification[key] == nil {
			ø.paramsModification[key] = []*paramsMod{}
		}

		num := 2
		if ø.HiJacker != nil {
			num = 3
		}

		db := DebugCall(num)
		m := &paramsMod{
			debugCall: db,
			Typ:       "set",
			Val:       i,
		}

		ø.paramsModification[key] = append(ø.paramsModification[key], m)
	}
	ø.params[key] = i
}

func (ø *Rack) Unset(key string) {
	if DEBUG {
		if ø.paramsModification[key] == nil {
			ø.paramsModification[key] = []*paramsMod{}
		}
		num := 2
		if ø.HiJacker != nil {
			num = 3
		}
		db := DebugCall(num)
		m := &paramsMod{
			debugCall: db,
			Typ:       "unset",
		}

		ø.paramsModification[key] = append(ø.paramsModification[key], m)
	}
	delete(ø.params, key)
}

func (ø *Rack) GetMiddleware(i int) Middleware {
	if i < len(ø.middlewares) {
		return ø.middlewares[i]
	}
	return nil
}

func (ø *Rack) Pos() (i int, e error) {
	i = ø.pointer
	if ø.pointer >= len(ø.middlewares) {
		e = ŧ.Errorf("pointer %v not in middleware stack", ø.pointer)
	}
	return
}

func (ø *Rack) Name() string {
	if ø.HiJacker != nil && ø.name == "" {
		return NameOf(ø.HiJacker)
	}
	if ø.name == "" {
		return NameOf(ø)
	}
	return ø.name
}

func (ø *Rack) Debug() {
	if !DEBUG {
		return
	}
	ø.Json(ø.DebugRaw())
}

func (ø *Rack) DebugDetails() {
	if !DEBUG {
		return
	}
	ø.Json(ø.DebugRawDetails())
}

func (ø *Rack) debugMiddleware(pos int) (r map[string]interface{}) {
	mw := ø.middlewares[pos]
	r = map[string]interface{}{}
	r["Name"] = mw.Name()
	r["Pushed"] = ø.debugMw[pos]
	return
}

func (ø *Rack) debugMiddlewares() (r []map[string]interface{}) {
	r = []map[string]interface{}{}
	for i, _ := range ø.middlewares {
		r = append(r, ø.debugMiddleware(i))
	}
	return
}

func (ø *Rack) DebugRaw() (r map[string]interface{}) {
	r = map[string]interface{}{}
	r["Name"] = ø.Name()
	dbMw := ø.debugMiddlewares()
	if ø.app != nil {
		r["App"] = ø.app.Name()
	}
	r["Params"] = ø.params
	called := []interface{}{}
	for _, i := range ø.debugMwCalls {
		if ra, ok := ø.middlewares[i].(RackerFull); ok {
			rdb := ra.DebugRaw()
			delete(rdb, "Params")
			called = append(called, rdb)
			continue
		}
		called = append(called, dbMw[i]["Name"].(string))
	}
	r["Called"] = called
	return
}

func (ø *Rack) DebugRawDetails() (r map[string]interface{}) {
	r = map[string]interface{}{}
	d := ø.debugCall
	r["Name"] = ø.Name()
	r["Created"] = d
	dbMw := ø.debugMiddlewares()
	r["Middlewares"] = dbMw
	if ø.app != nil {
		r["App"] = ø.app.Name()
	}
	r["Params"] = ø.params
	called := []interface{}{}
	for _, i := range ø.debugMwCalls {
		if ra, ok := ø.middlewares[i].(RackerFull); ok {
			rdb := ra.DebugRawDetails()
			called = append(called, rdb)
			continue
		}
		called = append(called, dbMw[i]["Name"].(string))
	}

	mods := map[string][]map[string]interface{}{}
	for kn, mod := range ø.paramsModification {
		if mods[kn] == nil {
			mods[kn] = []map[string]interface{}{}
		}

		for _, mm := range mod {
			mmm := map[string]interface{}{}
			mmm["File"] = mm.File
			mmm["Line"] = mm.Line
			mmm["Function"] = mm.Function
			mmm["Typ"] = mm.Typ
			if mm.Val != nil {
				mmm["Val"] = mm.Val
			}
			mods[kn] = append(mods[kn], mmm)
		}
	}
	r["ParamsModifications"] = mods
	r["Called"] = called
	return
}

func (ø *Rack) SetContentTypeJson() {
	ø.SetContentType("application/json; charset=utf-8")
}

func (ø *Rack) Json(v interface{}) {
	ø.SetContentTypeJson()
	enc := json.NewEncoder(ø)
	ſ := enc.Encode(v)
	if ſ != nil {
		panic("can't encode json: " + ſ.Error())
	}
	return
}

func (ø *Rack) JsonString(jsonStr string) {
	if ø.CheckResponse {
		r := bytes.NewBufferString(jsonStr)
		dec := json.NewDecoder(r)
		var v interface{}
		ſ := dec.Decode(v)
		if ſ != nil {
			panic("not a valid json: " + ſ.Error())
		}
	}
	ø.SetContentTypeJson()
	ø.Write([]byte(jsonStr))
}

func (ø *Rack) Xml(v interface{}) {
	// contenttype is autodetected
	enc := xml.NewEncoder(ø)
	ſ := enc.Encode(v)
	if ſ != nil {
		panic("can't encode xml: " + ſ.Error())
	}
	return
}

func (ø *Rack) XmlString(xmlStr string) {
	if ø.CheckResponse {
		r := bytes.NewBufferString(xmlStr)
		dec := xml.NewDecoder(r)
		var v interface{}
		ſ := dec.Decode(v)
		if ſ != nil {
			panic("not a valid xml: " + ſ.Error())
		}
	}
	// contenttype is autodetected
	ø.Write([]byte(xmlStr))
}

func (ø *Rack) push(mw Middleware) {
	ø.middlewares = append(ø.middlewares, mw)
}

func (ø *Rack) PushFunc(fn func(Racker)) {
	mw := MiddlewareFunc(fn)
	if DEBUG {
		num := 2
		if ø.debugMw == nil {
			ø.debugMw = map[int]*debugCall{}
		}
		ø.debugMw[len(ø.middlewares)] = DebugCall(num)
	}
	ø.push(mw)
}

func (ø *Rack) Push(mw Middleware) {
	if DEBUG {
		num := 2
		if ø.debugMw == nil {
			ø.debugMw = map[int]*debugCall{}
		}
		ø.debugMw[len(ø.middlewares)] = DebugCall(num)
	}
	ø.push(mw)
}

func (ø *Rack) JumpToApp() {
	if ø.app == nil {
		return
	}
	ø.pointer = len(ø.middlewares)
	ø.callMiddleware()
}

func (ø *Rack) Run() {
	ø.error = nil
	ø.errorStatus = 0
	ø.recovered = false
	ø.backtrace = []string{}
	ø.finished = false
	ø.pointer = 0
	if DEBUG {
		ø.debugMwCalls = []int{}
	}
	ø.callMiddleware()
}

func (ø *Rack) InitBuffer() {
	ø.Buffer = &bytes.Buffer{}
	ø.readwriter = bufio.NewReadWriter(
		bufio.NewReader(ø.Buffer),
		bufio.NewWriter(ø.ResponseWriter))
}

func (ø *Rack) RaiseRecovered(err error, backtrace ...string) {
	if ø.parent != nil {
		ø.parent.RaiseRecovered(err, backtrace...)
		return
	}
	ø.recovered = true
	ø.Raise(500, err, backtrace...)
}

// only call for errors that can't be handled on their own
func (ø *Rack) Raise(status int, err error, backtrace ...string) {
	if ø.parent != nil {
		ø.parent.Raise(status, err, backtrace...)
		return
	}
	if status == 0 {
		status = ħ.StatusBadRequest
	}
	if len(backtrace) == 0 {
		for i := 0; i < 100; i++ {
			_, file, line, _ := runtime.Caller(1 + i)
			if file == "" {
				continue
			}
			backtrace = append(backtrace, ŧ.Sprintf("%v: %v", file, line))
		}
	}

	if ø.errorHandler == nil {
		log.Println("No Error Handler defined, taking default")
		if ø.mode == "development" {
			ø.HtmlString(ŧ.Sprintf("<html><body><h1>Error %s (%v - %s)</h1><p>%s</p></body></html>", err.Error(), status, ħ.StatusText(status), strings.Join(backtrace, "<br />")))
		} else {
			log.Printf(`If you want the backtrace in the html output, set mode to "development" (is currently: %#v)`, ø.mode)
			ø.HtmlString(ŧ.Sprintf("<html><body><h1>Error (%v - %s)</h1></body></html>", status, ħ.StatusText(status)))
		}
		log.Printf("Error %s (%v - %s)\n%v\n", err.Error(), status, ħ.StatusText(status), strings.Join(backtrace, "\n"))
		return
	}
	ø.error = err
	ø.errorStatus = status
	ø.backtrace = backtrace
	ø.finished = true
	ø.errorHandler.Call(ø)
}

func (ø *Rack) Call(r Racker) {
	parent := r.(RackerFull)
	ø.SetRequest(parent.Request())
	ø.ResetParams()
	ø.SetParams(parent.Params())
	ø.SetResponseWriter(parent.GetResponseWriter())
	ø.SetParent(parent)
	ø.SetMode(parent.Mode())
	ø.Run()
	if ø.IsFinished() {
		parent.Finish()
	}
}

func (ø *Rack) Delegate(target RackerFull) {
	target.SetRequest(ø.Request())
	target.ResetParams()
	target.SetParams(ø.Params())
	target.SetResponseWriter(ø.ResponseWriter)
	target.SetMode(ø.Mode())
	target.Run()
	if target.IsFinished() {
		ø.Finish()
	}
}

func (ø *Rack) Inject(target RackerFull) {
	target.SetRequest(ø.Request())
	target.ResetParams()
	target.SetParams(ø.Params())
	target.SetResponseWriter(ø.ResponseWriter)
	target.SetParent(ø)
	target.SetMode(ø.Mode())
	target.Run()
	if target.IsFinished() {
		ø.Finish()
	}
}

func (ø *Rack) ParseParams() {
	ſ := ø.request.ParseMultipartForm(0)
	if ſ != nil {
		ø.Raise(ħ.StatusForbidden, ſ)
	}
	vals := ø.request.Form
	for k, val := range vals {
		if len(k) > 2 && k[len(k)-2:len(k)] == "[]" {
			ø.params[k[:len(k)-2]] = val
		} else {
			ø.params[k] = val[0]
		}
	}
	return
}

func (ø *Rack) finish() {
	if !ø.finished && ø.parent != nil {
		ø.parent.Next()
		ø.parent = nil
	}
}

func (ø *Rack) callMiddleware() {
	if ø.HiJacker != nil {
		ø.HiJacker.CallMiddleware()
		return
	}
	ø.CallMiddleware()
}

func (ø *Rack) recover() {
	// This recovers from a panic if one occurred.
	if p := recover(); p != nil {
		err := ŧ.Errorf(ŧ.Sprintf("%v", p))
		backtrace := []string{}
		for i := 0; i < 100; i++ {
			_, file, line, _ := runtime.Caller(2 + i)
			if file == "" {
				continue
			}
			// trace := fmt.Sprintf("File: %#v\nLine: %#v\n", file, line)
			backtrace = append(backtrace, ŧ.Sprintf("%v: %v", file, line))
		}

		ø.RaiseRecovered(err, backtrace...)
		// err := fmt.Sprintf("PANIC: %#v\n"+strings.Join(backtrace, "\n")+"\n", p)

		//fmt.Fprintln(w, err)
		//log.Println(err)
	}
}

func (ø *Rack) _call(mw Middleware, r Racker) {
	defer ø.recover()
	mw.Call(r)
}

func (ø *Rack) CallMiddleware() {
	var arg Racker
	arg = ø
	if ø.HiJacker != nil {
		arg = ø.HiJacker
	}
	if ø.finished {
		ø.finish()
		return
	}
	if ø.pointer == len(ø.middlewares) {
		if ø.app != nil {
			ø._call(ø.app, arg)
			//ø.app.Call(arg)
		}
		ø.finish()
		return
	}
	if ø.pointer > len(ø.middlewares) {
		return
	}
	if DEBUG {
		ø.debugMwCalls = append(ø.debugMwCalls, ø.pointer)
	}
	ø._call(ø.middlewares[ø.pointer], arg)
	// ø.middlewares[ø.pointer].Call(arg)
}

func (ø *Rack) ServeHTTP(writer ħ.ResponseWriter, request *ħ.Request) {
	ø.SetResponseWriter(writer)
	ø.SetRequest(request)
	ø.ResetParams()
	ø.ParseParams()
	ø.Run()
}

// stolen from https://bitbucket.org/kardianos/staticserv/src/66d4675d9ed9897fe218d2dac5eca4b18e05cad2/main.go?at=default
// and changed
func (ø *Rack) uploadedFile(fh *multipart.FileHeader, uploadDir string) (orignalname string, tmpFile string) {
	//formfilename := fh.Filename
	//formHead := fh.Header
	formFile, err := fh.Open()
	// formFile, formHead, err := r.FormFile(formfilename)
	if err != nil {
		return
	}
	defer formFile.Close()

	//remove any directory names in the filename
	//START: work around IE sending full filepath and manually get filename
	itemHead := fh.Header["Content-Disposition"][0]
	lookfor := "filename=\""
	fileIndex := strings.Index(itemHead, lookfor)
	if fileIndex < 0 {
		panic("FileUpload: no filename")
	}
	filename := itemHead[fileIndex+len(lookfor):]
	filename = filename[:strings.Index(filename, "\"")]

	slashIndex := strings.LastIndex(filename, "\\")
	if slashIndex > 0 {
		filename = filename[slashIndex+1:]
	}
	slashIndex = strings.LastIndex(filename, "/")
	if slashIndex > 0 {
		filename = filename[slashIndex+1:]
	}
	_, saveToFilename := filepath.Split(filename)
	//END: work around IE sending full filepath

	//join the filename to the upload dir
	//ŧ.Printf("options: %#v\n", ø.options)

	//saveToFilePath := filepath.Join(uploadDir, saveToFilename)

	//osFile, err := os.Create(saveToFilePath)
	osFile, err := ioutil.TempFile(uploadDir, "_tmp_upload")
	if err != nil {
		panic(err.Error())
	}
	savedFile := osFile.Name()
	defer osFile.Close()

	count, err := io.Copy(osFile, formFile)
	if err != nil {
		panic(err.Error())
	}
	_ = count
	//ŧ.Printf("ALLOW: %s SAVE: %s (%d)\n", r.RemoteAddr, saveToFilename, count)
	//w.Write([]byte("Upload Complete for " + filename))
	orignalname = saveToFilename
	tmpFile = savedFile
	//return map[string]string{saveToFilename: savedFile}
	return
}

func (ø *Rack) UploadedFiles(uploadDir string) (files map[string]string) {
	files = map[string]string{}
	for _, formfiles := range ø.request.MultipartForm.File {
		for _, fileHeader := range formfiles {
			name, tmpFile := ø.uploadedFile(fileHeader, uploadDir)
			files[name] = tmpFile
		}
	}
	return
}

func (ø *Rack) ResetParams() {
	ø.params = Params(map[string]interface{}{})
	if DEBUG {
		ø.paramsModification = map[string][]*paramsMod{}
	}
}

func (ø *Rack) Error() (status int, err error, backtrace []string, recovered bool) {
	status = ø.errorStatus
	err = ø.error
	backtrace = ø.backtrace
	recovered = ø.recovered
	return
}

func (ø *Rack) Finish() {
	ø.finished = true
}

func (ø *Rack) FlushBuffer() (ſ error) {
	ø.InitBuffer()
	if ø.Buffer == nil || ø.readwriter == nil {
		ſ = ŧ.Errorf("no buffer available, call InitBuffer() first")
		return
	}
	return ø.readwriter.Flush()
}

func (ø *Rack) WriteString(s string) (int, error) {
	if ø.Buffer != nil {
		return ø.Buffer.WriteString(s)
	}
	return ø.ResponseWriter.Write([]byte(s))
}

func (ø *Rack) Write(b []byte) (int, error) {
	if ø.Buffer != nil {
		return ø.Buffer.Write(b)
	}
	return ø.ResponseWriter.Write(b)
}

func (ø *Rack) HasError() bool                       { return ø.error != nil }
func (ø *Rack) Get(key string) (i interface{})       { return ø.params[key] }
func (ø *Rack) IsSet(key string) bool                { return ø.params.IsSet(key) }
func (ø *Rack) GetString(key string) string          { return ø.params.String(key) }
func (ø *Rack) GetBytes(key string) []byte           { return ø.params.Bytes(key) }
func (ø *Rack) GetFloat(key string) float32          { return ø.params.Float(key) }
func (ø *Rack) GetFloats(key string) []float32       { return ø.params.Floats(key) }
func (ø *Rack) GetBool(key string) bool              { return ø.params.Bool(key) }
func (ø *Rack) GetInt(key string) int                { return ø.params.Int(key) }
func (ø *Rack) GetStrings(key string) []string       { return ø.params.Strings(key) }
func (ø *Rack) GetBools(key string) []bool           { return ø.params.Bools(key) }
func (ø *Rack) GetInts(key string) []int             { return ø.params.Ints(key) }
func (ø *Rack) Mode() string                         { return ø.mode }
func (ø *Rack) Len() int                             { return len(ø.middlewares) }
func (ø *Rack) HasApp() bool                         { return ø.app != nil }
func (ø *Rack) SetMode(mode string)                  { ø.mode = mode }
func (ø *Rack) SetName(name string)                  { ø.name = name }
func (ø *Rack) Params() Params                       { return ø.params }
func (ø *Rack) Parent() RackerFull                   { return ø.parent }
func (ø *Rack) SetParent(r RackerFull)               { ø.parent = r }
func (ø *Rack) Request() *ħ.Request                  { return ø.request }
func (ø *Rack) SetApp(app Middleware)                { ø.app = app }
func (ø *Rack) SetErrorHandler(h Middleware)         { ø.errorHandler = h }
func (ø *Rack) GetResponseWriter() ħ.ResponseWriter  { return ø.ResponseWriter }
func (ø *Rack) SetResponseWriter(w ħ.ResponseWriter) { ø.ResponseWriter = w }
func (ø *Rack) SetRequest(r *ħ.Request)              { ø.request = r }
func (ø *Rack) IsFinished() bool                     { return ø.finished }
func (ø *Rack) HtmlString(html string)               { ø.Write([]byte(html)) } // contenttype is autodetected
func (ø *Rack) SetHeader(k string, v string)         { h := ø.Header(); h.Set(k, v) }
func (ø *Rack) SetContentType(t string)              { ø.SetHeader("Content-Type", t) }
func (ø *Rack) TextString(text string)               { ø.SetContentType("text/plain"); ø.Write([]byte(text)) }
func (ø *Rack) Skip(n int)                           { ø.pointer = ø.pointer + n; ø.Next() }
func (ø *Rack) GoTo(p int)                           { ø.pointer = p; ø.callMiddleware() } // goto a position p within the middleware stack, handle with care
func (ø *Rack) Next()                                { ø.pointer++; ø.callMiddleware() }
func (ø *Rack) Prev()                                { ø.pointer--; ø.callMiddleware() }
func (ø *Rack) SetParams(p Params)                   { ø.params = p }
