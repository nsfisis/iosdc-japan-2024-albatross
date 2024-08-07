package gen

//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc generate
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config oapi-codegen.yaml ../../openapi.yaml
//go:generate go run ./api/handler_wrapper_gen.go -i ../api/generated.go -o ../api/handler_wrapper.go
