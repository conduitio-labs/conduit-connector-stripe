VERSION=$(shell git describe --tags --dirty --always)

.PHONY: build
build:
	go build -ldflags "-X 'github.com/conduitio-labs/conduit-connector-stripe.version=${VERSION}'" -o conduit-connector-stripe cmd/connector/main.go

.PHONY: test
test:
	go test $(GOTEST_FLAGS) -race ./...

.PHONY: generate
generate:
	go generate ./...

.PHONY: lint
lint:
	golangci-lint run -v

.PHONY: dep
dep:
	go mod download
	go mod tidy

.PHONY: mockgen
mockgen:
	mockgen -package mock -source source/source.go -destination source/mock/source.go
	mockgen -package mock -source source/iterator/iterator.go -destination source/iterator/mock/iterator.go

.PHONY: fmt
fmt:
	gofumpt -l -w .

.PHONY: install-tools
install-tools:
	@echo Installing tools from tools.go
	@go list -e -f '{{ join .Imports "\n" }}' tools.go | xargs -tI % go install %
	@go mod tidy