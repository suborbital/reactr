test:
	go test -v --count=1 -p=1 ./...

wasm:
	wasm-pack build
	cp ./pkg/wasm_runner_bg.wasm ./wasm/

deps:
	go get -u -d ./...

.PHONY: wasm deps