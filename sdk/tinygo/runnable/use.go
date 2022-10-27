//go:build tinygo.wasm

package runnable

import (
	"github.com/suborbital/sdk/api/tinygo/runnable/internal/ffi"
	"github.com/suborbital/sdk/api/tinygo/runnable/runnable"
)

func Use(runnable runnable.Runnable) {
	ffi.Use(runnable)
}
