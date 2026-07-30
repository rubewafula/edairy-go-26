.PHONY: build run test

build:
	go build -o bin/server cmd/api/main.go

run:
	go run cmd/api/main.go

test:
	go test ./...
