test:
	go test -v --count=1 ./...

wasm:
	wasm-pack build
	cp ./pkg/wasm_runner_bg.wasm ./wasm/

.PHONY: wasm