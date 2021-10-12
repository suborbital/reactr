package wasmtest

import (
	"os"
	"testing"

	"github.com/suborbital/reactr/rcap"
	"github.com/suborbital/reactr/rt"
	"github.com/suborbital/reactr/rwasm"
)

func TestDisabledHTTP(t *testing.T) {
	config := rcap.DefaultCapabilityConfig()
	config.HTTP = &rcap.HTTPConfig{Enabled: false}

	r := rt.NewWithConfig(config)

	// test a WASM module that is loaded directly instead of through the bundle
	doWasm := r.Register("wasm", rwasm.NewRunner("../testdata/fetch/fetch.wasm"))

	_, err := doWasm("https://1password.com").Then()
	if err != nil {
		if err.Error() != `failed to execute Wasm Runnable: {"code":1,"message":"capability is not enabled"}` {
			t.Error("received incorrect error", err.Error())
		}
	} else {
		t.Error("runnable should have failed")
	}
}

func TestDisabledGraphQL(t *testing.T) {
	// bail out if GitHub auth is not set up (i.e. in Travis)
	// we want the Runnable to fail because graphql is disabled,
	// not because auth isn't set up correctly
	if _, ok := os.LookupEnv("GITHUB_TOKEN"); !ok {
		return
	}

	config := rcap.DefaultCapabilityConfig()
	config.GraphQL = &rcap.GraphQLConfig{Enabled: false}
	config.Auth = &rcap.AuthConfig{
		Enabled: true,
		Headers: map[string]rcap.AuthHeader{
			"api.github.com": {
				HeaderType: "bearer",
				Value:      "env(GITHUB_TOKEN)",
			},
		},
	}

	r := rt.NewWithConfig(config)

	r.Register("rs-graphql", rwasm.NewRunner("../testdata/rs-graphql/rs-graphql.wasm"))

	_, err := r.Do(rt.NewJob("rs-graphql", nil)).Then()
	if err != nil {
		if err.Error() != `failed to execute Wasm Runnable: {"code":1,"message":"capability is not enabled"}` {
			t.Error("incorrect error ", err.Error())
		}
	} else {
		t.Error("runnable should have produced an error")
	}
}
