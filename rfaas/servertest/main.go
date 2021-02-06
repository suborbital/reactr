package main

import (
	"errors"
	"log"

	"github.com/suborbital/reactr/rfaas"
	"github.com/suborbital/reactr/rt"
	"github.com/suborbital/vektor/vk"
)

func main() {
	server := rfaas.New(vk.UseInsecureHTTP(8080), vk.UseAppName("rfaas test"))

	server.Handle("generic", generic{})

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

type generic struct{}

// Run runs a generic job
func (g generic) Run(job rt.Job, ctx *rt.Ctx) (interface{}, error) {
	if string(job.Bytes()) == "first" {
		return ctx.Do(rt.NewJob("generic", []byte("second"))), nil
	} else if string(job.Bytes()) == "second" {
		return ctx.Do(rt.NewJob("generic", []byte("last"))), nil
	}

	if string(job.Bytes()) == "error" {
		return nil, errors.New("bad")
	}

	return job.Bytes(), nil
}

func (g generic) OnChange(change rt.ChangeEvent) error {
	return nil
}
