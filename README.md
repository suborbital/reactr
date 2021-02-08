![logo_transparent_wide](https://user-images.githubusercontent.com/5942370/107126087-ca589080-687b-11eb-820e-c6161c355eba.png)

Reactr is a fast, performant function scheduling system. Reactr is designed to be flexible, with the ability to run embedded in your Go applications or as a standalone FaaS server, and has first-class support for Wasm/WASI bundles.

Reactr runs functions called Runnables, and transparently spawns workers to process jobs. Each worker processes jobs in sequence, using Runnables to execute them. Reactr jobs are arbitrary data, and they return arbitrary data (or an error). Jobs are scheduled, and their results can be retreived at a later time.

## Wasm

Reactr has support for Wasm-packaged Runnables. The `rwasm` package contains a multi-tenant Wasm scheduler, an API to grant capabilities to Wasm Runnables, and support for several languages including Rust (stable) and Swift (experimental). See [wasm](./docs/wasm.md) and the [subo CLI](https://github.com/suborbital/subo) for details.

## FaaS

Reactr has early (read: alpha) support for acting as a Functions-as-a-Service system. Reactr can be run as a server, accepting jobs from HTTP/S and making the job results available to be fetched later. See [faas](./docs/faas.md) for details.

### The Basics

First, install Reactr's core package `rt`:
```bash
go get github.com/suborbital/reactr/rt
```

And then get started by defining something `Runnable`:
```golang
type generic struct{}

// Run runs a generic job
func (g generic) Run(job rt.Job, ctx *rt.Ctx) (interface{}, error) {
	fmt.Println("doing job:", job.String()) // get the string value of the job's data

	// do your work here

	return fmt.Sprintf("finished %s", job.String()), nil
}

// OnChange is called when Reactr starts or stops a worker to handle jobs,
// and allows the Runnable to set up before receiving jobs or tear down if needed.
func (g generic) OnChange(change rt.ChangeEvent) error {
	return nil
}
```
A `Runnable` is something that can take care of a job, all it needs to do is conform to the `Runnable` interface as you see above.

Once you have a Runnable, create a Reactr instance, register it, and `Do` some work:
```golang
package main

import (
	"fmt"
	"log"

	"github.com/suborbital/reactr/rt"
)

func main() {
	r := rt.New()

	r.Handle("generic", generic{})

	res := r.Do(r.Job("generic", "hard work"))

	res, err := res.Then()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("done!", res.(string))
}
```
When you `Do` some work, you get a `Result`. A result is like a Rust future or a JavaScript promise, it is something you can get the job's result from once it is finished.

Calling `Then()` will block until the job is complete, and then give you the return value from the Runnable's `Run`. Cool, right?

## Reactr has some very powerful capabilities, visit the [get started guide](./docs/getstarted.md) to learn more.

Reactr is being actively developed and has planned improvements, including optimized memory usage, library stability, data persistence, and more. Cheers!

Copyright Suborbital contributors 2020
