package grack

import (
	"strconv"
	"strings"
)

type Params map[string]interface{}

func (ø Params) IsSet(key string) bool {
	return ø[key] != nil
}

func (ø Params) toBool(i interface{}) bool {
	if str, ok := i.(string); ok {
		str = strings.ToLower(str)
		if str == "false" || str == "true" {
			return str == "true"
		}
		return str == "1"
	}
	return false
}

func (ø Params) Bool(key string) (r bool) {
	i := ø[key]
	if i == nil {
		return
	}
	return ø.toBool(i)
}

func (ø Params) Bytes(key string) (r []byte) {
	i := ø[key]
	if i == nil {
		return
	}
	return ø.toBytes(i)
}

func (ø Params) toBytes(i interface{}) []byte {
	if b, ok := i.([]byte); ok {
		return b
	}
	str := ø.toString(i)
	return []byte(str)
}

func (ø Params) toInt(i interface{}) int {
	if f, ok := i.(int); ok {
		return f
	}
	if f, ok := i.(int32); ok {
		return int(f)
	}
	if f, ok := i.(int64); ok {
		return int(f)
	}
	if s, ok := i.(string); ok {
		if s == "" {
			return 0
		}

		in, ſ := strconv.Atoi(s)
		if ſ != nil {
			return 0
		}
		return in
	}
	return 0
}

func (ø Params) Int(key string) (r int) {
	i := ø[key]
	if i == nil {
		return
	}
	return ø.toInt(i)
}

func (ø Params) Float(key string) (r float32) {
	i := ø[key]
	if i == nil {
		return
	}
	return ø.toFloat(i)
}

func (ø Params) toFloat(i interface{}) float32 {
	if f, ok := i.(float32); ok {
		return f
	}
	if f, ok := i.(float64); ok {
		return float32(f)
	}
	if f, ok := i.(int); ok {
		return float32(f)
	}
	if f, ok := i.(int32); ok {
		return float32(f)
	}
	if f, ok := i.(int64); ok {
		return float32(f)
	}
	if s, ok := i.(string); ok {
		if s == "" {
			return 0
		}

		f, ſ := strconv.ParseFloat(s, 0)
		if ſ == nil {
			return float32(f)
		}
		in, ſ := strconv.Atoi(s)
		if ſ != nil {
			return 0
		}
		return float32(in)
	}
	return 0
}

func (ø Params) toString(i interface{}) string {
	s, ok := i.(string)
	if ok {
		return s
	}
	return ""
}

func (ø Params) String(key string) (r string) {
	i := ø[key]
	if i == nil {
		return
	}
	return ø.toString(i)
}

func (ø Params) Bools(key string) (r []bool) {
	if ø[key] == nil {
		return
	}

	if r, ok := ø[key].([]bool); ok {
		return r
	}

	if b, ok := ø[key].(bool); ok {
		return []bool{b}
	}

	strs, ok := ø[key].([]string)
	if !ok {
		return []bool{ø.toBool(ø[key])}
	}

	r = []bool{}

	for _, s := range strs {
		r = append(r, ø.toBool(s))
	}
	return
}

func (ø Params) Ints(key string) (r []int) {
	if ø[key] == nil {
		return
	}

	if r, ok := ø[key].([]int); ok {
		return r
	}

	if i, ok := ø[key].(int); ok {
		return []int{i}
	}

	strs, ok := ø[key].([]string)
	if !ok {
		return []int{ø.toInt(ø[key])}
	}

	r = []int{}

	for _, s := range strs {
		r = append(r, ø.toInt(s))
	}
	return
}

func (ø Params) Floats(key string) (r []float32) {
	if ø[key] == nil {
		return
	}

	if r, ok := ø[key].([]float32); ok {
		return r
	}

	if f, ok := ø[key].(float32); ok {
		return []float32{f}
	}

	strs, ok := ø[key].([]string)
	if !ok {
		return []float32{ø.toFloat(ø[key])}
	}

	r = []float32{}

	for _, s := range strs {
		r = append(r, ø.toFloat(s))
	}
	return
}

func (ø Params) Strings(key string) (r []string) {
	if ø[key] == nil {
		return
	}

	r, ok := ø[key].([]string)
	if !ok {
		return []string{ø.toString(ø[key])}
	}
	return
}
