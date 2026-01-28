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

// Package enumiota provides an analyzer that enforces the iota pattern for enum types.
// Enum types should use uint64 with iota, not string constants.
package enumiota

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const (
	analyzerName = "attgo_enum_iota"
	doc          = `enforces iota pattern for enum types

Enum types should use uint64 (or another integer type) with iota, not string constants.
The string representation should be provided via a String() method.

Bad:
    type SANType string
    const (
        SANTypeDNS   SANType = "dns"
        SANTypeEmail SANType = "email"
    )

Good:
    type SANType uint64
    const (
        SANTypeUnknown SANType = iota
        SANTypeDNS
        SANTypeEmail
    )

    func (s SANType) String() string {
        return [...]string{"unknown", "dns", "email"}[s]
    }`
)

// NewAnalyzer creates a new enum-iota analyzer with the given enum type suffixes.
func NewAnalyzer(enumTypeSuffixes []string) *analysis.Analyzer {
	r := &runner{
		enumTypeSuffixes: enumTypeSuffixes,
	}

	return &analysis.Analyzer{
		Name: analyzerName,
		Doc:  doc,
		Run:  r.run,
	}
}

type runner struct {
	enumTypeSuffixes []string
}

func (r *runner) run(pass *analysis.Pass) (any, error) {
	// First pass: collect type definitions that look like enums (have enum-like suffixes).
	enumTypes := make(map[string]*ast.TypeSpec)

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				// Check if the type name has an enum-like suffix.
				if r.isEnumTypeName(typeSpec.Name.Name) {
					enumTypes[typeSpec.Name.Name] = typeSpec
				}
			}
		}
	}

	// Second pass: check const declarations that use these types.
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.CONST {
				continue
			}

			r.checkConstDecl(pass, genDecl, enumTypes)
		}
	}

	return nil, nil
}

// isEnumTypeName checks if a type name appears to be an enum type based on suffix.
func (r *runner) isEnumTypeName(name string) bool {
	for _, suffix := range r.enumTypeSuffixes {
		if strings.HasSuffix(name, suffix) {
			return true
		}
	}

	return false
}

// checkConstDecl checks a const declaration for string-based enum patterns.
func (r *runner) checkConstDecl(pass *analysis.Pass, genDecl *ast.GenDecl, enumTypes map[string]*ast.TypeSpec) {
	// Check each const spec.
	for _, spec := range genDecl.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}

		// Get the type of the const.
		if len(valueSpec.Names) == 0 {
			continue
		}

		obj := pass.TypesInfo.ObjectOf(valueSpec.Names[0])
		if obj == nil {
			continue
		}

		// Get the underlying type.
		constType := obj.Type()
		if constType == nil {
			continue
		}

		// Check if this const uses a named type.
		named, ok := constType.(*types.Named)
		if !ok {
			continue
		}

		typeName := named.Obj().Name()

		// Only check if it's one of our enum types.
		if _, isEnum := enumTypes[typeName]; !isEnum {
			continue
		}

		// Check if the underlying type is string.
		if isStringType(named.Underlying()) {
			// Check if this const has a string literal value.
			if hasStringLiteralValue(valueSpec) {
				pass.Reportf(valueSpec.Pos(),
					"enum constant %q uses string value; consider using uint64 with iota pattern instead",
					valueSpec.Names[0].Name)
			}
		}
	}
}

// isStringType checks if a type is string.
func isStringType(t types.Type) bool {
	basic, ok := t.(*types.Basic)
	if !ok {
		return false
	}

	return basic.Kind() == types.String
}

// hasStringLiteralValue checks if a value spec has a string literal value.
func hasStringLiteralValue(vs *ast.ValueSpec) bool {
	if len(vs.Values) == 0 {
		return false
	}

	for _, val := range vs.Values {
		lit, ok := val.(*ast.BasicLit)
		if ok && lit.Kind == token.STRING {
			return true
		}
	}

	return false
}
