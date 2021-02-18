package rwasm

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/suborbital/reactr/rt"
)

func BenchmarkRunnable(b *testing.B) {
	job := rt.NewJob("hello-echo", "my name is joe")

	for n := 0; n < b.N; n++ {
		_, err := sharedRT.Do(job).Then()
		if err != nil {
			b.Error(errors.Wrap(err, "failed to Then"))
		}
	}
}

func BenchmarkSwiftRunnable(b *testing.B) {
	job := rt.NewJob("hello-swift", "my name is joe")

	for n := 0; n < b.N; n++ {
		_, err := sharedRT.Do(job).Then()
		if err != nil {
			b.Error(errors.Wrap(err, "failed to Then"))
		}
	}
}

func BenchmarkRunnableFetch(b *testing.B) {
	job := rt.NewJob("bench-fetch", "https://google.com")

	for n := 0; n < b.N; n++ {
		_, err := sharedRT.Do(job).Then()
		if err != nil {
			b.Error(errors.Wrap(err, "failed to Then"))
		}
	}
}
