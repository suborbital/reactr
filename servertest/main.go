package main

import (
	"errors"
	"log"

	"github.com/suborbital/hive/hive"
	"github.com/suborbital/vektor/vk"
)

func main() {
	h := hive.New()

	h.Handle("generic", generic{})

	server := h.Server(vk.UseInsecureHTTP(8080), vk.UseAppName("hivetest"))

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

type generic struct{}

// Run runs a generic job
func (g generic) Run(job hive.Job, run hive.RunFunc) (interface{}, error) {
	if string(job.Bytes()) == "first" {
		return run(hive.NewJob("generic", []byte("second"))), nil
	} else if string(job.Bytes()) == "second" {
		return run(hive.NewJob("generic", []byte("last"))), nil
	}

	if string(job.Bytes()) == "error" {
		return nil, errors.New("bad!!")
	}

	return job.Bytes(), nil
}

func (g generic) OnStart() error {
	return nil
}
