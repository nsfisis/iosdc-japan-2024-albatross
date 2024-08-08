.PHONY: all
all: check lint

.PHONY: check
check:
	go build -o /dev/null ./...

.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run
