package main

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/suborbital/hive"
)

func main() {
	h := hive.New()

	h.Handle("generic", generic{})

	r := h.Do(h.Job("generic", "first"))

	res, err := r.Then()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("done!", res.(string))

	doMath := h.Handle("math", math{})

	for i := 1; i < 10; i++ {
		equals, _ := doMath(input{i, i * 3}).ThenInt()
		fmt.Println("result", equals)
	}

	grp := hive.NewGroup()
	grp.Add(doMath(input{5, 6}))
	grp.Add(doMath(input{7, 8}))
	grp.Add(doMath(input{9, 10}))
	if err := grp.Wait(); err != nil {
		log.Fatal(errors.Wrap(err, "failed to Wait"))
	} else {
		fmt.Println("all good!")
	}

	doGrp := h.Handle("group", groupWork{})
	if _, err := doGrp(nil).Then(); err != nil {
		log.Fatal(errors.Wrap(err, "failed to doGrp"))
	}
}

type generic struct{}

// Run runs a generic job
func (g generic) Run(job hive.Job, run hive.RunFunc) (interface{}, error) {
	fmt.Println("doing job:", job.String())

	if job.String() == "first" {
		return run(hive.NewJob("generic", "second")), nil
	} else if job.String() == "second" {
		return run(hive.NewJob("generic", "last")), nil
	}

	return fmt.Sprintf("finished %s", job.String()), nil
}

type input struct {
	First, Second int
}

type math struct{}

// Run runs a math job
func (g math) Run(job hive.Job, run hive.RunFunc) (interface{}, error) {
	in := job.Data().(input)

	fmt.Println("adding", in.First, "+", in.Second)

	return in.First + in.Second, nil
}

type groupWork struct{}

// Run runs a groupWork job
func (g groupWork) Run(job hive.Job, run hive.RunFunc) (interface{}, error) {
	grp := hive.NewGroup()

	grp.Add(run(hive.NewJob("generic", "first")))
	grp.Add(run(hive.NewJob("generic", "group work")))
	grp.Add(run(hive.NewJob("generic", "group work")))
	grp.Add(run(hive.NewJob("generic", "group work")))
	grp.Add(run(hive.NewJob("generic", "group work")))

	return grp, nil
}
