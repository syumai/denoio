//go:generate sh -c "GOOS=js GOARCH=wasm go build -o testmain.wasm ./ && cat testmain.wasm | deno run https://denopkg.com/syumai/binpack/mod.ts > testmainwasm.ts && rm testmain.wasm"
package main

import (
	"bytes"
	_ "embed"
	"io"
	"strings"
	"syscall/js"

	"github.com/syumai/denoio"
)

var (
	writer     bytes.Buffer
	writerSync bytes.Buffer
)

//go:embed example.txt
var exampleBytes []byte

func main() {
	// Go reader => JS reader (read)
	js.Global().Set("readAsync",
		js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			b := bytes.NewReader(exampleBytes)
			return denoio.NewJSReader(b)
		}))

	// JS syncReader => Go reader => JS reader (readSync)
	js.Global().Set("readSync",
		js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			f := denoio.NewReader(args[0])
			var buf bytes.Buffer
			if _, err := io.Copy(&buf, f); err != nil {
				panic(err)
			}
			return denoio.NewJSReader(&buf)
		}))

	// JS writer => Go writer (write)
	js.Global().Set("writeAsyncFromGo",
		js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			f := denoio.NewWriter(args[0])
			r := strings.NewReader("wrote async from Go")
			if _, err := io.Copy(f, r); err != nil {
				panic(err)
			}
			return js.Undefined()
		}))

	// JS syncWriter => Go writer (write)
	js.Global().Set("writeSyncFromGo",
		js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			f := denoio.NewWriter(args[0])
			r := strings.NewReader("wrote sync from Go")
			if _, err := io.Copy(f, r); err != nil {
				panic(err)
			}
			return js.Undefined()
		}))

	// Go writer => JS writer (write)
	js.Global().Set("getWriterAsyncFromGo",
		js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			return denoio.NewJSWriter(&writer)
		}))

	// Go writer => JS writer (writeSync)
	js.Global().Set("getWriterSyncFromGo",
		js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			return denoio.NewJSWriter(&writerSync)
		}))

	js.Global().Set("getWriterAsyncResult",
		js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			return js.ValueOf(writer.String())
		}))

	js.Global().Set("getWriterSyncResult",
		js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			return js.ValueOf(writerSync.String())
		}))

	// Go seeker => JS seeker (seek)
	js.Global().Set("seekAsync",
		js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			f := denoio.NewSeeker(args[0])
			return denoio.NewJSSeeker(f)
		}))

	// Go seeker => JS seeker (seekSync)
	js.Global().Set("seekSync",
		js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			f := denoio.NewSeeker(args[0])
			return denoio.NewJSSeeker(f)
		}))

	// JS closer => Go closer => JS closer (close)
	js.Global().Set("close",
		js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			f := denoio.NewCloser(args[0])
			return denoio.NewJSCloser(f)
		}))

	// block not to exit from program
	select {}
}
