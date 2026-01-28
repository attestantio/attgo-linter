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

// Package capitalcomment provides an analyzer that checks comments start with a capital letter.
package capitalcomment

import (
	"go/ast"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

const (
	analyzerName = "attgo_capital_comment"
	doc          = `checks that comments start with a capital letter

Comments should start with a capital letter for consistency and readability.
Exceptions are made for:
- Comments starting with code references (identifiers)
- nolint directives
- URLs
- Comments that start with punctuation

Bad:
    // this is a comment

Good:
    // This is a comment
    // someVariable is used for...`
)

// Analyzer is the capital comment analyzer.
var Analyzer = &analysis.Analyzer{
	Name: analyzerName,
	Doc:  doc,
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		for _, cg := range file.Comments {
			// Only check the first comment in each group.
			// Subsequent comments are continuations and may legitimately start lowercase.
			if len(cg.List) > 0 {
				checkComment(pass, cg.List[0])
			}
		}
	}

	return nil, nil
}

func checkComment(pass *analysis.Pass, c *ast.Comment) {
	text := c.Text

	// Remove comment prefix.
	text = strings.TrimPrefix(text, "//")
	text = strings.TrimPrefix(text, "/*")
	text = strings.TrimSuffix(text, "*/")
	text = strings.TrimSpace(text)

	if len(text) == 0 {
		return
	}

	// Get the first rune.
	firstRune := []rune(text)[0]

	// Skip if starts with punctuation or number.
	if unicode.IsPunct(firstRune) || unicode.IsDigit(firstRune) {
		return
	}

	// Skip special patterns.
	if shouldSkip(text) {
		return
	}

	// Check if it starts with lowercase.
	if unicode.IsLower(firstRune) {
		// Check if this might be an identifier reference.
		if looksLikeIdentifier(text) {
			return
		}

		pass.Reportf(c.Pos(), "comment should start with a capital letter")
	}
}

// shouldSkip returns true if the comment should be skipped from checking.
func shouldSkip(text string) bool {
	lowerText := strings.ToLower(text)

	// Skip nolint directives.
	if strings.HasPrefix(lowerText, "nolint") {
		return true
	}

	// Skip TODO, FIXME, etc. (case insensitive prefix).
	prefixes := []string{"todo", "fixme", "hack", "xxx", "bug"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(lowerText, prefix) {
			return true
		}
	}

	// Skip URLs.
	if strings.Contains(text, "://") {
		return true
	}

	// Skip build tags.
	if strings.HasPrefix(text, "+build") || strings.HasPrefix(text, "go:") {
		return true
	}

	// Skip standard license header text (Apache 2.0, MIT, BSD, etc.).
	licensePatterns := []string{
		"you may not use this file",
		"distributed under the license",
		"without warranties or conditions",
		"limitations under the license",
		"permission is hereby granted",
		"the above copyright notice",
		"in no event shall",
		"as is",
	}
	for _, pattern := range licensePatterns {
		if strings.Contains(lowerText, pattern) {
			return true
		}
	}

	return false
}

// commonEnglishWords are words that shouldn't be treated as identifiers.
var commonEnglishWords = map[string]bool{
	"this": true, "that": true, "these": true, "those": true,
	"it": true, "its": true, "the": true, "a": true, "an": true,
	"here": true, "there": true, "where": true, "when": true,
	"see": true, "use": true, "set": true, "get": true,
	"all": true, "any": true, "some": true, "each": true,
	"for": true, "not": true, "but": true, "and": true, "or": true,
}

// looksLikeIdentifier checks if the comment appears to start with an identifier.
// This handles patterns like "someFunc is...", "myVar contains...".
func looksLikeIdentifier(text string) bool {
	// Find the first word.
	words := strings.Fields(text)
	if len(words) == 0 {
		return false
	}

	firstWord := words[0]

	// Common English words are not identifiers.
	if commonEnglishWords[strings.ToLower(firstWord)] {
		return false
	}

	// Check if it looks like a code identifier:
	// - Contains underscore (snake_case)
	// - Contains mixed case after first char (camelCase)
	// - Is all lowercase and followed by "is", "are", "was", "contains", etc.

	if strings.Contains(firstWord, "_") {
		return true
	}

	// Check for camelCase (lowercase start, has uppercase within).
	hasInternalUpper := false

	for i, r := range firstWord {
		if i > 0 && unicode.IsUpper(r) {
			hasInternalUpper = true

			break
		}
	}

	if hasInternalUpper {
		return true
	}

	// Check if followed by common identifier-describing words.
	if len(words) >= 2 {
		followWord := strings.ToLower(words[1])
		identifierFollowers := []string{
			"is", "are", "was", "were", "has", "have", "had",
			"contains", "returns", "holds", "stores", "represents",
			"defines", "implements", "provides", "specifies",
		}

		for _, follower := range identifierFollowers {
			if followWord == follower {
				return true
			}
		}
	}

	return false
}
