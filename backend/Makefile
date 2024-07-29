.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: check
check:
	go build -o /dev/null ./...

.PHONY: lint
lint:
	go vet ./...

.PHONY: oapi-codegen
oapi-codegen:
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config oapi-codegen.yaml ../openapi.yaml

.PHONY: sqlc
sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc generate
