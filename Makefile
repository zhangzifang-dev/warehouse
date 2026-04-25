.PHONY: all build run test clean dev frontend-build

all: frontend-build build

frontend-build:
	cd web && npm install && npm run build

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
