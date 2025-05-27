.DEFAULT_GOAL := build

.PHONY: fmt vet build
fmt:
	@go fmt ./...

vet: fmt
	@go vet ./...

build: vet
	@go build ./...

clean:
	@go mod tidy
	@go clean

test:
	@go test ./.. -vet=off
