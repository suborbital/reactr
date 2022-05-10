package wasmtest

import (
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

	// using a TinyGo module
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

	resp := res.(*request.CoordinatedResponse)

	if string(resp.Output) != "hello whatsup" {
		t.Error(fmt.Errorf("expected 'hello whatsup', got %s", string(resp.Output)))
	}
}

func TestGoContentType(t *testing.T) {
	req := &request.CoordinatedRequest{
		Method: "POST",
		URL:    "/hello/world",
		ID:     uuid.New().String(),
		Body:   []byte("world"),
	}

	reqJSON, err := req.ToJSON()
	if err != nil {
		t.Error("failed to ToJSON", err)
	}

	r := rt.New()
	r.Register("content-type", rwasm.NewRunner("../testdata/tinygo-resp/tinygo-resp.wasm"))

	job := rt.NewJob("content-type", reqJSON)

	res, err := r.Do(job).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then"))
		return
	}

	resp := res.(*request.CoordinatedResponse)

	if resp.RespHeaders["Content-Type"] != "application/json" {
		t.Error(fmt.Errorf("expected 'Content-Type: application/json', got %s", resp.RespHeaders["Content-Type"]))
	}

	if resp.RespHeaders["X-Reactr"] != string(req.Body) {
		t.Error(fmt.Errorf("expected 'X-Reactr: %s', got %s", string(req.Body), string(req.Body)))
	}
}
