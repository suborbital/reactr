## Reactr has been deprecated. You can use the new [scheduler](https://github.com/suborbital/e2core/tree/main/scheduler) and [engine](https://github.com/suborbital/sat/tree/vmain/engine) packages, which are a drop-in replacements for this project. 

# Reactr ❤️ WebAssembly

Reactr has first-class support for WebAssembly-packaged Runnables. Wasm is an incredibly useful modern portable binary format that allows multiple languages to be compiled into .wasm modules.

The current supported languages are Rust (stable) and TypeScript/JavaScript. The Runnable API is available for each.

To create a Wasm runnable, check out the [subo CLI](https://suborbital,dev/subo). Once you've generated a `.wasm` file, you can use it with Reactr just like any other Runnable!

A multitude of example Runnables can be found in the [testdata directory](https://github.com/suborbital/reactr/tree/main/engine/testdata).

Due to the memory layout of WebAssembly, Wasm runners accept bytes (rather than arbitrary input) and return bytes. Reactr will handle the conversion of inputs and outputs automatically. Wasm runners cannot currently schedule other jobs.

To get started with Wasm Runnables, install Reactr's WebAssembly package `engine`:

```bash
go get github.com/suborbital/reactr/engine
```

```go
r := rt.New()

doWasm := r.Register("wasm", engine.NewRunner("path/to/runnable/file.wasm"))

res, err := doWasm("input_will_be_converted_to_bytes").Then()
if err != nil {
	log.Fatal(err)
}

fmt.Println(string(res.([]byte)))
```

By default, Reactr uses the Wasmer runtime internally, but supports the Wasmtime runtime as well. Pass `-tags wasmtime` to any `go` command to use Wasmtime. Wasmtime is not yet supported on ARM.

And that's it! You can schedule Wasm jobs as normal, and Wasm environments will be managed automatically to run your jobs.
