package main

import (
	"fmt"
	"log"

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

	return in.First + in.Second, nil
}
