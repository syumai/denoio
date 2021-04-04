package denoio

import "syscall/js"

type Callback = func(this js.Value, args []js.Value) interface{}

var global = js.Global()

func newObject() js.Value {
	o := global.Get("Object")
	return o.New()
}

func newUint8Array(size int) js.Value {
	ua := global.Get("Uint8Array")
	return ua.New(size)
}

func newPromise(fn Callback) js.Value {
	p := global.Get("Promise")
	return p.New(js.FuncOf(fn))
}
