.PHONY: build test

build:
	go build -o conduit-connector-stripe cmd/stripe/main.go

test:
	go test $(GOTEST_FLAGS) -race ./...

lint:
	golangci-lint run -c .golangci.yml --go=1.18

dep:
	go mod download
	go mod tidy

mockgen:
	mockgen -package mock -source source/iterator/interface.go -destination source/iterator/mock/stripe.go