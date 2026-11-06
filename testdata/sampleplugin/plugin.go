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

// Package sampleplugin registers two analyzers under their own names. The
// integration test compiles it into a custom golangci-lint binary next to the
// attgo plugin. Both analyzers report the same literals, so a //nolint
// directive that names one of them shows whether the other still reports.
package sampleplugin

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/attestantio/attgo-linter/internal/pername"
	"golang.org/x/tools/go/analysis"
)

func init() {
	pername.Register(newAnalyzer("sample_first"))
	pername.Register(newAnalyzer("sample_second"))
}

func newAnalyzer(name string) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: name,
		Doc:  "reports string literals that contain flag-me",
		Run: func(pass *analysis.Pass) (any, error) {
			for _, file := range pass.Files {
				ast.Inspect(file, func(n ast.Node) bool {
					lit, ok := n.(*ast.BasicLit)
					if ok && lit.Kind == token.STRING && strings.Contains(lit.Value, "flag-me") {
						pass.Reportf(lit.Pos(), "string literal contains flag-me")
					}

					return true
				})
			}

			return nil, nil
		},
	}
}
