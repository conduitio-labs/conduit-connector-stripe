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
	golangci-lint run

.PHONY: dep
dep:
	go mod download
	go mod tidy

.PHONY: fmt
fmt:
	gofumpt -l -w .

.PHONY: install-tools
install-tools:
	@echo Installing tools from tools/go.mod
	@go list -modfile=tools/go.mod tool | xargs -I % go list -modfile=tools/go.mod -f "%@{{.Module.Version}}" % | xargs -tI % go install %
	@go mod tidy
