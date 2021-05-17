package main

import (
	"fmt"
	"time"

	"github.com/suborbital/reactr/rt"
)

func main() {
	r := rt.New()

	r.Register("print", &printJob{})

	r.Do(rt.NewJob("print", "start")).Discard()

	r.Schedule(rt.Every(5, func() rt.Job {
		return rt.NewJob("print", "every 5")
	}))

	r.Schedule(rt.After(7, func() rt.Job {
		return rt.NewJob("print", "after 7")
	}))

	time.Sleep(time.Second * 25)
}

type printJob struct{}

func (p *printJob) Run(job rt.Job, ctx *rt.Ctx) (interface{}, error) {
	fmt.Println(job.String())
	return nil, nil
}

func (p *printJob) OnChange(c rt.ChangeEvent) error { return nil }
