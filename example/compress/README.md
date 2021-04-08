# compress example

## Usage

```
go generate
deno run --allow-read mod.js ${file name to compress} > compressed.gz
gzip -c -d compressed.gz > decompressed
```
