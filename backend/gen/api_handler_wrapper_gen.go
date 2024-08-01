package main

import (
	"bytes"
	"flag"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"slices"
	"strings"
	"text/template"
)

func main() {
	inputFile := flag.String("i", "", "input file")
	outputFile := flag.String("o", "", "output file")
	flag.Parse()

	if inputFile == nil || *inputFile == "" || outputFile == nil || *outputFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Parse the input file
	fileSet := token.NewFileSet()
	parsedFile, err := parser.ParseFile(fileSet, *inputFile, nil, parser.SkipObjectResolution)
	if err != nil {
		panic(err)
	}

	// Find methods in StrictServerInterface
	var methods []string
	for _, decl := range parsedFile.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if typeSpec.Name.Name != "StrictServerInterface" {
				continue
			}
			interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				continue
			}
			for _, method := range interfaceType.Methods.List {
				if len(method.Names) != 0 {
					methods = append(methods, method.Names[0].Name)
				}
			}
		}
	}
	if len(methods) == 0 {
		panic("StrictServerInterface not found")
	}
	slices.Sort(methods)

	type TemplateParameter struct {
		Name              string
		RequiresLogin     bool
		RequiresAdminRole bool
	}
	templateParameters := make([]TemplateParameter, len(methods))
	for i, method := range methods {
		templateParameters[i] = TemplateParameter{
			Name:              method,
			RequiresLogin:     method != "PostLogin",
			RequiresAdminRole: strings.Contains(method, "Admin"),
		}
	}

	// Generate code.
	tmpl, err := template.New("code").Parse(templateText)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, templateParameters)
	if err != nil {
		panic(err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(*outputFile, formatted, 0644)
	if err != nil {
		panic(err)
	}
}

const templateText = `// Code generated by go generate; DO NOT EDIT.

package api

import (
	"context"
	"errors"
	"strings"

	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/auth"
	"github.com/nsfisis/iosdc-japan-2024-albatross/backend/db"
)

var _ StrictServerInterface = (*ApiHandlerWrapper)(nil)

type ApiHandlerWrapper struct {
	innerHandler ApiHandler
}

func NewHandler(queries *db.Queries, hubs GameHubsInterface) *ApiHandlerWrapper {
	return &ApiHandlerWrapper{
		innerHandler: ApiHandler{
			q:    queries,
			hubs: hubs,
		},
	}
}

func parseJWTClaimsFromAuthorizationHeader(authorization string) (*auth.JWTClaims, error) {
	const prefix = "Bearer "
	if !strings.HasPrefix(authorization, prefix) {
		return nil, errors.New("invalid authorization header")
	}
	token := authorization[len(prefix):]
	claims, err := auth.ParseJWT(token)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

{{ range . }}
	func (h *ApiHandlerWrapper) {{ .Name }}(ctx context.Context, request {{ .Name }}RequestObject) ({{ .Name }}ResponseObject, error) {
		{{ if .RequiresLogin -}}
			user, err := parseJWTClaimsFromAuthorizationHeader(request.Params.Authorization)
			if err != nil {
				return {{ .Name }}401JSONResponse{
					UnauthorizedJSONResponse: UnauthorizedJSONResponse{
						Message: "Unauthorized",
					},
				}, nil
			}
			{{ if .RequiresAdminRole -}}
				if !user.IsAdmin {
					return {{ .Name }}403JSONResponse{
						ForbiddenJSONResponse: ForbiddenJSONResponse{
							Message: "Forbidden",
						},
					}, nil
				}
			{{ end -}}
			return h.innerHandler.{{ .Name }}(ctx, request, user)
		{{ else -}}
			return h.innerHandler.{{ .Name }}(ctx, request)
		{{ end -}}
	}
{{ end }}
`
