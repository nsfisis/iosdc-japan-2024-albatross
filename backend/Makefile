.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: check
check:
	go build -o /dev/null ./...

.PHONY: lint
lint:
	go vet ./...

.PHONY: gen
gen:
	go generate ./...
