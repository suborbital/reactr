# Hive ❤️ WASM

Hive has early support for WASM-packaged runnables. WASM runnables are not ready for production use, but should absolutely be tested to help us root out issues!

WASM support in Hive is powered by [Wasmer](https://github.com/wasmerio/go-ext-wasm), the hard work they've done to create a powerful WASM runtime that is extensible has been very much appreciated, and it's been very cool seeing that project grow.

The process for creating a WASM runnable is currently also _very rough_, so bear with us!

The currently "supported" language is Rust, but that only means we are providing the boilerplate needed to use Rust/WASM code. Any language that compiles to WASM can be used if the functions in `src/lib.rs` are re-created for that language. In the future, a proper WASM repo will be created for Hive that will include the boilerplate for a number of languages.

To get started, you'll need the `Cargo` toolchain installed (with a recent version of Rust), [`wasm-pack`](https://rustwasm.github.io/wasm-pack/installer/), and the `wasm32-unknown-unknown` rust compilation target installed. You should also clone this repo.

In the future, this will all be Dockerized so you won't need to concern yourself :)

Due to the memory limitations of WASM, WASM runners accept a string (rather than arbitrary input) and return a string. WASM runners cannot currently schedule other jobs, though support for that is coming.

In this repo, look at `src/run.rs`. You'll find the `run` function already defined for you. Add whatever code you want, just don't change the function signature.

Once you're done, run `make wasm`, which will generate `pkg/wasm_runner_bg.wasm` (among other things). This is your WASM runner file, and should be included wherever you want to run Hive with WASM.

Here's how to use it:
```golang
h := hive.New()

doWasm := h.Handle("wasm", hive.NewWasm("{path/to/wasm_runner_bg.wasm}"))

res, err := doWasm("input_must_be_a_string").Then()
if err != nil {
	log.Fatal(err)
}

fmt.Println(res.(string))
```

And that's it! You can schedule WASM jobs as normal, and WASM instances will be managed automatically to run your jobs.

Please file issues if you encounter anything, and please give the Wasmer team a shout-out for all the great work!