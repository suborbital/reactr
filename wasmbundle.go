package hive

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/suborbital/hivew/hivew/util"
)

// HandleBundle loads a .wasm.zip file into the hive instance
func (h *Hive) HandleBundle(path string) error {
	if !strings.HasSuffix(path, ".wasm.zip") {
		return fmt.Errorf("cannot load bundle %s, does not have .wasm.zip extension", filepath.Base(path))
	}

	bundle, err := util.ReadBundle(path)
	if err != nil {
		return errors.Wrap(err, "failed to ReadBundle")
	}

	for i, r := range bundle.Runnables {
		runner := newWasmFromRaw(&bundle.Runnables[i])

		jobName := strings.Replace(r.Name, ".wasm", "", -1)
		h.Handle(jobName, runner)
	}

	return nil
}
