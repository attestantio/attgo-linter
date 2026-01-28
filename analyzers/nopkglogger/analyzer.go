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

// Package nopkglogger provides an analyzer that detects package-level logger variables.
// Loggers should be struct fields, not package-level variables, to enable proper
// dependency injection and testability.
package nopkglogger

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const (
	analyzerName = "attgo_no_pkg_logger"
	doc          = `detects package-level logger variables

Loggers should be struct fields, not package-level variables. This enables:
- Proper dependency injection
- Better testability (inject mock loggers)
- Clear ownership of logging configuration

Bad:
    var log zerolog.Logger

Good:
    type Service struct {
        log zerolog.Logger
    }`
)

// NewAnalyzer creates a new no-pkg-logger analyzer with the given logger type patterns.
func NewAnalyzer(loggerTypePatterns []string) *analysis.Analyzer {
	r := &runner{
		loggerTypePatterns: loggerTypePatterns,
	}

	return &analysis.Analyzer{
		Name: analyzerName,
		Doc:  doc,
		Run:  r.run,
	}
}

type runner struct {
	loggerTypePatterns []string
}

func (r *runner) run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}

			for _, spec := range genDecl.Specs {
				valueSpec, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}

				// Check each variable in the declaration.
				for _, name := range valueSpec.Names {
					obj := pass.TypesInfo.ObjectOf(name)
					if obj == nil {
						continue
					}

					// Only check package-level variables.
					if obj.Parent() != obj.Pkg().Scope() {
						continue
					}

					// Check if the type matches any logger pattern.
					if r.isLoggerType(obj.Type()) {
						pass.Reportf(name.Pos(),
							"package-level logger %q detected; loggers should be struct fields for better dependency injection and testability",
							name.Name)
					}
				}
			}
		}
	}

	return nil, nil
}

// isLoggerType checks if the given type matches any of the configured logger patterns.
func (r *runner) isLoggerType(t types.Type) bool {
	typeName := typeString(t)

	for _, pattern := range r.loggerTypePatterns {
		if matchTypePattern(typeName, pattern) {
			return true
		}
	}

	return false
}

// typeString returns a string representation of the type suitable for pattern matching.
func typeString(t types.Type) string {
	// Get the full type string including package path.
	return types.TypeString(t, nil)
}

// matchTypePattern checks if a type string matches a pattern.
// Patterns can be:
// - Exact match: "zerolog.Logger"
// - Pointer: "*zerolog.Logger"
// The type string from types.TypeString includes the full package path,
// so we match against the suffix.
func matchTypePattern(typeName, pattern string) bool {
	// Handle pointer patterns.
	if strings.HasPrefix(pattern, "*") {
		if !strings.HasPrefix(typeName, "*") {
			return false
		}
		typeName = strings.TrimPrefix(typeName, "*")
		pattern = strings.TrimPrefix(pattern, "*")
	} else if strings.HasPrefix(typeName, "*") {
		// Pattern is not pointer but type is.
		return false
	}

	// Check if the type name ends with the pattern (handles full package paths).
	// e.g., "github.com/rs/zerolog.Logger" ends with "zerolog.Logger"
	if strings.HasSuffix(typeName, pattern) {
		// Ensure we match at a package boundary.
		prefix := strings.TrimSuffix(typeName, pattern)
		if prefix == "" || strings.HasSuffix(prefix, "/") || strings.HasSuffix(prefix, ".") {
			return true
		}
	}

	return false
}
