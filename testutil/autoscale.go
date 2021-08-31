package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/suborbital/reactr/rt"
	"github.com/suborbital/reactr/rwasm"
	"golang.org/x/crypto/pbkdf2"
)

func main() {
	start := time.Now()

	r := rt.New()

	r.Register("pbkdf2", &pbkdf2Job{}, rt.PoolSize(1), rt.Autoscale(0))
	r.Register("fetch", rwasm.NewRunner("./rwasm/testdata/as-fetch/as-fetch.wasm"), rt.PoolSize(1), rt.Autoscale(0))

	group := rt.NewGroup()

	go func() {
		for i := 0; i < 1000; i++ {
			group.Add(r.Do(rt.NewJob("pbkdf2", []byte("someinputtobehashed"))))
			group.Add(r.Do(rt.NewJob("fetch", "https://google.com")))
		}

		if err := group.Wait(); err != nil {
			log.Fatal(err)
		}

		duration := time.Since(start)

		fmt.Println("done!", duration.Seconds(), "s")
	}()

	go func() {
		for i := 0; i < 100; i++ {
			metrics := r.Metrics()
			metricsJSON, _ := json.MarshalIndent(metrics, "", "\t")
			fmt.Println(string(metricsJSON))

			time.Sleep(time.Second)
		}
	}()

	time.Sleep(time.Minute)
}

// an intentionally slow job to test scaling
type pbkdf2Job struct{}

func (p *pbkdf2Job) Run(job rt.Job, ctx *rt.Ctx) (interface{}, error) {
	// 100k rounds of PBKDF2 aught to slow things down
	pbkdf2.Key(job.Bytes(), []byte("neverdothis"), 100000, 32, sha256.New)
	return nil, nil
}

func (p *pbkdf2Job) OnChange(c rt.ChangeEvent) error { return nil }
