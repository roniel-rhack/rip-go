.PHONY: build run test lint clean install

build:
	go build -o bin/rip-go ./cmd/rip-go

run:
	go run ./cmd/rip-go $(ARGS)

test:
	go test ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/

deps:
	go mod tidy

install:
	go install ./cmd/rip-go
