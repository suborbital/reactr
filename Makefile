packages = $(shell go list ./... | grep -v github.com/suborbital/reactr/api/tinygo/runnable)

test:
	go test -v --count=1 -p=1 $(packages)

test/wasmtime:
	go test --tags wasmtime -v --count=1 -p=1 $(packages)

test/wasmedge:
	go test --tags wasmedge -v -count=1 -p=1 $(packages)

test/multi: test test/wasmtime test/wasmedge

testdata:
	subo build ./engine/testdata/ --native

testdata/docker:
	subo build ./engine/testdata/

testdata/docker/dev:
	subo build ./engine/testdata/ --builder-tag dev --mountpath $(PWD)

crate/publish:
	cargo publish --manifest-path ./api/rust/codegen/Cargo.toml --target=wasm32-wasi
	cargo publish --manifest-path ./api/rust/core/Cargo.toml --target=wasm32-wasi

npm/publish:
	npm publish ./api/assemblyscript
	npm publish ./api/typescript

deps:
	go get -u -d ./...

deps/wasmedge:
	wget -qO- https://raw.githubusercontent.com/WasmEdge/WasmEdge/master/utils/install.sh | bash -s -- -v 0.9.1

.PHONY: test test/wasmtime test/wasmedge test/multi testdata crate/publish deps deps/wasmedge
