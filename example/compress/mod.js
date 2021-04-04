import mainwasm from "./mainwasm.ts";
import { Go } from "../wasm_exec.js";
import { decode } from "https://deno.land/std@0.92.0/encoding/base64.ts";

const bytes = decode(mainwasm);
const go = new Go();
const result = await WebAssembly.instantiate(bytes, go.importObject);
go.run(result.instance);

const sleep = (ms) => new Promise(r => setTimeout(r, ms));
await sleep(1000);

/***
 * run compressFile function defined in Go
 */
// 1. open target file in Deno
const file = await Deno.open(Deno.args[0]);
// 2. call Go side's function and pass Deno side's file
const compressed = compressFile(file);
// 3. copy compressed file to Stdout
await Deno.copy(compressed, Deno.stdout);
file.close();
