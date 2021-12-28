//go:build go1.18

package rt

import (
	"sync"

	"golang.org/x/sync/errgroup"
)

// Group represents a group of job results
type Group[T any, R any] struct {
	results []*Result[R]
	sync.Mutex
}

// NewGroup creates a new Group
func NewGroup[T any, R any]() *Group[T, R] {
	g := &Group[T, R]{
		results: []*Result[R]{},
		Mutex:   sync.Mutex{},
	}

	return g
}

// Add adds a job result to the group
func (g *Group[T, R]) Add(result *Result[R]) {
	g.Lock()
	defer g.Unlock()

	if g.results == nil {
		g.results = []*Result[R]{}
	}

	g.results = append(g.results, result)
}

// Wait waits for all results to come in and returns an error if any arise
func (g *Group[T, R]) Wait() error {
	g.Lock()
	defer g.Unlock()

	wg := errgroup.Group{}

	for i := range g.results {
		res := g.results[i]

		wg.Go(func() error {
			_, err := res.Then()
			return err
		})
	}

	return wg.Wait()
}
