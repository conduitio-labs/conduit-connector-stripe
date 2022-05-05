.PHONY: build test

build:
	go build -o conduit-connector-stripe cmd/stripe/main.go

test:
	go test $(GOTEST_FLAGS) -race ./...

lint:
	golangci-lint run -c .golangci.yml --go=1.18

mockgen:
	mockgen -package mock -source clients/http/interface.go -destination clients/http/mock/http.go