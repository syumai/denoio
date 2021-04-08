//go:generate sh -c "GOOS=js GOARCH=wasm go build -o main.wasm ./ && cat main.wasm | deno run https://denopkg.com/syumai/binpack/mod.ts > mainwasm.ts && rm main.wasm"
package main

import (
	"bytes"
	"compress/gzip"
	"io"
	"syscall/js"

	"github.com/syumai/denoio"
)

func main() {
	// register function as "compressFile"
	js.Global().Set("compressFile",
		js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
			var buf bytes.Buffer
			zw := gzip.NewWriter(&buf)

			// 1. Convert Deno (JS) side's Deno.File to Go's io.Reader interface
			f := denoio.NewReader(args[0])
			if _, err := io.Copy(zw, f); err != nil {
				panic(err)
			}
			if err := zw.Flush(); err != nil {
				panic(err)
			}
			if err := zw.Close(); err != nil {
				panic(err)
			}

			// 2. Convert Go side's io.Reader to Deno (JS) side's Deno.Reader interface.
			return denoio.NewJSReader(&buf)
		}))

	// block not to exit from program
	select {}
}
