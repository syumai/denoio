# denoio

* denoio is a Go package to bind Deno and Go's I/O interfaces.

## Status

WIP

## Usage

* gzip compressing example using tinygo
  - see example/compress for detail

### Go side

  ```go
// register function as "compressFile"
js.Global().Set("compressFile",
	js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		var buf bytes.Buffer
		zw := gzip.NewWriter(&buf)

		// 1. Convert Deno (JS) side's Deno.File to Go(Wasm)'s io.Reader interface
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
```

### Deno (JS) side

```js
// 1. open target file in Deno
const file = await Deno.open(Deno.args[0]);
// 2. call Go side's function and pass Deno side's file
const compressed = compressFile(file);
// 3. copy compressed file to Stdout
Deno.copy(compressed, Deno.stdout);
```

## Author

syumai

## License

MIT
