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

// Package rawstring provides an analyzer that suggests raw strings over escaped strings.
package rawstring

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const (
	analyzerName = "attgo_raw_string"
	doc          = `suggests raw strings over escaped double-quoted strings

Prefer raw strings (backticks) over double-quoted strings with escape sequences.
Raw strings are more readable when the string contains quotes or backslashes.

Bad:
    query := "vouch_relay_execution_config_total{result=\"succeeded\"}"

Good:
    query := ` + "`" + `vouch_relay_execution_config_total{result="succeeded"}` + "`" + `

Exceptions:
- Strings containing backticks (cannot use raw string)
- Strings with actual newlines intended as \n
- Short strings with minimal escaping`
)

// Analyzer is the raw string preference analyzer.
var Analyzer = &analysis.Analyzer{
	Name: analyzerName,
	Doc:  doc,
	Run:  run,
}

// minEscapesForWarning is the minimum number of escape sequences to trigger a warning.
const minEscapesForWarning = 3

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			lit, ok := n.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return true
			}

			checkStringLiteral(pass, lit)

			return true
		})
	}

	return nil, nil
}

func checkStringLiteral(pass *analysis.Pass, lit *ast.BasicLit) {
	// Only check double-quoted strings.
	if !strings.HasPrefix(lit.Value, `"`) {
		return // Already a raw string.
	}

	value := lit.Value

	// Check if it contains backticks - can't convert to raw string.
	// Need to check the interpreted value.
	interpreted := interpretString(value)
	if strings.Contains(interpreted, "`") {
		return
	}

	// Count escape sequences.
	escapeCount := countEscapes(value)

	if escapeCount >= minEscapesForWarning {
		pass.Reportf(lit.Pos(),
			"string has %d escape sequences; consider using a raw string (backticks) for better readability",
			escapeCount)
	}
}

// countEscapes counts the number of escape sequences in a double-quoted string literal.
func countEscapes(s string) int {
	if len(s) < 2 {
		return 0
	}

	// Remove the surrounding quotes.
	s = s[1 : len(s)-1]

	count := 0
	i := 0

	for i < len(s) {
		if s[i] == '\\' && i+1 < len(s) {
			nextChar := s[i+1]
			// Count meaningful escapes that raw strings would eliminate.
			switch nextChar {
			case '"', '\\':
				count++
			case 'n', 't', 'r':
				// Don't count \n, \t, \r as these serve a purpose
				// that raw strings can't provide (except for literal newlines).
			default:
				// Other escapes like \x, \u, etc.
				count++
			}
			// Skip the escape sequence.
			i += 2

			continue
		}

		i++
	}

	return count
}

// interpretString interprets a Go string literal.
func interpretString(s string) string {
	if len(s) < 2 {
		return ""
	}

	// Remove surrounding quotes.
	s = s[1 : len(s)-1]

	var result strings.Builder

	i := 0

	for i < len(s) {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case 'n':
				result.WriteByte('\n')
			case 't':
				result.WriteByte('\t')
			case 'r':
				result.WriteByte('\r')
			case '\\':
				result.WriteByte('\\')
			case '"':
				result.WriteByte('"')
			case '\'':
				result.WriteByte('\'')
			case '`':
				result.WriteByte('`')
			default:
				// For simplicity, just write the escaped char.
				result.WriteByte(s[i+1])
			}

			i += 2

			continue
		}

		result.WriteByte(s[i])
		i++
	}

	return result.String()
}
