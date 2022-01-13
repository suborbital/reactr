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

func TestWasmRunnerTinyGo(t *testing.T) {
	r := rt.New()

	// test a WASM module that is loaded directly instead of through the bundle
	doWasm := r.Register("wasm", rwasm.NewRunner("../testdata/tinygo-hello-echo/tinygo-hello-echo.wasm"))

	res, err := doWasm("world").Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then"))
		return
	}

	if string(res.([]byte)) != "Hello, world" {
		t.Errorf("expected Hello world got %q", string(res.([]byte)))
	}
}

func TestWasmFileGetStaticTinyGo(t *testing.T) {
	config := rcap.DefaultCapabilityConfig()
	config.File = fileConfig

	r, _ := rt.NewWithConfig(config)
	r.Register("tinygo-get-static", rwasm.NewRunner("../testdata/tinygo-get-static/tinygo-get-static.wasm"))

	getJob := rt.NewJob("tinygo-get-static", "")

	res, err := r.Do(getJob).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Do tinygo-get-static job"))
		return
	}

	result := string(res.([]byte))

	expected := "# Hello, World\n\nContents are very important"

	if result != expected {
		t.Error("failed, got:\n", result, "\nexpected:\n", expected)
	}
}

func TestGoURLQuery(t *testing.T) {
	r := rt.New()

	// using a Rust module
	doWasm := r.Register("wasm", rwasm.NewRunner("../testdata/tinygo-urlquery/tinygo-urlquery.wasm"))

	req := &request.CoordinatedRequest{
		Method: "GET",
		URL:    "/hello/world?message=whatsup",
		ID:     uuid.New().String(),
		Body:   []byte{},
	}

	res, err := doWasm(req).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then"))
		return
	}

	resp := &request.CoordinatedResponse{}
	if err := json.Unmarshal(res.([]byte), resp); err != nil {
		t.Error("failed to Unmarshal response")
	}

	if string(resp.Output) != "hello whatsup" {
		t.Error(fmt.Errorf("expected 'hello whatsup', got %s", string(resp.Output)))
	}
}
