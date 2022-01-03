GOROOT := $(shell go env GOROOT)
SRCS := $(shell go list ./... | grep -v wasm)

.PHONY: wasm build test

wasm:
	cp $(GOROOT)/misc/wasm/wasm_exec.js docs
	GOOS=js GOARCH=wasm go build -o docs/repl_lib.wasm wasm/main.go

build:
	go build -v ./...

test:
	go test -v ./...
