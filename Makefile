.PHONY: build run test clean dev

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server

test:
	go test -v ./...

clean:
	rm -rf bin/

dev:
	go run ./cmd/server
