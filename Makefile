packages = $(shell go list ./... | grep -v github.com/suborbital/reactr/api/tinygo/runnable)

test:
	go test -v --count=1 -p=1 $(packages)

test/multi: test
	go test --tags wasmtime -v --count=1 -p=1 $(packages)

testdata:
	subo build ./rwasm/testdata/ --native

testdata/docker:
	subo build ./rwasm/testdata/

testdata/docker/dev:
	subo build ./rwasm/testdata/ --builder-tag dev --mountpath $(PWD)

crate/check:
	cargo publish --manifest-path ./api/rust/codegen/Cargo.toml --target=wasm32-wasi --dry-run
	cargo publish --manifest-path ./api/rust/core/Cargo.toml --target=wasm32-wasi --dry-run

crate/publish:
	cargo publish --manifest-path ./api/rust/codegen/Cargo.toml --target=wasm32-wasi
	cargo publish --manifest-path ./api/rust/core/Cargo.toml --target=wasm32-wasi

npm/publish:
	npm publish ./api/assemblyscript

deps:
	go get -u -d ./...

.PHONY: test testdata crate/check crate/publish deps
