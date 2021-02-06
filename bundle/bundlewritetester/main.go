package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/suborbital/reactr/bundle"
	"github.com/suborbital/reactr/directive"
)

func main() {
	files := []os.File{}
	for _, filename := range []string{"fetch.wasm", "log_example.wasm", "example.wasm"} {
		path := filepath.Join("./", "wasm", "testdata", filename)

		file, err := os.Open(path)
		if err != nil {
			log.Fatal("failed to open file", err)
		}

		files = append(files, *file)
	}

	directive := &directive.Directive{
		Identifier:  "dev.suborbital.appname",
		AppVersion:  "v0.1.1",
		AtmoVersion: "v0.0.6",
		Runnables: []directive.Runnable{
			{
				Name:      "fetch",
				Namespace: "default",
			},
			{
				Name:      "log_example",
				Namespace: "default",
			},
			{
				Name:      "example",
				Namespace: "default",
			},
		},
		Handlers: []directive.Handler{
			{
				Input: directive.Input{
					Type:     directive.InputTypeRequest,
					Method:   "GET",
					Resource: "/api/v1/user",
				},
				Steps: []directive.Executable{
					{
						Group: []directive.CallableFn{
							{
								Fn: "fetch",
								As: "ghData",
							},
							{
								Fn: "log_example",
							},
						},
					},
					{
						CallableFn: directive.CallableFn{
							Fn: "example",
							With: []string{
								"data: ghData",
							},
						},
					},
				},
				Response: "ghData",
			},
		},
	}

	if err := directive.Validate(); err != nil {
		log.Fatal("failed to validate directive: ", err)
	}

	if err := bundle.Write(directive, files, "./runnables.wasm.zip"); err != nil {
		log.Fatal("failed to WriteBundle", err)
	}

	fmt.Println("done ✨")
}
