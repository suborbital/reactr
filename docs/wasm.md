# Hive ❤️ Wasm

Hive has early support for Wasm-packaged runnables. Wasm runnables are not ready for production use, but should absolutely be tested to help us root out issues!

Wasm support in Hive is powered by [Wasmer](https://github.com/wasmerio/wasmer-go), the hard work they've done to create a powerful Wasm runtime that is extensible has been very much appreciated, and it's been very cool seeing that project grow.

The `subo` CLI is also in its early days, so please bear with us!

The currently "supported" language is Rust, but that only means we are providing the boilerplate needed to use Rust/Wasm code. Any language that compiles to Wasm can be used if the functions in [lib.rs](https://github.com/suborbital/hive-wasm/blob/master/rs-builder/src/lib.rs) are re-created for that language.

To create a Wasm runnable, check out the [subo CLI](https://github.com/suborbital/subo). Once you've generated a `.wasm` file, you can use it with Hive just like any other Runnable!

Due to the memory limitations of Wasm, Wasm runners accept bytes (rather than arbitrary input) and return bytes. Hive will handle the conversion of inputs and outputs automatically. Wasm runners cannot currently schedule other jobs, though support for that is coming.

Here's how to use hivew Runnables:

First, install hive-wasm's wasm package:
```bash
go get github.com/suborbital/hive-wasm/wasm
```

```golang
h := hive.New()

doWasm := h.Handle("wasm", wasm.NewRunner("path/to/Runnable/file.wasm"))

res, err := doWasm("input_will_be_converted_to_bytes").Then()
if err != nil {
	log.Fatal(err)
}

fmt.Println(string(res.([]byte)))
```

## Bundles
If you use `subo` to create a [bundle](https://github.com/suborbital/subo/blob/main/docs/wasm.md#bundles), you can load the entire bundle with all of its runnables into your Hive instance:
```golang
if err := wasm.HandleBundle(h, "path/to/runnables.wasm.zip"); err != nil {
	//handle failure
}

res := h.Do(hive.NewJob("name_of_runnable", "input_will_be_converted_to_bytes"))
[...]
```
The name of each runnable will be the name of the directory that the original source was found in.

And that's it! You can schedule Wasm jobs as normal, and Wasm runtimes will be managed automatically to run your jobs.

Please file issues if you encounter anything, and please give the Wasmer team a shout-out for all the great work!
