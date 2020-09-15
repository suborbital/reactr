# Hive ❤️ WASM

Hive has early support for WASM-packaged Runnables. WASM Runnables are not ready for production use, but should absolutely be tested to help us root out issues!

WASM support in Hive is powered by [Wasmer](https://github.com/wasmerio/go-ext-wasm), the hard work they've done to create a powerful WASM runtime that is extensible has been very much appreciated, and it's been very cool seeing that project grow.

The `hivew` toolchain is also in its early days, so please bear with us!

The currently "supported" language is Rust, but that only means we are providing the boilerplate needed to use Rust/WASM code. Any language that compiles to WASM can be used if the functions in [lib.rs](https://github.com/suborbital/hivew-rs-builder/blob/master/src/lib.rs) are re-created for that language.

To create a WASM Runnable, check out the [hivew CLI](https://github.com/suborbital/hivew). Once you've generated a `.wasm` file, you can use it with Hive just like any other Runnable!

WASM Runnables accept `[]byte` as input, and so jobs scheduled for a WASM Runnable must be able to convert to `[]byte`. hivew will automatically make this conversion, with `string` being cast to `[]byte` and `struct` types being JSON marshalled. WASM Runnables always return `[]byte`. The hivew FFI Runnable API does not currently support scheduling new jobs from a WASM Runnable, but that ability is planned.

Here's how to use hivew Runnables:

First, install hivew's wasm package:
```bash
go get github.com/suborbital/hivew/wasm
```

```golang
h := hive.New()

doWasm := h.Handle("wasm", wasm.NewRunner("path/to/Runnable/file.wasm"))

res, err := doWasm("some_input_that_can_become_bytes").Then()
if err != nil {
	log.Fatal(err)
}

resStr := string(res.([]byte]))

fmt.Println(resStr)
```

## Bundles
If you use `hivew` to create a [bundle](https://github.com/suborbital/hivew#bundles), you can load the entire bundle with all of its Runnables into your Hive instance:
```golang
if err := wasm.HandleBundle(h, "path/to/runnables.wasm.zip"); err != nil {
	//handle failure
}

res := h.Do(hive.NewJob("name_of_Runnable", "some_input_that_can_become_bytes"))
[...]
```
The name of each Runnable will be the name of the directory that the original source was found in. 

And that's it! You can schedule WASM jobs as usual, and WASM runtimes will be managed automatically to run your jobs.

Please file issues if you encounter anything, and please give the Wasmer team a shout-out for all the great work!
