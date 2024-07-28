package api

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// Work-around for this issue:
// https://stackoverflow.com/questions/70087465/echo-groups-not-working-with-openapi-generated-code-using-oapi-codegen
func GetSwaggerWithPrefix(prefix string) (*openapi3.T, error) {
	spec, err := GetSwagger()
	if err != nil {
		return nil, err
	}

	var prefixedPaths openapi3.Paths = openapi3.Paths{}
	for key, value := range spec.Paths.Map() {
		prefixedPaths.Set(prefix+key, value)
	}

	spec.Paths = &prefixedPaths
	return spec, nil
}
