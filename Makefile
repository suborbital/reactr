test:
	go test -v --count=1 ./...

wasm:
	wasm-pack build