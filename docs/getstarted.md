# Get started with Hive 🚀

Once you've gotten the basics of Hive, follow along here to learn what makes it so powerful.

## Runnables pt. 2

There are some more complicated things you can do with Runnables:
```golang
type recursive struct{}

// Run runs a recursive job
func (r recursive) Run(job hive.Job, do hive.DoFunc) (interface{}, error) {
	fmt.Println("doing job:", job.String())

	if job.String() == "first" {
		return do(hive.NewJob("recursive", "second")), nil
	} else if job.String() == "second" {
		return do(hive.NewJob("recursive", "last")), nil
	}

	return fmt.Sprintf("finished %s", job.String()), nil
}

func (r recursive) OnChange(change ChangeEvent) error {
	return nil
}
```
The `hive.DoFunc` that you see there is a way for your Runnable to, well, run more things!

Calling the `DoFunc` will schedule another job to be executed and give you a `Result`. If you return a `Result` from `Run`, then the caller will recursively recieve that `Result` when they call `Then()`!

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
The ability to chain jobs is quite powerful!

You won't always need or care about a job's output, and in those cases, make sure to call `Discard()` on the result to allow the underlying resources to be deallocated!
```golang
h.Do(h.Job("recursive", "first")).Discard()
```

To to something asynchronously with the `Result` once it completes, call `ThenDo` on the result:
```golang
h.Do(h.Job("generic", "first")).ThenDo(func(res interface{}, err error) {
	if err != nil {
		// do something with the error
	}

	//do something with the result
})
```
`ThenDo` will return immediately, and provided callback will be run on a background goroutine. This is useful for handling results that don't need to be consumed by your main program execution.

### Groups

A hive `Group` is a set of `Result`s that belong together. If you're familiar with Go's `errgroup.Group{}`, it is similar. Adding results to a group will allow you to evaluate them all together at a later time.
```golang
grp := hive.NewGroup()

grp.Add(do(hive.NewJob("recursive", "first")))
grp.Add(do(hive.NewJob("generic", "group work")))
grp.Add(do(hive.NewJob("generic", "group work")))

if err := grp.Wait(); err != nil {
	log.Fatal(err)
}
```
Will print: 
```
doing job: first
doing job: group work
doing job: group work
doing job: second
doing job: last
```
As you can see, the "recursive" jobs from the `generic` runner get queued up after the two jobs that don't recurse.

Note that you cannot get result values from result groups, the error returned from `Wait()` will be the first error from any of the results in the group, if any. To get result values from a group of jobs, put them in an array and call `Then` on them individually.

**TIP** If you return a group from a Runnable's `Run`, calling `Then()` on the result will recursively call `Wait()` on the group and return the error to the original caller! You can easily chain jobs and job groups in various orders.

### Pools
Each `Runnable` that you register is given a worker to process their jobs. By default, each worker has one work thread processing jobs in sequence. If you want a particular worker to process more than one job concurrently, you can increase its `PoolSize`:
```golang
doGeneric := h.Handle("generic", generic{}, hive.PoolSize(3))

grp := hive.NewGroup()
grp.Add(doGeneric("first"))
grp.Add(doGeneric("second"))
grp.Add(doGeneric("random"))

if err := grp.Wait(); err != nil {
	log.Fatal(err)
}
```
Passing `PoolSize(3)` will spawn three work threads to process `generic` jobs.

### Timeouts
By default, if a job becomes stuck and is blocking execution, it will block forever. If you want to have a worker time out after a certain amount of seconds on a stuck job, pass `hive.TimeoutSeconds` to Handle:
``` golang
h := hive.New()

doTimeout := h.Handle("timeout", timeoutRunner{}, hive.TimeoutSeconds(3))
```
When `TimeoutSeconds` is set and a job executes for longer than the provided number of seconds, the worker will move on to the next job and `ErrJobTimeout` will be returned to the Result. The failed job will continue to execute in the background, but its result will be discarded.

### Advanced Runnables

The `Runnable` interface defines an `OnChange` function which gives the Runnable a chance to prepare itself for changes to the worker running it. For example, when a Runnable is registered with a pool size greater than 1, the Runnable may need to provision resources for itself to enable handling jobs concurrently, and `OnChange` will be called once each time a new worker starts up. Our [Wasm implementation](https://github.com/suborbital/hive-wasm/blob/master/wasm/wasmrunnable.go) is a good example of this. 

Most Runnables can return `nil` from this function, however returning an error will cause the worker start to be paused and retried until the required pool size has been acheived. The number of seconds between retries (default 3) and the maximum number of retries (default 5) can be configured when registering a Runnable:
```golang
doBad := h.Handle("badRunner", badRunner{}, hive.RetrySeconds(1), hive.MaxRetries(10))
```
Any error from a failed worker will be returned to the first job that is attempted for that Runnable.

### Shortcuts

There are also some shortcuts to make working with Hive a bit easier:
```golang
type input struct {
	First, Second int
}

type math struct{}

// Run runs a math job
func (g math) Run(job hive.Job, do hive.DoFunc) (interface{}, error) {
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

## Additional features

Hive can integrate with [Grav](https://github.com/suborbital/grav), which is the decentralized message bus developed as part of the Suborbital Development Platform. Read about the integration on [the grav documentation page.](./grav.md)

Hive provides the building blocks for scalable asynchronous systems. This should be everything you need to help you improve the performance of your application. When you are looking to take advantage of Hive's other features, check out its [FaaS](./faas.md) and [Wasm](./wasm.md) capabilities!