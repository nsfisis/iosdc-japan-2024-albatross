package gen

//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc generate
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config oapi-codegen.api-server.yaml ../../openapi/api-server.yaml
//go:generate go run ./api/handler_wrapper_gen.go -i ../api/generated.go -o ../api/handler_wrapper.go
//go:generate go run ./taskqueue/processor_wrapper_gen.go -i ../taskqueue/tasks.go -o ../taskqueue/processor_wrapper.go
