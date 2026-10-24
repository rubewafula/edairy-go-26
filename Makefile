.PHONY: build run test

build:
	go build -o bin/server cmd/api/main.go

run:
	go run cmd/api/main.go

test:
	go test ./...

seed-route-centers:
	go run ./cmd/seed-route-centers

seed-credit-limits:
	go run ./cmd/seed-credit-limits
