# hive

Hive is a job scheduler, plain and simple. It is designed to allow multiple workers to exist in paralell, with each worker processing jobs in sequence.

Hive jobs are arbitrary data, and they return arbitrary data (or an error). Jobs are scheduled, and their results can be retreived at a later time.

To get started, define something `Runnable`:
```golang
type generic struct{}

// Run runs a generic job
func (g generic) Run(job hive.Job, run hive.RunFunc) (interface{}, error) {
	fmt.Println("doing job:", job.String()) // get the string value of the job's data

	// do your work here

	return fmt.Sprintf("finished %s", job.String()), nil
}
```
A `Runnable` is something that can take care of a job, all it needs to do is conform to the `Runnable` interface as you see above.

Once you have a Runnable, create a hive, register it, and `Do` some work:
```golang
package main

import (
	"fmt"
	"log"

	"github.com/suborbital/hive"
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

Calling `Then()` will block until the job is complete, and then give you the return value from the Runnable's `Run`. Make sense?

There are some more advanced things you can do with Runnables:
```golang
type recursive struct{}

// Run runs a recursive job
func (g recursive) Run(job hive.Job, run hive.RunFunc) (interface{}, error) {
	fmt.Println("doing job:", job.String())

	if job.String() == "first" {
		return run(hive.NewJob("recursive", "second")), nil
	} else if job.String() == "second" {
		return run(hive.NewJob("recursive", "last")), nil
	}

	return fmt.Sprintf("finished %s", job.String()), nil
}
```
The `hive.RunFunc` that you see there is a way for your Runnable to, well, run more things!

Calling the `RunFunc` will schedule another job to be executed and give you a `Result`. If you return a `Result` from `Run`, then the caller will recursively recieve that `Result` when they call `Then()`!

For example:
```golang
r := h.Do(h.Job("recursive", "first"))

res, err := r.Then()
if err != nil {
	log.Fatal(err)
}

fmt.Println("done!", res.(string))
```
Will cause this output:
```
doing job: first
doing job: second
doing job: last
done! finished last
```
Think about that for a minute, and let it sink in, it can be quite powerful!

There are also some shortcuts to make working with Hive a bit easier:
```golang
type input struct {
	First, Second int
}

type math struct{}

// Run runs a math job
func (g math) Run(job hive.Job, run hive.RunFunc) (interface{}, error) {
	in := job.Data().(input)

	return in.First + in.Second, nil
}
```
```golang
doMath := h.Handle("math", math{})

for i := 1; i < 10; i++ {
	equals, _ := doMath(input{i, i * 3}).ThenInt()
	fmt.Println("result", equals)
}
```
The `Handle` function returns an optional helper function. Instead of passing a job name and full `Job` into `h.Do`, you can use the helper function to instead just pass the input data for the job, and you receive a `Result` as normal. `doMath`!

More to come, including better performance, library stability, etc. Cheers!