// +build js
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
	if item == js.Undefined {
		val = nil
		return
	}
	str := item.String()
	json.Unmarshal([]byte(str), &val)
}
func Pluralize(count int, word string) string {
	if count == 1 {
		return word
	}
	return word + "s"
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

//router (Director.js)
type Router struct {
	*js.Object
}

func NewRouter() *Router {
	return &Router{js.Global.Get("Router").New()}
}
func (r *Router) On(path string, handler func(string)) {
	r.Call("on", path, handler)
}

func (r *Router) Init(path string) {
	r.Call("init", path)
}
