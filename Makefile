packages = $(shell go list ./... | grep -v github.com/suborbital/sdk/sdk/tinygo/runnable)

test:
	go test -v --count=1 -p=1 $(packages)

crate/publish:
	cargo publish --manifest-path ./sdk/rust/codegen/Cargo.toml --target=wasm32-wasi
	cargo publish --manifest-path ./sdk/rust/core/Cargo.toml --target=wasm32-wasi

npm/publish:
	npm publish ./sdk/assemblyscript
	npm publish ./sdk/typescript

deps:
	go get -u -d ./...

.PHONY: test crate/publish deps
