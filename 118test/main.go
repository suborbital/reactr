//go:build go1.18

package main

import (
	"fmt"
	"log"
	"github.com/suborbital/reactr/rt"
)

func main() {
	r := rt.New()

	r.Register("118test", &runner[rt.Input, rt.Output]{})

	_, err := r.Do(rt.NewJob[rt.Input, rt.Output]("118test", "hello")).Then()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("done!")
}

type runner[T rt.Input, R rt.Output] struct{}

func (r *runner[T, R]) Run(rt.Job[T, R], *rt.Ctx) (R, error) {
	
	return strOut("hi"), nil
}

func (r *runner[T, R]) OnChange(_ rt.ChangeEvent) error { return nil }

func strOut[R rt.Output](in string) R {
	return in
}