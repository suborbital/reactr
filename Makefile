test:
	go test -v --count=1 -p=1 ./...

testdata:
	subo build ./rwasm/testdata/ --bundle --native

crate/check:
	cargo publish --manifest-path ./api/rust/suborbital/Cargo.toml --target=wasm32-wasi --dry-run

crate/publish:
	cargo publish --manifest-path ./api/rust/suborbital/Cargo.toml --target=wasm32-wasi

deps:
	go get -u -d ./...

.PHONY: test testdata crate/check crate/publish deps