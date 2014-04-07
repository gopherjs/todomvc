package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gopherjs/gopherjs/js"
	"math/rand"
	"time"
)

func Store(key string, val interface{}) {
	byteArr, _ := json.Marshal(val)
	str := string(byteArr)
	js.Global.Get("localStorage").Call("setItem", key, str)
}
func Retrieve(key string, val interface{}) {
	item := js.Global.Get("localStorage").Call("getItem", key)
	if item.IsNull() {
		val = nil
		return
	}
	str := item.Str()
	json.Unmarshal([]byte(str), &val)
}
func Pluralize(count int, word string) string {
	if count == 1 {
		return word
	}
	return word + "s"
}

func UuidJS() string {
	uuid := ""
	for i := 0; i < 32; i++ {
		rand := int(js.Global.Get("Math").Call("random").Float()*16) | 0
		switch i {
		case 8, 12, 16, 20:
			uuid += "-"
		}
		switch i {
		case 12:
			uuid += "4"
		case 16:
			uuid += js.Global.Get("Number").New(rand&3|8).Call("toString", 16).Str()
		default:
			uuid += js.Global.Get("Number").New(rand).Call("toString", 16).Str()
		}
	}
	return uuid
}
func Uuid() (uuid string) {
	for i := 0; i < 32; i++ {
		rand.Seed(time.Now().UnixNano() + int64(i))
		random := rand.Intn(16)
		switch i {
		case 8, 12, 16, 20:
			uuid += "-"
		}
		switch i {
		case 12:
			uuid += fmt.Sprintf("%x", 4)
		case 16:
			uuid += fmt.Sprintf("%x", random&3|8)
		default:
			uuid += fmt.Sprintf("%x", random)
		}
	}
	return
}

//handlebar templates
type Handlebar struct {
	js.Object
}

func CompileHandlebar(template string) *Handlebar {
	h := js.Global.Get("Handlebars").Call("compile", template)
	return &Handlebar{h}
}
func RenderHandlebar(hb *Handlebar, i interface{}) string {
	return hb.Object.Invoke(i).Str()
}
func RegisterHandlebarsHelper() {
	fn := func(a, b, options js.Object) js.Object {
		if a.Str() == b.Str() {
			return options.Call("fn", js.This)
		} else {
			return options.Call("inverse", js.This)
		}
	}
	js.Global.Get("Handlebars").Call("registerHelper", "eq", fn)
}

//router (Director.js)
type Router struct {
	js.Object
}

func NewRouter() Router {
	return Router{Object: js.Global.Get("Router").New()}
}
func (r Router) On(path string, handler func(string)) {
	r.Call("on", path, handler)
}

func (r Router) Init(path string) {
	r.Call("init", path)
}
