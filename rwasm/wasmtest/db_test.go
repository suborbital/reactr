package wasmtest

import (
	"fmt"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/suborbital/reactr/rcap"
	"github.com/suborbital/reactr/rt"
	"github.com/suborbital/reactr/rwasm"
	"github.com/suborbital/vektor/vlog"
)

func TestDBInsertQuery(t *testing.T) {
	dbConnString, exists := os.LookupEnv("REACTR_DB_CONN_STRING")
	if !exists {
		t.Skip("skipping as conn string env var not set")
	}

	config := rcap.DefaultConfigWithDB(vlog.Default(), dbConnString)
	r := rt.NewWithConfig(config)

	doWasm := r.Register("rs-dbtest", rwasm.NewRunner("../testdata/rs-dbtest/rs-dbtest.wasm"))

	res, err := doWasm(nil).Then()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to doWasm"))
	}

	fmt.Println(res)
}
