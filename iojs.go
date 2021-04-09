package denoio

import (
	"io"
	"syscall/js"
)

type (
	jsReader struct {
		r io.Reader
	}
	jsWriter struct {
		w io.Writer
	}
	jsSeeker struct {
		s io.Seeker
	}
	jsCloser struct {
		c io.Closer
	}
)

func registerReadFunc(obj js.Value, r io.Reader) {
	jr := &jsReader{r}
	readFunc := js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		var cb js.Func
		cb = js.FuncOf(func(_ js.Value, pArgs []js.Value) interface{} {
			defer cb.Release()
			resolve := pArgs[0]
			go func() {
				n, err := jr.Read(args[0])
				if err != nil {
					panic(err)
				}
				resolve.Invoke(n)
			}()
			return js.Undefined()
		})
		return newPromise(cb)
	})
	readSyncFunc := js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		n, err := jr.Read(args[0])
		if err != nil {
			panic(err)
		}
		return n
	})
	obj.Set("read", readFunc)
	obj.Set("readSync", readSyncFunc)
}

func registerWriteFunc(obj js.Value, w io.Writer) {
	jw := &jsWriter{w}
	writeFunc := js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		var cb js.Func
		cb = js.FuncOf(func(_ js.Value, pArgs []js.Value) interface{} {
			defer cb.Release()
			resolve := pArgs[0]
			go func() {
				n, err := jw.Write(args[0])
				if err != nil {
					panic(err)
				}
				resolve.Invoke(n)
			}()
			return js.Undefined()
		})
		return newPromise(cb)
	})
	writeSyncFunc := js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		n, err := jw.Write(args[0])
		if err != nil {
			panic(err)
		}
		return n
	})
	obj.Set("write", writeFunc)
	obj.Set("writeSync", writeSyncFunc)
}

func registerSeekFunc(obj js.Value, s io.Seeker) {
	jss := &jsSeeker{s}
	seekFunc := js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		var cb js.Func
		cb = js.FuncOf(func(_ js.Value, pArgs []js.Value) interface{} {
			defer cb.Release()
			ints := make([]interface{}, len(pArgs))
			for i, v := range pArgs {
				ints[i] = v
			}
			resolve := pArgs[0]
			go func() {
				offset, err := jss.Seek(args[0], args[1])
				if err != nil {
					panic(err)
				}
				resolve.Invoke(offset)
			}()
			return js.Undefined()
		})
		return newPromise(cb)
	})
	seekSyncFunc := js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		offset, err := jss.Seek(args[0], args[1])
		if err != nil {
			panic(err)
		}
		return offset
	})
	obj.Set("seek", seekFunc)
	obj.Set("seekSync", seekSyncFunc)
}

func registerCloseFunc(obj js.Value, c io.Closer) {
	jc := &jsCloser{c}
	closeFunc := js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		err := jc.Close()
		if err != nil {
			panic(err)
		}
		return nil
	})
	obj.Set("close", closeFunc)
}

func NewJSReader(v io.Reader) js.Value {
	obj := newObject()
	registerReadFunc(obj, v)
	return obj
}

func NewJSWriter(v io.Writer) js.Value {
	obj := newObject()
	registerWriteFunc(obj, v)
	return obj
}

func NewJSSeeker(v io.Seeker) js.Value {
	obj := newObject()
	registerSeekFunc(obj, v)
	return obj
}

func NewJSCloser(v io.Closer) js.Value {
	obj := newObject()
	registerCloseFunc(obj, v)
	return obj
}

func NewJSReadWriter(v io.ReadWriter) js.Value {
	obj := newObject()
	registerReadFunc(obj, v)
	registerWriteFunc(obj, v)
	return obj
}

func NewJSReadSeeker(v io.ReadSeeker) js.Value {
	obj := newObject()
	registerReadFunc(obj, v)
	registerSeekFunc(obj, v)
	return obj
}

func NewJSReadCloser(v io.ReadCloser) js.Value {
	obj := newObject()
	registerReadFunc(obj, v)
	registerCloseFunc(obj, v)
	return obj
}

func NewJSWriteCloser(v io.WriteCloser) js.Value {
	obj := newObject()
	registerWriteFunc(obj, v)
	registerCloseFunc(obj, v)
	return obj
}

func NewJSWriteSeeker(v io.WriteSeeker) js.Value {
	obj := newObject()
	registerWriteFunc(obj, v)
	registerSeekFunc(obj, v)
	return obj
}

func NewJSReadWriteSeeker(v io.ReadWriteSeeker) js.Value {
	obj := newObject()
	registerReadFunc(obj, v)
	registerWriteFunc(obj, v)
	registerSeekFunc(obj, v)
	return obj
}

func NewJSReadWriteCloser(v io.ReadWriteCloser) js.Value {
	obj := newObject()
	registerReadFunc(obj, v)
	registerWriteFunc(obj, v)
	registerCloseFunc(obj, v)
	return obj
}

func NewJSFile(v readWriteSeekCloser) js.Value {
	obj := newObject()
	registerReadFunc(obj, v)
	registerWriteFunc(obj, v)
	registerSeekFunc(obj, v)
	registerCloseFunc(obj, v)
	return obj
}

func (jr *jsReader) Read(p js.Value) (js.Value, error) {
	b := make([]byte, p.Length())
	result, err := jr.r.Read(b)
	if err == io.EOF {
		return js.Null(), nil
	}
	if err != nil {
		return js.Undefined(), err
	}
	_ = js.CopyBytesToJS(p, b)
	return js.ValueOf(result), nil
}

func (jw *jsWriter) Write(p js.Value) (js.Value, error) {
	b := make([]byte, p.Length())
	_ = js.CopyBytesToGo(b, p)
	result, err := jw.w.Write(b)
	if err != nil {
		return js.Undefined(), err
	}
	return js.ValueOf(result), nil
}

func (jss *jsSeeker) Seek(offset js.Value, whence js.Value) (js.Value, error) {
	result, err := jss.s.Seek(int64(offset.Int()), whence.Int())
	if err != nil {
		return js.Undefined(), err
	}
	return js.ValueOf(result), nil
}

func (jc *jsCloser) Close() error {
	return jc.c.Close()
}
