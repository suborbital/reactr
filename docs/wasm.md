# Reactr ❤️ WebAssembly

Reactr has first-class support for WebAssembly-packaged Runnables. Wasm is an incredibly useful modern portable binary format that allows multiple languages to be compiled into .wasm modules.

The current supported languages are Rust (stable), TypeScript/AssemblyScript (beta) and Swift (alpha). The Runnable API is available for each. More languages such as Go and C++ are coming soon!

To create a Wasm runnable, check out the [subo CLI](https://github.com/suborbital/subo). Once you've generated a `.wasm` file, you can use it with Reactr just like any other Runnable!

A multitude of example Runnables can be found in the [testdata directory](../rwasm/testdata).

Due to the memory layout of WebAssembly, Wasm runners accept bytes (rather than arbitrary input) and return bytes. Reactr will handle the conversion of inputs and outputs automatically. Wasm runners cannot currently schedule other jobs.

To get started with Wasm Runnables, install Reactr's WebAssembly package `rwasm`:
```bash
go get github.com/suborbital/reactr/rwasm
```

Install [WasmEdge](https://github.com/WasmEdge/WasmEdge/blob/master/docs/install.md) for using WasmEdge runtime.

For now, WasmEdge runtime has not been merged into main repo, so we should replace it with your local repo in the `go.mod`, something like:
```
replace github.com/suborbital/reactr v0.12.0 => ../second-state/reactr
```

And the latest release of [WasmEdge-go](https://github.com/second-state/WasmEdge-go) has an issue for loading wasm from buffer, so update the dependency to the resolved commit:
```bash
go get github.com/second-state/WasmEdge-go@3b00645
```

```golang
r := rt.New()

doWasm := r.Register("wasm", rwasm.NewRunner("path/to/runnable/file.wasm"))

res, err := doWasm("input_will_be_converted_to_bytes").Then()
if err != nil {
	log.Fatal(err)
}

fmt.Println(string(res.([]byte)))
```

By default, Reactr uses the Wasmer runtime internally, but supports the Wasmtime and WasmEdge runtime as well. Pass `-tags wasmtime` to any `go` command to use Wasmtime. Wasmtime is not yet supported on ARM. Pass `-tags wasmedge` to any `go` command to use WasmEdge.

And that's it! You can schedule Wasm jobs as normal, and Wasm environments will be managed automatically to run your jobs.
