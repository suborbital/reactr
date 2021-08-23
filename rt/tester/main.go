package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/suborbital/reactr/rt"
	"golang.org/x/crypto/pbkdf2"
)

func main() {
	start := time.Now()

	r := rt.New()

	r.Register("pbkdf2", &pbkdf2Job{}, rt.PoolSize(1), rt.Autoscale(0))

	go func() {
		for i := 0; i < 1000; i++ {
			r.Do(rt.NewJob("pbkdf2", []byte("someinputtobehashed")))
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			metrics := r.Metrics()
			metricsJSON, _ := json.MarshalIndent(metrics, "", "\t")
			fmt.Println(string(metricsJSON))

			if metrics.TotalThreadCount > 0 && metrics.TotalJobCount == 0 {
				duration := time.Since(start)

				fmt.Println("done!", duration.Seconds(), "s")

				os.Exit(0)
			}

			time.Sleep(time.Second)
		}
	}()

	time.Sleep(time.Second * 100)
}

// an intentionally slow job to test scaling
type pbkdf2Job struct{}

func (p *pbkdf2Job) Run(job rt.Job, ctx *rt.Ctx) (interface{}, error) {
	// 100k rounds of PBKDF2 aught to slow things down
	pbkdf2.Key(job.Bytes(), []byte("neverdothis"), 100000, 32, sha256.New)
	return nil, nil
}

func (p *pbkdf2Job) OnChange(c rt.ChangeEvent) error { return nil }
