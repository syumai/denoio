package denoio

import "syscall/js"

var (
	global          = js.Global()
	objectClass     = global.Get("Object")
	promiseClass    = global.Get("Promise")
	uint8ArrayClass = global.Get("Uint8Array")
)

func newObject() js.Value {
	return objectClass.New()
}

func newUint8Array(size int) js.Value {
	return uint8ArrayClass.New(size)
}

func newPromise(fn js.Func) js.Value {
	return promiseClass.New(fn)
}
