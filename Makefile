GOROOT := $(shell go env GOROOT)
SRCS := $(shell go list ./... | grep -v wasm)

.PHONY: wasm build test

wasm:
	cp $(GOROOT)/misc/wasm/wasm_exec.js docs
	GOOS=js GOARCH=wasm go build -o docs/repl_lib.wasm ./wasm/repl/
	GOOS=js GOARCH=wasm go build -o docs/vm_lib.wasm ./wasm/vm/

build:
	go build -v $(SRCS)

test:
	go test -v $(SRCS)
