package wasmtest

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/suborbital/reactr/rcap"
	"github.com/suborbital/reactr/rt"
	"github.com/suborbital/reactr/rwasm"
	"github.com/suborbital/vektor/vlog"
)

func TestDBInsertQuery(t *testing.T) {
	config := rcap.DefaultConfigWithDB(vlog.Default(), "bvzxim39dfti:pscale_pw_7gfxZr0DqAedpAkhJpMfNX5wXKS2eDgZovmzbBoxnns@tcp(ww5mgqwa0v0z.us-east-4.psdb.cloud)/suborbital-compute-network?tls=true")
	r := rt.NewWithConfig(config)

	doWasm := r.Register("rs-dbtest", rwasm.NewRunner("../testdata/rs-dbtest/rs-dbtest.wasm"))

	res, err := doWasm(nil).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to doWasm"))
	}

	fmt.Println(res)
}
