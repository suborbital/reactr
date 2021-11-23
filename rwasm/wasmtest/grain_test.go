package wasmtest

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/suborbital/reactr/rcap"
	"github.com/suborbital/reactr/request"
	"github.com/suborbital/reactr/rt"
	"github.com/suborbital/reactr/rwasm"
)

func TestGrainEcho(t *testing.T) {
	r := rt.New()

	// test a WASM module that is loaded directly instead of through the bundle
	doWasm := r.Register("grain-echo", rwasm.NewRunner("../testdata/grain-echo/grain-echo.wasm"))

	res, err := doWasm("from Grain!").Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then"))
		return
	}

	fmt.Println(string(res.([]byte)))

	if string(res.([]byte)) != "hello, from Grain!" {
		t.Error("grain-echo failed, got:", string(res.([]byte)))
	}
}

func TestGrainFetch(t *testing.T) {
	r := rt.New()

	// test a WASM module that is loaded directly instead of through the bundle
	doWasm := r.Register("grain-fetch", rwasm.NewRunner("../testdata/grain-fetch/grain-fetch.wasm"))

	res, err := doWasm("https://1password.com").Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then"))
		return
	}

	if string(res.([]byte)[:100]) != "<!doctype html><html lang=en data-language-url=/><head><meta charset=utf-8><meta name=viewport conte" {
		t.Error("grain-fetch failed, got:", string(res.([]byte)[:100]))
	}
}

func TestGrainGraphql(t *testing.T) {
	// bail out if GitHub auth is not set up (i.e. in Travis)
	if _, ok := os.LookupEnv("GITHUB_TOKEN"); !ok {
		return
	}

	r := rt.New()

	caps := r.DefaultCaps()
	caps.Auth = rcap.DefaultAuthProvider(rcap.AuthConfig{
		Enabled: true,
		Headers: map[string]rcap.AuthHeader{
			"api.github.com": {
				HeaderType: "bearer",
				Value:      "env(GITHUB_TOKEN)",
			},
		},
	})

	// test a WASM module that is loaded directly instead of through the bundle
	r.RegisterWithCaps("grain-graphql", rwasm.NewRunner("../testdata/grain-graphql/grain-graphql.wasm"), caps)

	res, err := r.Do(rt.NewJob("grain-graphql", nil)).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then"))
		return
	}

	if string(res.([]byte)) != `{"data":{"repository":{"name":"reactr","nameWithOwner":"suborbital/reactr"}}}` {
		t.Error("grain-graphql failed, got:", string(res.([]byte)))
	}
}

func TestGrainLargeData(t *testing.T) {
	r := rt.New()

	// test a WASM module that is loaded directly instead of through the bundle
	doWasm := r.Register("grain-echo", rwasm.NewRunner("../testdata/grain-echo/grain-echo.wasm"))

	res, err := doWasm(largeInput).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then"))
		return
	}

	if string(res.([]byte)) != "hello, "+largeInput {
		t.Error("grain-test failed, got:", string(res.([]byte)))
	}
}

func TestGrainRunnerWithRequest(t *testing.T) {
	r := rt.New()

	doWasm := r.Register("wasm", rwasm.NewRunner("../testdata/grain-req/grain-req.wasm"))

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

	res, err := doWasm(reqJSON).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Then"))
		return
	}

	resp := &request.CoordinatedResponse{}
	if err := json.Unmarshal(res.([]byte), resp); err != nil {
		t.Error("failed to Unmarshal response")
	}

	if string(resp.Output) != "hello what is up" {
		t.Error(fmt.Errorf("expected 'hello, what is up', got %s", string(res.([]byte))))
	}
}
