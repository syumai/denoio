package denoio

import (
	"io"
	"syscall/js"
)

type (
	Reader struct {
		v           js.Value
		syncEnabled bool
	}
	Writer struct {
		v           js.Value
		syncEnabled bool
	}
	Seeker struct {
		v           js.Value
		syncEnabled bool
	}
	Closer struct {
		v js.Value
	}
	ReadWriter struct {
		*Reader
		*Writer
	}
	ReadSeeker struct {
		*Reader
		*Seeker
	}
	ReadCloser struct {
		*Reader
		*Closer
	}
	WriteCloser struct {
		*Writer
		*Closer
	}
	WriteSeeker struct {
		*Writer
		*Seeker
	}
	ReadWriteSeeker struct {
		*Reader
		*Writer
		*Seeker
	}
	ReadWriteCloser struct {
		*Reader
		*Writer
		*Closer
	}
	File struct {
		*Reader
		*Writer
		*Seeker
		*Closer
	}
)

// readWriteSeekCloser is interface type defined to check *File interface implementation.
type readWriteSeekCloser interface {
	io.Reader
	io.Writer
	io.Seeker
	io.Closer
}

var (
	_ io.Reader           = (*Reader)(nil)
	_ io.Writer           = (*Writer)(nil)
	_ io.Seeker           = (*Seeker)(nil)
	_ io.Closer           = (*Closer)(nil)
	_ io.ReadWriter       = (*ReadWriter)(nil)
	_ io.ReadSeeker       = (*ReadSeeker)(nil)
	_ io.ReadCloser       = (*ReadCloser)(nil)
	_ io.WriteCloser      = (*WriteCloser)(nil)
	_ io.WriteSeeker      = (*WriteSeeker)(nil)
	_ io.ReadWriteSeeker  = (*ReadWriteSeeker)(nil)
	_ io.ReadWriteCloser  = (*ReadWriteCloser)(nil)
	_ readWriteSeekCloser = (*File)(nil)
)

func NewReader(v js.Value) *Reader {
	f := v.Get("readSync")
	syncEnabled := f.Type() == js.TypeFunction
	return &Reader{v, syncEnabled}
}

func NewWriter(v js.Value) *Writer {
	f := v.Get("writeSync")
	syncEnabled := f.Type() == js.TypeFunction
	return &Writer{v, syncEnabled}
}

func NewSeeker(v js.Value) *Seeker {
	f := v.Get("seekSync")
	syncEnabled := f.Type() == js.TypeFunction
	return &Seeker{v, syncEnabled}
}

func NewCloser(v js.Value) *Closer {
	return &Closer{v}
}

func NewReadWriter(v js.Value) *ReadWriter {
	return &ReadWriter{
		Reader: NewReader(v),
		Writer: NewWriter(v),
	}
}

func NewReadSeeker(v js.Value) *ReadSeeker {
	return &ReadSeeker{
		Reader: NewReader(v),
		Seeker: NewSeeker(v),
	}
}

func NewReadCloser(v js.Value) *ReadCloser {
	return &ReadCloser{
		Reader: NewReader(v),
		Closer: NewCloser(v),
	}
}

func NewWriteCloser(v js.Value) *WriteCloser {
	return &WriteCloser{
		Writer: NewWriter(v),
		Closer: NewCloser(v),
	}
}

func NewWriteSeeker(v js.Value) *WriteSeeker {
	return &WriteSeeker{
		Writer: NewWriter(v),
		Seeker: NewSeeker(v),
	}
}

func NewReadWriteSeeker(v js.Value) *ReadWriteSeeker {
	return &ReadWriteSeeker{
		Reader: NewReader(v),
		Writer: NewWriter(v),
		Seeker: NewSeeker(v),
	}
}

func NewReadWriteCloser(v js.Value) *ReadWriteCloser {
	return &ReadWriteCloser{
		Reader: NewReader(v),
		Writer: NewWriter(v),
		Closer: NewCloser(v),
	}
}

func NewFile(v js.Value) *File {
	return &File{
		Reader: NewReader(v),
		Writer: NewWriter(v),
		Seeker: NewSeeker(v),
		Closer: NewCloser(v),
	}
}

func read(v js.Value, p []byte) (int, error) {
	ua := newUint8Array(len(p))
	promise := v.Call("read", ua)
	resultCh := make(chan js.Value)
	eofCh := make(chan struct{})
	promise.Call("then", js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		result := args[0]
		if result.IsNull() {
			eofCh <- struct{}{}
		}
		resultCh <- result
		return js.Undefined()
	}))
	select {
	case result := <-resultCh:
		_ = js.CopyBytesToGo(p, ua)
		return result.Int(), nil
	case <-eofCh:
		return 0, io.EOF
	}
}

func readSync(v js.Value, p []byte) (int, error) {
	ua := newUint8Array(len(p))
	result := v.Call("readSync", ua)
	if result.IsNull() {
		return 0, io.EOF
	}
	_ = js.CopyBytesToGo(p, ua)
	return result.Int(), nil
}

func write(v js.Value, p []byte) (int, error) {
	ua := newUint8Array(len(p))
	_ = js.CopyBytesToJS(ua, p)
	promise := v.Call("write", ua)
	resultCh := make(chan js.Value)
	promise.Call("then", js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		resultCh <- args[0]
		return js.Undefined()
	}))
	result := <-resultCh
	return result.Int(), nil
}

func writeSync(v js.Value, p []byte) (int, error) {
	ua := newUint8Array(len(p))
	_ = js.CopyBytesToJS(ua, p)
	result := v.Call("writeSync", ua)
	return result.Int(), nil
}

func seek(v js.Value, offset int64, whence int) (int64, error) {
	promise := v.Call("seek", js.ValueOf(offset), js.ValueOf(whence))
	resultCh := make(chan js.Value)
	promise.Call("then", js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		resultCh <- args[0]
		return js.Undefined()
	}))
	result := <-resultCh
	return int64(result.Int()), nil
}

func seekSync(v js.Value, offset int64, whence int) (int64, error) {
	result := v.Call("seekSync", js.ValueOf(offset), js.ValueOf(whence))
	return int64(result.Int()), nil
}

func (f *Reader) Read(p []byte) (int, error) {
	if f.syncEnabled {
		return readSync(f.v, p)
	}
	return read(f.v, p)
}

func (f *Writer) Write(p []byte) (int, error) {
	if f.syncEnabled {
		return writeSync(f.v, p)
	}
	return write(f.v, p)
}

// Seek
// whence: SeekStart = 0 / SeekCurrent = 1 / SeekEnd = 2
func (f *Seeker) Seek(offset int64, whence int) (int64, error) {
	if f.syncEnabled {
		return seekSync(f.v, offset, whence)
	}
	return seek(f.v, offset, whence)
}

func (f *Closer) Close() error {
	f.v.Call("close")
	return nil
}
