# denoio

- denoio is a Go package to bind Deno and Go(Wasm)'s I/O interfaces.

## Usage

- In Go, `denoio.NewReader(v)` converts JS's `Deno.Reader` into Go's `io.Reader`.
- Similarly, in Go, `denoio.NewJSReader(r)` converts Go's `io.Reader` into JS's
  `Deno.Reader`.
- In the same way, this package can convert `io.Writer`, `io.Seeker` and
  `io.Closer`.

## Example

gzip compressing example. (see example/compress for detail)

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

## Status

### Converting JS (Deno) to Go (Wasm)

- [ ] JS reader to Go reader
- [x] JS syncReader to Go reader
- [ ] JS writer to Go writer
- [x] JS syncWriter to Go writer
- [x] JS seeker to Go seeker
- [x] JS syncSeeker to Go seeker
- [x] JS closer to Go closer

### Converting Go (Wasm) to JS (Deno)

- [x] Go reader to JS reader
- [x] Go reader to JS syncReader
- [x] Go writer to JS writer
- [x] Go writer to JS syncWriter
- [x] Go seeker to JS seeker
- [x] Go seeker to JS syncSeeker
- [x] Go closer to JS closer

## Author

syumai

## License

MIT
