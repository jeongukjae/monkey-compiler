.PHONY: build test

build:
	go build -v ./...

test:
	go test -v ./...
