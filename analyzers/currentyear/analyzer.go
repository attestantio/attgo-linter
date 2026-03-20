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
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
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

// copyrightYearPattern matches common copyright year formats.
// Matches patterns like:
// - Copyright © 2024
// - Copyright 2024
// - Copyright (c) 2024
// - Copyright © 2023-2024 (captures both first and last year in range)
// Group 1: first year (optional, only present in ranges)
// Group 2: last year (always present)
var copyrightYearPattern = regexp.MustCompile(`[Cc]opyright\s*(?:©|\(c\))?\s*(?:(\d{4})\s*-\s*)?(\d{4})`)

// runner holds the configuration and cached git state for the current year analyzer.
type runner struct {
	currentYear  int
	repoRoot     string
	changedFiles map[string]fileStatus
	gitAvailable bool
}

// NewAnalyzer creates a new current year copyright analyzer.
// If baseRef is non-empty, only files changed relative to that git ref are checked.
// Git results are resolved once here and cached, so that run() (called once per
// package by the analysis framework) never spawns subprocesses.
func NewAnalyzer(baseRef string) *analysis.Analyzer {
	r := &runner{
		currentYear: time.Now().Year(),
	}

	if baseRef != "" {
		repoRoot, changedFiles, err := resolveChangedFiles(baseRef)
		if err != nil {
			fmt.Fprintf(os.Stderr, "attgo_current_year: warning: %v; checking all files\n", err)
		} else {
			r.repoRoot = filepath.Clean(repoRoot)
			r.changedFiles = changedFiles
			r.gitAvailable = true
		}
	}

	return &analysis.Analyzer{
		Name: analyzerName,
		Doc:  doc,
		Run:  r.run,
	}
}

func (r *runner) run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if r.gitAvailable {
			filePath := pass.Fset.Position(file.Package).Filename
			status := resolveFileStatus(filePath, r.repoRoot, r.changedFiles)

			if status == fileStatusUnchanged {
				// File is unchanged; skip it.
				continue
			}

			checkFile(pass, file, r.currentYear, status)
		} else {
			// No baseRef or git not available: check all files.
			checkFile(pass, file, r.currentYear, fileStatusUnchanged)
		}
	}

	return nil, nil
}

func checkFile(pass *analysis.Pass, file *ast.File, currentYear int, status fileStatus) {
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

	if len(matches) < 3 {
		// No copyright year found in header - that's ok, goheader linter handles format.
		return
	}

	lastYearStr := matches[2]

	lastYear, err := strconv.Atoi(lastYearStr)
	if err != nil {
		return
	}

	// Check if the year is current.
	if lastYear >= currentYear {
		return
	}

	// Determine the first year (for modified file range suggestions).
	firstYear := lastYear
	if matches[1] != "" {
		if fy, err := strconv.Atoi(matches[1]); err == nil {
			firstYear = fy
		}
	}

	switch status {
	case fileStatusModified:
		pass.Reportf(copyrightComment.Pos(),
			"copyright year %d is outdated; should be %d-%d for modified files",
			lastYear, firstYear, currentYear)
	case fileStatusNew:
		pass.Reportf(copyrightComment.Pos(),
			"copyright year %d is outdated; should be %d for new files",
			lastYear, currentYear)
	case fileStatusUnchanged:
		// No baseRef or git not available: backward-compat message.
		pass.Reportf(copyrightComment.Pos(),
			"copyright year %d is outdated; should be %d for new or modified files",
			lastYear, currentYear)
	}
}
