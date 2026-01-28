// Copyright © 2026 Attestant Limited.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package funcopts provides an analyzer that suggests functional options pattern for service types.
package funcopts

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const (
	analyzerName = "attgo_func_opts"
	doc          = `suggests functional options pattern for service types

Service types should use the functional options pattern for configuration.
This provides a clean, extensible API for construction.

Bad:
    func NewService(log Logger, db DB, timeout time.Duration) *Service

Good:
    type Option func(*Service)

    func WithLogger(log Logger) Option {
        return func(s *Service) { s.log = log }
    }

    func New(opts ...Option) *Service`
)

// Analyzer is the functional options analyzer.
var Analyzer = &analysis.Analyzer{
	Name: analyzerName,
	Doc:  doc,
	Run:  run,
}

// serviceTypeSuffixes are suffixes that identify service types.
var serviceTypeSuffixes = []string{
	"Service",
	"Manager",
	"Handler",
	"Controller",
	"Provider",
	"Client",
	"Server",
}

func run(pass *analysis.Pass) (any, error) {
	// Collect service types.
	serviceTypes := make(map[string]bool)

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				// Check if it's a struct type with service-like name.
				if _, isStruct := typeSpec.Type.(*ast.StructType); !isStruct {
					continue
				}

				if isServiceTypeName(typeSpec.Name.Name) {
					serviceTypes[typeSpec.Name.Name] = true
				}
			}
		}
	}

	// Check constructor functions.
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			// Look for constructor functions (New..., Create...).
			if funcDecl.Recv != nil {
				continue // Skip methods.
			}

			name := funcDecl.Name.Name
			if !strings.HasPrefix(name, "New") && !strings.HasPrefix(name, "Create") {
				continue
			}

			// Check if returns a service type.
			returnType := getReturnTypeName(pass, funcDecl)
			if returnType == "" || !serviceTypes[returnType] {
				continue
			}

			// Check parameters - warn if more than 2 non-context parameters.
			if shouldSuggestFuncOpts(funcDecl) {
				pass.Reportf(funcDecl.Name.Pos(),
					"constructor %q has many parameters; consider using functional options pattern",
					name)
			}
		}
	}

	return nil, nil
}

// isServiceTypeName checks if a type name looks like a service.
func isServiceTypeName(name string) bool {
	for _, suffix := range serviceTypeSuffixes {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}

	return false
}

// getReturnTypeName extracts the type name from the function's return type.
func getReturnTypeName(pass *analysis.Pass, fn *ast.FuncDecl) string {
	if fn.Type.Results == nil || len(fn.Type.Results.List) == 0 {
		return ""
	}

	// Look for pointer to struct or struct.
	for _, result := range fn.Type.Results.List {
		switch t := result.Type.(type) {
		case *ast.StarExpr:
			if ident, ok := t.X.(*ast.Ident); ok {
				return ident.Name
			}
		case *ast.Ident:
			return t.Name
		}
	}

	return ""
}

// shouldSuggestFuncOpts determines if functional options should be suggested.
func shouldSuggestFuncOpts(fn *ast.FuncDecl) bool {
	if fn.Type.Params == nil {
		return false
	}

	nonContextParams := 0

	for _, param := range fn.Type.Params.List {
		// Check if this is a context parameter.
		if isContextParam(param) {
			continue
		}

		// Check if this is already a variadic options parameter.
		if isOptionsParam(param) {
			return false // Already using func opts pattern.
		}

		// Count non-context parameters (accounting for multiple names).
		names := len(param.Names)
		if names == 0 {
			names = 1
		}

		nonContextParams += names
	}

	// Suggest func opts if more than 3 non-context parameters.
	return nonContextParams > 3
}

// isContextParam checks if a parameter is a context.Context.
func isContextParam(param *ast.Field) bool {
	sel, ok := param.Type.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "context" && sel.Sel.Name == "Context"
}

// isOptionsParam checks if a parameter looks like a functional option.
func isOptionsParam(param *ast.Field) bool {
	// Check for variadic.
	ellipsis, ok := param.Type.(*ast.Ellipsis)
	if !ok {
		return false
	}

	// Check if the element type is a function or named Option type.
	switch t := ellipsis.Elt.(type) {
	case *ast.Ident:
		name := t.Name

		return strings.HasSuffix(name, "Option") || strings.HasSuffix(name, "Opt")
	case *ast.FuncType:
		return true
	}

	return false
}

// Ensure types package is used for type info.
var _ types.Type
