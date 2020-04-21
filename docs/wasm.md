# Hive ❤️ WASM

Hive has early support for WASM-packaged runnables. WASM runnables are not ready for production use, but should absolutely be tested to help us root out issues!

WASM support in Hive is powered by [Wasmer](https://github.com/wasmerio/go-ext-wasm), the hard work they've done to create a powerful WASM runtime that is extensible has been very much appreciated, and it's been very cool seeing that project grow.

The `hivew` toolchain is also in its early days, so please bear with us!

The currently "supported" language is Rust, but that only means we are providing the boilerplate needed to use Rust/WASM code. Any language that compiles to WASM can be used if the functions in [lib.rs](https://github.com/suborbital/hivew/blob/master/rs-builder/src/lib.rs) are re-created for that language.

To create a WASM runnable, check out the [hivew CLI](https://github.com/suborbital/hivew). Once you've generated a `.wasm` file, you can use it with Hive just like any other Runnable!

Due to the memory limitations of WASM, WASM runners can only accept a string (rather than arbitrary input) and return a string. WASM runners cannot currently schedule other jobs, though support for that is coming.

Here's how to use it:

First, install hivew's wasm package:
```bash
go get github.com/suborbital/hivew/wasm
```

```golang
h := hive.New()

doWasm := h.Handle("wasm", wasm.NewRunner("path/to/runnable/file.wasm"))

res, err := doWasm("input_must_be_a_string").Then()
if err != nil {
	log.Fatal(err)
}

fmt.Println(res.(string))
```

## Bundles
If you use `hivew` to create a [bundle](https://github.com/suborbital/hivew#bundles), you can load the entire bundle with all of its runnables into your Hive instance:
```golang
if err := wasm.HandleBundle(h, "path/to/runnables.wasm.zip"); err != nil {
	//handle failure
}

res := h.Do(hive.NewJob("name_of_runnable", "input_must_be_a_string"))
[...]
```
The name of each runnable will be the name of the directory that the original source was found in. 

And that's it! You can schedule WASM jobs as normal, and WASM runtimes will be managed automatically to run your jobs.

Please file issues if you encounter anything, and please give the Wasmer team a shout-out for all the great work!