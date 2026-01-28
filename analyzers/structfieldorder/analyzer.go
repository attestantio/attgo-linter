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

// Package structfieldorder provides an analyzer that checks struct field ordering.
// Fields should be ordered: logger, metrics, dependencies, data, synchronization.
package structfieldorder

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const (
	analyzerName = "attgo_struct_field_order"
	doc          = `checks struct field ordering convention

Struct fields should be ordered by category:
1. Logger fields (log, logger)
2. Metrics fields (metrics, monitor)
3. Dependency fields (services, clients, external deps)
4. Data fields (configuration, state)
5. Synchronization fields (mutex, wg, channels)

This creates a predictable structure that makes code easier to navigate.

Example:
    type Service struct {
        // Logger
        log zerolog.Logger

        // Metrics
        metrics *prometheus.Registry

        // Dependencies
        client    *http.Client
        db        Database

        // Data
        config    Config
        cache     map[string]Value

        // Synchronization
        mu        sync.Mutex
        done      chan struct{}
    }`
)

// Analyzer is the struct field order analyzer.
var Analyzer = &analysis.Analyzer{
	Name: analyzerName,
	Doc:  doc,
	Run:  run,
}

// fieldCategory represents the category of a struct field.
type fieldCategory int

const (
	categoryUnknown fieldCategory = iota
	categoryLogger
	categoryMetrics
	categoryDependency
	categoryData
	categorySync
)

func (c fieldCategory) String() string {
	switch c {
	case categoryLogger:
		return "logger"
	case categoryMetrics:
		return "metrics"
	case categoryDependency:
		return "dependency"
	case categoryData:
		return "data"
	case categorySync:
		return "synchronization"
	default:
		return "unknown"
	}
}

func run(pass *analysis.Pass) (any, error) {
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

				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				checkStructFieldOrder(pass, typeSpec.Name.Name, structType)
			}
		}
	}

	return nil, nil
}

func checkStructFieldOrder(pass *analysis.Pass, structName string, st *ast.StructType) {
	if st.Fields == nil || len(st.Fields.List) == 0 {
		return
	}

	var lastCategory fieldCategory

	var lastCategoryField string

	for _, field := range st.Fields.List {
		if len(field.Names) == 0 {
			continue // Embedded field.
		}

		for _, name := range field.Names {
			cat := categorizeField(name.Name, field.Type)
			if cat == categoryUnknown {
				continue
			}

			if cat < lastCategory {
				pass.Reportf(name.Pos(),
					"field %q (%s) should come before %q (%s) in struct %q",
					name.Name, cat, lastCategoryField, lastCategory, structName)
			}

			lastCategory = cat
			lastCategoryField = name.Name
		}
	}
}

// categorizeField determines the category of a field based on name and type.
func categorizeField(name string, typ ast.Expr) fieldCategory {
	lowerName := strings.ToLower(name)

	// Logger fields.
	if lowerName == "log" || lowerName == "logger" || strings.HasSuffix(lowerName, "log") || strings.HasSuffix(lowerName, "logger") {
		return categoryLogger
	}

	// Metrics fields.
	if lowerName == "metrics" || lowerName == "monitor" || strings.HasSuffix(lowerName, "metrics") {
		return categoryMetrics
	}

	// Sync fields - check type.
	if isSyncType(typ) {
		return categorySync
	}

	// Check name patterns for sync.
	if lowerName == "mu" || lowerName == "mutex" || lowerName == "lock" ||
		lowerName == "wg" || strings.HasSuffix(lowerName, "mu") ||
		strings.HasSuffix(lowerName, "lock") || strings.HasSuffix(lowerName, "mutex") {
		return categorySync
	}

	// Channel fields.
	if _, ok := typ.(*ast.ChanType); ok {
		return categorySync
	}

	// Dependency fields - common patterns.
	if strings.HasSuffix(lowerName, "client") || strings.HasSuffix(lowerName, "service") ||
		strings.HasSuffix(lowerName, "provider") || strings.HasSuffix(lowerName, "handler") ||
		lowerName == "db" || lowerName == "database" || lowerName == "store" ||
		lowerName == "cache" || lowerName == "repo" || lowerName == "repository" {
		return categoryDependency
	}

	// Default to data for other fields.
	return categoryData
}

// isSyncType checks if a type is from the sync package.
func isSyncType(typ ast.Expr) bool {
	sel, ok := typ.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	if ident.Name == "sync" {
		switch sel.Sel.Name {
		case "Mutex", "RWMutex", "WaitGroup", "Once", "Cond", "Pool", "Map":
			return true
		}
	}

	return false
}
