.PHONY: build test

VERSION=$(shell git describe --tags --dirty --always)

build:
	go build -ldflags "-X 'github.com/conduitio-labs/conduit-connector-stripe.version=${VERSION}'" -o conduit-connector-stripe cmd/connector/main.go

test:
	go test $(GOTEST_FLAGS) -race ./...

lint:
	golangci-lint run -c .golangci.yml --go=1.18

dep:
	go mod download
	go mod tidy

mockgen:
	mockgen -package mock -source source/source.go -destination source/mock/source.go
	mockgen -package mock -source source/iterator/iterator.go -destination source/iterator/mock/iterator.go
