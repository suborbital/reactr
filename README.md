![logo_transparent](https://user-images.githubusercontent.com/5942370/88548780-87288580-cfed-11ea-8239-991b6ac420e3.png)

![Testapalooza](https://github.com/suborbital/hive/workflows/Testapalooza/badge.svg)

Hive is a fast, performant job scheduling system, plain and simple. Hive is designed to be flexible, with the ability to run embedded in your Go applications or as a standalone FaaS server, and has early support for running WASM/WASI bundles.

Hive transparently spawns workers to process jobs, with each worker processing jobs in sequence. Hive jobs are arbitrary data, and they return arbitrary data (or an error). Jobs are scheduled by clients, and their results can be retreived at a later time.

## WASM

Hive has early (read: alpha) support for acting as a Functions-as-a-Service system. Hive can be run as a server, accepting jobs from HTTP/S and making the job results available to be fetched later. See [faas](./docs/faas.md) for details. gRPC support is planned.

## FaaS

Hive has early (read: alpha) support for Wasm-packaged runnables. This is actively being worked on, as Wasm is an exciting new standard that makes cross-language and cross-platform code just a bit easier :) See [wasm](./docs/wasm.md) and the [subo CLI](https://github.com/suborbital/subo) for details.

### The Basics

First, install Hive:
```bash
go get github.com/suborbital/hive/hive
```

And then get started by defining something `Runnable`:
```golang
type generic struct{}

// Run runs a generic job
func (g generic) Run(job hive.Job, do hive.DoFunc) (interface{}, error) {
	fmt.Println("doing job:", job.String()) // get the string value of the job's data

	// do your work here

	return fmt.Sprintf("finished %s", job.String()), nil
}

// OnStart is called when Hive starts up a worker to handle jobs,
// and allows the Runnable to set itself up before receiving jobs
func (g generic) OnStart() error {
	return nil
}
```
A `Runnable` is something that can take care of a job, all it needs to do is conform to the `Runnable` interface as you see above.

Once you have a Runnable, create a hive, register it, and `Do` some work:
```golang
package main

import (
	"fmt"
	"log"

	"github.com/suborbital/hive/hive"
)

func main() {
	h := hive.New()

	h.Handle("generic", generic{})

	r := h.Do(h.Job("generic", "hard work"))

	res, err := r.Then()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("done!", res.(string))
}
```
When you `Do` some work, you get a `Result`. A result is like a Rust future or a JavaScript promise, it is something you can get the job's result from once it is finished.

Calling `Then()` will block until the job is complete, and then give you the return value from the Runnable's `Run`. Cool, right?

### Hive has some very powerful capabilities, visit the [get started guide](./docs/getstarted.md) to learn more.

Hive is being actively developed and has planned improvements, including optimized memory usage, library stability, data persistence, and more. Cheers!

Copyright Suborbital contributors 2020
