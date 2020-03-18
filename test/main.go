package main

import (
	"fmt"
	"log"

	"github.com/suborbital/hive"
)

func main() {
	h := hive.New()

	h.Handle("generic", generic{})

	r := h.Do(hive.NewJob("generic", "first"))

	res, err := r.Then()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("done!", res.(string))
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
