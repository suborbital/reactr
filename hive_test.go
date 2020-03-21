package hive

import (
	"log"
	"testing"

	"github.com/pkg/errors"
)

type generic struct{}

// Run runs a generic job
func (g generic) Run(job Job, run RunFunc) (interface{}, error) {
	if job.String() == "first" {
		return run(NewJob("generic", "second")), nil
	} else if job.String() == "second" {
		return run(NewJob("generic", "last")), nil
	}

	return job.String(), nil
}

func TestHiveJob(t *testing.T) {
	h := New()

	h.Handle("generic", generic{})

	r := h.Do(h.Job("generic", "first"))

	res, err := r.Then()
	if err != nil {
		log.Fatal(err)
	}

	if res.(string) != "last" {
		t.Error("generic job failed, expected 'last', got", res.(string))
	}
}

type input struct {
	First, Second int
}

type math struct{}

// Run runs a math job
func (g math) Run(job Job, run RunFunc) (interface{}, error) {
	in := job.Data().(input)

	return in.First + in.Second, nil
}

func TestHiveJobHelperFunc(t *testing.T) {
	h := New()

	doMath := h.Handle("math", math{})

	for i := 1; i < 10; i++ {
		answer := i + i*3

		equals, _ := doMath(input{i, i * 3}).ThenInt()
		if equals != answer {
			t.Error("failed to get math right, expected", answer, "got", equals)
		}
	}
}

func TestHiveJobGroup(t *testing.T) {
	h := New()

	doMath := h.Handle("math", math{})

	grp := NewGroup()
	grp.Add(doMath(input{5, 6}))
	grp.Add(doMath(input{7, 8}))
	grp.Add(doMath(input{9, 10}))

	if err := grp.Wait(); err != nil {
		t.Error(errors.Wrap(err, "failed to grp.Wait"))
	}
}

type groupWork struct{}

// Run runs a groupWork job
func (g groupWork) Run(job Job, run RunFunc) (interface{}, error) {
	grp := NewGroup()

	grp.Add(run(NewJob("generic", "first")))
	grp.Add(run(NewJob("generic", "group work")))
	grp.Add(run(NewJob("generic", "group work")))
	grp.Add(run(NewJob("generic", "group work")))
	grp.Add(run(NewJob("generic", "group work")))

	return grp, nil
}

func TestHiveChainedGroup(t *testing.T) {
	h := New()

	h.Handle("generic", generic{})
	doGrp := h.Handle("group", groupWork{})

	if _, err := doGrp(nil).Then(); err != nil {
		t.Error(errors.Wrap(err, "failed to doGrp"))
	}
}
