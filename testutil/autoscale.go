package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/suborbital/reactr/engine"
	"github.com/suborbital/reactr/scheduler"
	"golang.org/x/crypto/pbkdf2"
)

func main() {
	start := time.Now()

	e := engine.New()

	e.Scheduler.Register("pbkdf2", &pbkdf2Job{}, scheduler.PoolSize(1), scheduler.Autoscale(0))
	e.RegisterFromFile("fetch", "./engine/testdata/as-fetch/as-fetch.wasm", scheduler.PoolSize(1), scheduler.Autoscale(0))

	group := scheduler.NewGroup()

	go func() {
		for i := 0; i < 1000; i++ {
			group.Add(e.Do(scheduler.NewJob("pbkdf2", []byte("someinputtobehashed"))))
			group.Add(e.Do(scheduler.NewJob("fetch", "https://google.com")))
		}

		if err := group.Wait(); err != nil {
			log.Fatal(err)
		}

		duration := time.Since(start)

		fmt.Println("done!", duration.Seconds(), "s")
	}()

	go func() {
		for i := 0; i < 100; i++ {
			metrics := e.Metrics()
			metricsJSON, _ := json.MarshalIndent(metrics, "", "\t")
			fmt.Println(string(metricsJSON))

			time.Sleep(time.Second)
		}
	}()

	time.Sleep(time.Minute)
}

// an intentionally slow job to test scaling
type pbkdf2Job struct{}

func (p *pbkdf2Job) Run(job scheduler.Job, ctx *scheduler.Ctx) (interface{}, error) {
	// 100k rounds of PBKDF2 aught to slow things down
	pbkdf2.Key(job.Bytes(), []byte("neverdothis"), 100000, 32, sha256.New)
	return nil, nil
}

func (p *pbkdf2Job) OnChange(c scheduler.ChangeEvent) error { return nil }
