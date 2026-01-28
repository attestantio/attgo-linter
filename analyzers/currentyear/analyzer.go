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

// Package currentyear provides an analyzer that checks for current year in copyright headers.
package currentyear

import (
	"go/ast"
	"regexp"
	"strconv"
	"time"

	"golang.org/x/tools/go/analysis"
)

const (
	analyzerName = "attgo_current_year"
	doc          = `checks for current year in copyright headers

New files should have the current year in their copyright header.
This helps maintain accurate copyright information.

Bad (in 2025):
    // Copyright © 2024 Attestant Limited.

Good:
    // Copyright © 2026 Attestant Limited.

Also acceptable (year ranges ending in current year):
    // Copyright © 2023-2025 Attestant Limited.`
)

// Analyzer is the current year copyright analyzer.
var Analyzer = &analysis.Analyzer{
	Name: analyzerName,
	Doc:  doc,
	Run:  run,
}

// copyrightYearPattern matches common copyright year formats.
// Matches patterns like:
// - Copyright © 2024
// - Copyright 2024
// - Copyright (c) 2024
// - Copyright © 2023-2024 (captures last year in range)
var copyrightYearPattern = regexp.MustCompile(`[Cc]opyright\s*(?:©|\(c\))?\s*(?:\d{4}\s*-\s*)?(\d{4})`)

func run(pass *analysis.Pass) (any, error) {
	currentYear := time.Now().Year()

	for _, file := range pass.Files {
		checkFile(pass, file, currentYear)
	}

	return nil, nil
}

func checkFile(pass *analysis.Pass, file *ast.File, currentYear int) {
	// Get the first comment group (copyright header).
	if len(file.Comments) == 0 {
		return
	}

	// Look at the first comment group that appears before the package declaration.
	var copyrightComment *ast.CommentGroup

	for _, cg := range file.Comments {
		if cg.Pos() < file.Package {
			copyrightComment = cg

			break
		}
	}

	if copyrightComment == nil {
		return
	}

	// Extract year from copyright comment.
	text := copyrightComment.Text()
	matches := copyrightYearPattern.FindStringSubmatch(text)

	if len(matches) < 2 {
		// No copyright year found in header - that's ok, goheader linter handles format.
		return
	}

	yearStr := matches[1]

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return
	}

	// Check if the year is current.
	if year < currentYear {
		pass.Reportf(copyrightComment.Pos(),
			"copyright year %d is outdated; should be %d for new or modified files",
			year, currentYear)
	}
}
