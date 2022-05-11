package wasmtest

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/suborbital/reactr/engine"
	"github.com/suborbital/reactr/engine/api"
	"github.com/suborbital/reactr/engine/capabilities"
	"github.com/suborbital/reactr/request"
	"github.com/suborbital/reactr/scheduler"
)

func TestWasmRunnerWithFetch(t *testing.T) {
	e := engine.New()

	// test a WASM module that is loaded directly instead of through the bundle
	doWasm := e.RegisterFromFile("wasm", "../testdata/fetch/fetch.wasm")

	res, err := doWasm("https://1password.com").Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then"))
		return
	}

	if len(res.([]byte)) < 100 {
		t.Errorf("expected 1password.com HTML, got %q", string(res.([]byte)))
	}

	if string(res.([]byte))[:100] != "{\"args\":{},\"data\":{\"message\":\"testing the echo!\"},\"files\":{},\"form\":{},\"headers\":{\"x-forwarded-proto" {
		t.Errorf("expected echo response, got %q", string(res.([]byte))[:100])
	}
}

func TestGraphQLRunner(t *testing.T) {
	// bail out if GitHub auth is not set up (i.e. in Travis)
	if _, ok := os.LookupEnv("GITHUB_TOKEN"); !ok {
		return
	}

	config := capabilities.DefaultCapabilityConfig()
	config.Auth = &capabilities.AuthConfig{
		Enabled: true,
		Headers: map[string]capabilities.AuthHeader{
			"api.github.com": {
				HeaderType: "bearer",
				Value:      "env(GITHUB_TOKEN)",
			},
		},
	}

	api, _ := api.NewWithConfig(config)

	e := engine.NewWithAPI(api)

	e.RegisterFromFile("rs-graphql", "../testdata/rs-graphql/rs-graphql.wasm")

	res, err := e.Do(scheduler.NewJob("rs-graphql", nil)).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then"))
		return
	}

	if string(res.([]byte)) != `{"data":{"repository":{"name":"reactr","nameWithOwner":"suborbital/reactr"}}}` {
		t.Error("as-graphql failed, got:", string(res.([]byte)))
	}
}

func TestWasmRunnerReturnError(t *testing.T) {
	e := engine.New()
	e.RegisterFromFile("return-err", "../testdata/return-err/return-err.wasm")

	job := scheduler.NewJob("return-err", "asdf")

	_, err := e.Do(job).Then()
	if err == nil {
		t.Error("expected error, got none")
		return
	}

	if runErr, ok := err.(scheduler.RunErr); !ok || runErr.Error() != `{"code":400,"message":"job failed"}` {
		t.Error("expected RunErr JSON, got", err.Error())
	}
}

func TestWasmRunnerWithRequest(t *testing.T) {
	e := engine.New()

	// using a Rust module
	doWasm := e.RegisterFromFile("wasm", "../testdata/log/log.wasm")

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

	res, err := doWasm(req).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then"))
		return
	}

	resp := res.(*request.CoordinatedResponse)

	if string(resp.Output) != "hello what is up" {
		t.Error(fmt.Errorf("expected 'hello, what is up', got %s", string(res.([]byte))))
	}
}

func TestRustURLQuery(t *testing.T) {
	e := engine.New()

	// using a Rust module
	doWasm := e.RegisterFromFile("wasm", "../testdata/rust-urlquery/rust-urlquery.wasm")

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

func TestWasmRunnerSetRequest(t *testing.T) {
	e := engine.New()

	// using a Rust module
	doWasm := e.RegisterFromFile("wasm", "../testdata/rs-reqset/rs-reqset.wasm")

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
		Headers: map[string]string{},
	}

	_, err := doWasm(req).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then"))
		return
	}

	if val, ok := req.Headers["X-REACTR-TEST"]; !ok {
		t.Error("header was not set correctly")
	} else if val != "test successful!" {
		t.Error(fmt.Errorf("expected 'test successful!', got %s", val))
	}
}

func TestEmptyRequestBodyJSON(t *testing.T) {
	e := engine.New()

	// using a Rust module
	doWasm := e.RegisterFromFile("wasm", "../testdata/log/log.wasm")

	req := &request.CoordinatedRequest{
		Method: "GET",
		URL:    "/hello/world",
		ID:     uuid.New().String(),
		Body:   []byte{},
		State: map[string][]byte{
			"hello": []byte("what is up"),
		},
	}

	reqJSON, err := req.ToJSON()
	if err != nil {
		t.Error("failed to ToJSON", err)
	}

	res, err := doWasm(reqJSON).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then"))
		return
	}

	resp := res.(*request.CoordinatedResponse)

	if string(resp.Output) != "hello what is up" {
		t.Error(fmt.Errorf("expected 'hello, what is up', got %s", string(resp.Output)))
	}
}

func TestContentType(t *testing.T) {
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

	e := engine.New()
	e.RegisterFromFile("content-type", "../testdata/content-type/content-type.wasm")

	job := scheduler.NewJob("content-type", reqJSON)

	res, err := e.Do(job).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then"))
		return
	}

	resp := res.(*request.CoordinatedResponse)

	if resp.RespHeaders["Content-Type"] != "application/json" {
		t.Error("unexpected ctype, actually is", resp.RespHeaders["Content-Type"])
	}
}

func TestWasmRunnerDataConversion(t *testing.T) {
	e := engine.New()

	doWasm := e.RegisterFromFile("wasm", "../testdata/hello-echo/hello-echo.wasm")

	res, err := doWasm("my name is joe").Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then"))
	}

	if string(res.([]byte)) != "hello my name is joe" {
		t.Error(fmt.Errorf("expected 'hello my name is joe', got %s", string(res.([]byte))))
	}
}

func TestWasmRunnerGroup(t *testing.T) {
	e := engine.New()

	doWasm := e.RegisterFromFile("wasm", "../testdata/hello-echo/hello-echo.wasm")

	grp := scheduler.NewGroup()
	for i := 0; i < 50000; i++ {
		grp.Add(doWasm([]byte(fmt.Sprintf("world %d", i))))
	}

	if err := grp.Wait(); err != nil {
		t.Error(err)
	}
}

func TestWasmBundle(t *testing.T) {
	e := engine.New()
	e.RegisterFromFile("hello-echo", "../testdata/hello-echo/hello-echo.wasm")

	res, err := e.Do(scheduler.NewJob("hello-echo", []byte("wasmWorker!"))).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "Then returned error"))
		return
	}

	if string(res.([]byte)) != "hello wasmWorker!" {
		t.Error(fmt.Errorf("expected 'hello wasmWorker!', got %s", string(res.([]byte))))
	}
}

func TestWasmLargeData(t *testing.T) {
	e := engine.New()

	doWasm := e.RegisterFromFile("wasm", "../testdata/hello-echo/hello-echo.wasm")

	res := doWasm([]byte(largeInput))

	result, err := res.Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then for large input"))
	}

	if len(string(result.([]byte))) < 64000 {
		t.Error(fmt.Errorf("large input job . too small, got %d", len(string(result.([]byte)))))
	}

	if string(result.([]byte)) != fmt.Sprintf("hello %s", largeInput) {
		t.Error(fmt.Errorf("large input result did not match"))
	}
}

func TestWasmLargeDataGroup(t *testing.T) {
	e := engine.New()

	doWasm := e.RegisterFromFile("wasm", "../testdata/hello-echo/hello-echo.wasm")

	grp := scheduler.NewGroup()
	for i := 0; i < 5000; i++ {
		grp.Add(doWasm([]byte(largeInput)))
	}

	if err := grp.Wait(); err != nil {
		t.Error("group returned an error")
	}
}

func TestWasmLargeDataGroupWithPool(t *testing.T) {
	e := engine.New()

	doWasm := e.RegisterFromFile("wasm", "../testdata/hello-echo/hello-echo.wasm", scheduler.PoolSize(5))

	grp := scheduler.NewGroup()
	for i := 0; i < 5000; i++ {
		grp.Add(doWasm([]byte(largeInput)))
	}

	if err := grp.Wait(); err != nil {
		t.Error("group returned an error")
	}
}

func TestWasmFileGetStatic(t *testing.T) {
	config := capabilities.DefaultCapabilityConfig()
	config.File = fileConfig

	api, _ := api.NewWithConfig(config)

	e := engine.NewWithAPI(api)

	e.RegisterFromFile("get-static", "../testdata/get-static/get-static.wasm")

	getJob := scheduler.NewJob("get-static", "important.md")

	res, err := e.Do(getJob).Then()
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
