import testmainwasm from "./testmainwasm.ts";
import { Go } from "../lib/wasm_exec.js";
import { decode } from "https://deno.land/std@0.139.0/encoding/base64.ts";
import { Buffer } from "https://deno.land/std@0.139.0/io/buffer.ts";
import {
  assertEquals,
  assertThrows,
} from "https://deno.land/std@0.139.0/testing/asserts.ts";
const dec = new TextDecoder();

const bytes = decode(testmainwasm);
const go = new Go();
const result = await WebAssembly.instantiate(bytes, go.importObject);
go.run(result.instance);

const exampleFile = Deno.openSync("./example.txt");
const exampleBytes = Deno.readAllSync(exampleFile);
const exampleStr = dec.decode(exampleBytes);
exampleFile.close();

Deno.test("readAsyncFromGo", async () => {
  const f = await Deno.open("./example.txt");
  const reader = {
    read() {
      return f.read(...arguments);
    },
  };
  const result = await readAsync(reader);
  f.close();
  assertEquals(result, exampleStr);
});

Deno.test("readerAsyncFromGo", async () => {
  const goReader = getReaderAsyncFromGo()
  const result = await Deno.readAll(goReader);
  assertEquals(result, exampleBytes);
});

Deno.test("readSync", () => {
  const f = Deno.openSync("./example.txt");
  const reader = {
    readSync() {
      return f.readSync(...arguments);
    },
  };
  const goReader = readSync(reader);
  f.close();
  const result = Deno.readAllSync(goReader);
  assertEquals(result, exampleBytes);
});

Deno.test("writeAsyncFromGo", async () => {
    const buf = new Buffer();
    const writer = {
        write() {
            return buf.write(...arguments)
        },
    };
    await writeAsyncFromGo(writer);
    assertEquals(dec.decode(Deno.readAllSync(buf)), "wrote async from Go");
});

Deno.test("writeSyncFromGo", () => {
  const buf = new Buffer();
  const writer = {
    writeSync() {
      return buf.writeSync(...arguments);
    },
  };
  writeSyncFromGo(writer);
  assertEquals(dec.decode(Deno.readAllSync(buf)), "wrote sync from Go");
});

Deno.test("writerAsyncFromGo", async () => {
  const goWriter = getWriterAsyncFromGo();
  await Deno.writeAll(goWriter, exampleBytes);
  const result = getWriterAsyncResult();
  assertEquals(result, exampleStr);
});

Deno.test("writerSyncFromGo", () => {
  const goWriter = getWriterSyncFromGo();
  Deno.writeAllSync(goWriter, exampleBytes);
  const result = getWriterSyncResult();
  assertEquals(result, exampleStr);
});

Deno.test("seekAsync", async () => {
  const f = await Deno.open("./example.txt");
  const seeker = {
    seek() {
      return f.seek(...arguments);
    },
  };
  const goSeeker = seekAsync(seeker);
  const buf = new Uint8Array(5);

  await goSeeker.seek(-5, 2);
  await f.read(buf);
  const last5 = dec.decode(buf);
  assertEquals(last5, "orum.");

  await goSeeker.seek(-10, 2);
  await goSeeker.seek(-5, 1);
  await f.read(buf);
  const last15 = dec.decode(buf);
  assertEquals(last15, "id es");

  await goSeeker.seek(6, 0);
  await f.read(buf);
  const first6 = dec.decode(buf);
  assertEquals(first6, "ipsum");
  f.close();
});

Deno.test("seekSync", () => {
  const f = Deno.openSync("./example.txt");
  const seeker = {
    seekSync() {
      return f.seekSync(...arguments);
    },
  };
  const goSeeker = seekSync(seeker);
  const buf = new Uint8Array(5);

  goSeeker.seekSync(-5, 2);
  f.readSync(buf);
  const last5 = dec.decode(buf);
  assertEquals(last5, "orum.");

  goSeeker.seekSync(-10, 2);
  goSeeker.seekSync(-5, 1);
  f.readSync(buf);
  const last15 = dec.decode(buf);
  assertEquals(last15, "id es");

  goSeeker.seekSync(6, 0);
  f.readSync(buf);
  const first6 = dec.decode(buf);
  assertEquals(first6, "ipsum");
  f.close();
});

Deno.test("close", () => {
  const f = Deno.openSync("./example.txt");
  const goCloser = close(f);
  goCloser.close();
  assertThrows(() => {
    f.close(); // closed file cannot be closed again.
  });
});
