package main

import "syscall/js"

var (
	global     = js.Global()
	object     = global.Get("Object")
	promise    = global.Get("Promise")
	uint8Array = global.Get("Uint8Array")
)

func newObject() js.Value {
	return object.New()
}

func newUint8Array(size int) js.Value {
	return uint8Array.New(size)
}

func newPromise(fn js.Func) js.Value {
	return promise.New(fn)
}
