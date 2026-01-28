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

// Package interfacecheck provides an analyzer that suggests compile-time interface compliance checks.
package interfacecheck

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const (
	analyzerName = "attgo_interface_check"
	doc          = `suggests compile-time interface compliance checks

When a struct type implements an interface, add a compile-time check
to ensure the implementation remains complete.

Pattern:
    var _ Interface = (*Struct)(nil)

This catches missing method implementations at compile time rather than runtime.

Example:
    type Reader interface {
        Read(p []byte) (n int, err error)
    }

    type MyReader struct{}

    // Compile-time check
    var _ Reader = (*MyReader)(nil)

    func (r *MyReader) Read(p []byte) (int, error) {
        return 0, nil
    }`
)

// Analyzer is the interface check analyzer.
var Analyzer = &analysis.Analyzer{
	Name: analyzerName,
	Doc:  doc,
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	// Collect all interfaces and structs defined in this package.
	interfaces := make(map[string]*types.Interface)
	structs := make(map[string]*types.Struct)

	for _, name := range pass.Pkg.Scope().Names() {
		obj := pass.Pkg.Scope().Lookup(name)
		if obj == nil {
			continue
		}

		typeName, ok := obj.(*types.TypeName)
		if !ok {
			continue
		}

		switch t := typeName.Type().Underlying().(type) {
		case *types.Interface:
			if t.NumMethods() > 0 { // Skip empty interfaces.
				interfaces[name] = t
			}
		case *types.Struct:
			structs[name] = t
		}
	}

	// Collect existing interface checks (var _ Interface = (*Struct)(nil)).
	existingChecks := collectExistingChecks(pass)

	// For each struct, check which interfaces it implements.
	for structName := range structs {
		structObj := pass.Pkg.Scope().Lookup(structName)
		if structObj == nil {
			continue
		}

		structType := structObj.Type()
		ptrType := types.NewPointer(structType)

		for ifaceName, iface := range interfaces {
			// Check if the struct (or pointer to struct) implements the interface.
			if !types.Implements(structType, iface) && !types.Implements(ptrType, iface) {
				continue
			}

			// Check if there's already a compliance check.
			key := ifaceName + ":" + structName

			if existingChecks[key] {
				continue
			}

			// Find the struct definition to report the diagnostic.
			pos := findStructPos(pass, structName)
			if pos == token.NoPos {
				continue
			}

			pass.Reportf(pos,
				"struct %q implements interface %q; consider adding: var _ %s = (*%s)(nil)",
				structName, ifaceName, ifaceName, structName)
		}
	}

	return nil, nil
}

// collectExistingChecks finds all var _ Interface = (*Struct)(nil) patterns.
func collectExistingChecks(pass *analysis.Pass) map[string]bool {
	checks := make(map[string]bool)

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.VAR {
				continue
			}

			for _, spec := range genDecl.Specs {
				valueSpec, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}

				// Check for blank identifier.
				if len(valueSpec.Names) != 1 || valueSpec.Names[0].Name != "_" {
					continue
				}

				// Get the interface name from the type.
				ifaceName := getInterfaceName(valueSpec.Type)
				if ifaceName == "" {
					continue
				}

				// Get the struct name from the value.
				structName := getStructNameFromNilCast(valueSpec)
				if structName == "" {
					continue
				}

				key := ifaceName + ":" + structName
				checks[key] = true
			}
		}
	}

	return checks
}

// getInterfaceName extracts the interface name from a type expression.
func getInterfaceName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return t.Sel.Name
	}

	return ""
}

// getStructNameFromNilCast extracts the struct name from (*Struct)(nil) pattern.
func getStructNameFromNilCast(vs *ast.ValueSpec) string {
	if len(vs.Values) != 1 {
		return ""
	}

	// Expect (*Type)(nil).
	call, ok := vs.Values[0].(*ast.CallExpr)
	if !ok {
		return ""
	}

	// Check for nil argument.
	if len(call.Args) != 1 {
		return ""
	}

	ident, ok := call.Args[0].(*ast.Ident)
	if !ok || ident.Name != "nil" {
		return ""
	}

	// Get the type from (*Type).
	paren, ok := call.Fun.(*ast.ParenExpr)
	if !ok {
		return ""
	}

	star, ok := paren.X.(*ast.StarExpr)
	if !ok {
		return ""
	}

	switch x := star.X.(type) {
	case *ast.Ident:
		return x.Name
	case *ast.SelectorExpr:
		return x.Sel.Name
	}

	return ""
}

// findStructPos finds the position of a struct type definition.
func findStructPos(pass *analysis.Pass, name string) token.Pos {
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

				if typeSpec.Name.Name == name {
					return typeSpec.Name.Pos()
				}
			}
		}
	}

	return token.NoPos
}

// Ensure strings is used.
var _ = strings.Contains
