package engine

import (
	"github.com/suborbital/reactr/engine/api"
	"github.com/suborbital/reactr/engine/moduleref"
	"github.com/suborbital/reactr/scheduler"
)

// Engine is a Webassembly job scheduler with configurable host APIs
type Engine struct {
	*scheduler.Scheduler
	api api.HostAPI
}

// New creates a new Engine with the default API
func New() *Engine {
	return NewWithAPI(api.New())
}

// NewWithAPI creates a new Engine with the given API
func NewWithAPI(api api.HostAPI) *Engine {
	e := &Engine{
		Scheduler: scheduler.New(),
		api:       api,
	}

	return e
}

// Register registers a Wasm module by reference
func (e *Engine) Register(name string, ref *moduleref.WasmModuleRef, opts ...scheduler.Option) scheduler.JobFunc {
	runner := newRunnerFromRef(ref, e.api)

	return e.Scheduler.Register(name, runner)
}

// RegisterFromFile registers a Wasm module by reference
func (e *Engine) RegisterFromFile(name, filename string, opts ...scheduler.Option) scheduler.JobFunc {
	runner := newRunnerFromFile(filename, e.api)

	return e.Scheduler.Register(name, runner)
}
