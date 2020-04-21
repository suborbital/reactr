package hive

import (
	"log"
	"testing"

	"github.com/pkg/errors"
)

func TestHiveJobWithPool(t *testing.T) {
	h := New()

	doGeneric := h.Handle("generic", generic{}, PoolSize(3))

	grp := NewGroup()
	grp.Add(doGeneric("first"))
	grp.Add(doGeneric("first"))
	grp.Add(doGeneric("first"))

	if err := grp.Wait(); err != nil {
		log.Fatal(err)
	}
}

type badRunner struct{}

// Run runs a badRunner job
func (g badRunner) Run(job Job, run RunFunc) (interface{}, error) {
	return job.String(), nil
}

func (g badRunner) OnStart() error {
	return errors.New("fail")
}

func TestRunnerWithError(t *testing.T) {
	h := New()

	doBad := h.Handle("badRunner", badRunner{})

	_, err := doBad(nil).Then()
	if err == nil {
		t.Error("expected error, did not get one")
	}
}

func TestRunnerWithOptionsAndError(t *testing.T) {
	h := New()

	doBad := h.Handle("badRunner", badRunner{}, RetrySeconds(1), MaxRetries(1))

	_, err := doBad(nil).Then()
	if err == nil {
		t.Error("expected error, did not get one")
	}
}
