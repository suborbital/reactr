# Reactr ❤️ Wasm

Reactr has first-class support for Wasm-packaged Runnables. Wasm is an incredibly useful modern portable binary format that allows multiple languages to be compiled into .wasm modules.

Wasm support in Reactr is powered by [Wasmer](https://github.com/wasmerio/wasmer-go), the hard work they've done to create a powerful Wasm runtime that is extensible has been very much appreciated, and it's been very cool seeing that project grow.

The current supported languages are Rust and Swift, and the Runnable API is available for each. More languages such as AssemblyScript, Go, and C++ are coming soon!

To create a Wasm runnable, check out the [subo CLI](https://github.com/suborbital/subo). Once you've generated a `.wasm` file, you can use it with Reactr just like any other Runnable!

Due to the memory limitations of Wasm, Wasm runners accept bytes (rather than arbitrary input) and return bytes. Reactr will handle the conversion of inputs and outputs automatically. Wasm runners cannot currently schedule other jobs, though support for that is coming.

Here's how to use Wasm Runnables:

First, install Reactr's Wasm package `rwasm`:
```bash
go get github.com/suborbital/reactr/rwasm
```

```golang
r := rt.New()

doWasm := r.Handle("wasm", rwasm.NewRunner("path/to/runnable/file.wasm"))

res, err := doWasm("input_will_be_converted_to_bytes").Then()
if err != nil {
	log.Fatal(err)
}

fmt.Println(string(res.([]byte)))
```

## Bundles
If you use `subo` to create a [bundle](https://github.com/suborbital/subo/blob/main/docs/wasm.md#bundles), you can load the entire bundle with all of its runnables into your Reactr instance:
```golang
if err := rwasm.HandleBundle(r, "path/to/runnables.wasm.zip"); err != nil {
	//handle failure
}

res := r.Do(rt.NewJob("name_of_runnable", "input_will_be_converted_to_bytes"))
[...]
```
The name of each runnable is defined in the `.runnable.yaml` file it was built with.

And that's it! You can schedule Wasm jobs as normal, and Wasm environments will be managed automatically to run your jobs.

Please file issues if you encounter anything, and please give the Wasmer team a shout-out for all the great work!
