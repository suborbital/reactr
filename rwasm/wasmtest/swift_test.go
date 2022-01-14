package wasmtest

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/suborbital/reactr/rcap"
	"github.com/suborbital/reactr/request"
	"github.com/suborbital/reactr/rt"
	"github.com/suborbital/reactr/rwasm"
)

func TestWasmRunnerWithFetchSwift(t *testing.T) {
	r := rt.New()
	r.Register("fetch-swift", rwasm.NewRunner("../testdata/fetch-swift/fetch-swift.wasm"))

	job := rt.NewJob("fetch-swift", "https://1password.com")

	res, err := r.Do(job).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then"))
		return
	}

	if string(res.([]byte))[:100] != "<!doctype html><html lang=en data-language-url=/><head><meta charset=utf-8><meta name=viewport conte" {
		t.Error(fmt.Errorf("expected 1password.com HTML, got %q", string(res.([]byte))[:100]))
	}
}

func TestWasmRunnerEchoSwift(t *testing.T) {
	r := rt.New()
	r.Register("hello-swift", rwasm.NewRunner("../testdata/hello-swift/hello-swift.wasm"))

	job := rt.NewJob("hello-swift", "Connor")

	res, err := r.Do(job).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then"))
		return
	}

	if string(res.([]byte)) != "hello Connor" {
		t.Error(fmt.Errorf("hello Connor, got %s", string(res.([]byte))))
	}
}

func TestWasmRunnerSwift(t *testing.T) {
	body := testBody{
		Username: "cohix",
	}

	bodyJSON, _ := json.Marshal(body)

	req := &request.CoordinatedRequest{
		Method: "GET",
		URL:    "/hello/world",
		ID:     uuid.New().String(),
		Body:   bodyJSON,
		State: map[string][]byte{
			"hello": []byte("what is up"),
		},
	}

	reqJSON, err := req.ToJSON()
	if err != nil {
		t.Error("failed to ToJSON", err)
	}

	r := rt.New()
	r.Register("swift-log", rwasm.NewRunner("../testdata/swift-log/swift-log.wasm"))

	job := rt.NewJob("swift-log", reqJSON)

	res, err := r.Do(job).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then"))
		return
	}

	resp := res.(*request.CoordinatedResponse)

	if string(resp.Output) != "hello what is up" {
		t.Error(fmt.Errorf("expected 'hello, what is up', got %s", string(res.([]byte))))
	}
}

func TestWasmFileGetStaticSwift(t *testing.T) {
	config := rcap.DefaultCapabilityConfig()
	config.File = fileConfig

	r, _ := rt.NewWithConfig(config)
	r.Register("get-static-swift", rwasm.NewRunner("../testdata/get-static-swift/get-static-swift.wasm"))

	getJob := rt.NewJob("get-static-swift", "")

	res, err := r.Do(getJob).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Do get-static job"))
		return
	}

	result := string(res.([]byte))

	expected := "# Hello, World\n\nContents are very important"

	if result != expected {
		t.Error("failed, got:\n", result, "\nexpeted:\n", expected)
	}
}
