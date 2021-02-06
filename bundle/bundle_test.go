package bundle

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
)

func TestRead(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Error(errors.Wrap(err, "failed to get CWD"))
	}

	bundle, err := Read(filepath.Join(cwd, "../rwasm/testdata/runnables.wasm.zip"))
	if err != nil {
		t.Error(errors.Wrap(err, "failed to Read"))
		return
	}

	if len(bundle.Runnables) == 0 {
		t.Error("bundle had 0 runnables")
		return
	}

	hasDefault := false
	for _, r := range bundle.Runnables {
		if r.Name == "hello-echo.wasm" {
			hasDefault = true
		}
	}

	if !hasDefault {
		t.Error("hello-echo.wasm runnable not found in bundle")
	}
}
