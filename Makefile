TAGS ?= ""
GO_BIN ?= "go"

install:
	$(GO_BIN) install -tags ${TAGS} -v ./.
	make tidy

tidy:
	$(GO_BIN) mod tidy

build:
	$(GO_BIN) build -v .
	make tidy

test:
	$(GO_BIN) test -cover -race -tags ${TAGS} ./...
	make tidy

lint:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run --enable-all
	make tidy
