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

package currentyear

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis"
)

func TestCheckFileMessages(t *testing.T) {
	tests := []struct {
		name        string
		status      fileStatus
		year        int
		noComment   bool
		currentYear int
		wantMsg     string
	}{
		{
			name:        "new file message",
			status:      fileStatusNew,
			year:        2020,
			currentYear: 2026,
			wantMsg:     "copyright year 2020 is outdated; should be 2026 for new files",
		},
		{
			name:        "modified file message",
			status:      fileStatusModified,
			year:        2020,
			currentYear: 2026,
			wantMsg:     "copyright year 2020 is outdated; should be 2020-2026 for modified files",
		},
		{
			name:        "empty status backward compat message",
			status:      "",
			year:        2020,
			currentYear: 2026,
			wantMsg:     "copyright year 2020 is outdated; should be 2026 for new or modified files",
		},
		{
			name:        "modified file with range",
			status:      fileStatusModified,
			year:        2023,
			currentYear: 2026,
			wantMsg:     "copyright year 2023 is outdated; should be 2023-2026 for modified files",
		},
		{
			name:        "current year does not report",
			status:      fileStatusNew,
			year:        2026,
			currentYear: 2026,
			wantMsg:     "",
		},
		{
			name:        "no comments does not report",
			status:      fileStatusNew,
			noComment:   true,
			currentYear: 2026,
			wantMsg:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			f := fset.AddFile("test.go", -1, 1000)

			// Build a minimal ast.File with a copyright comment.
			var comments []*ast.CommentGroup
			if !tt.noComment {
				comments = []*ast.CommentGroup{
					{
						List: []*ast.Comment{
							{
								Slash: token.Pos(f.Base()),
								Text:  fmt.Sprintf("// Copyright © %d Attestant Limited.", tt.year),
							},
						},
					},
				}
			}

			astFile := &ast.File{
				Package:  token.Pos(f.Base() + 100),
				Comments: comments,
			}

			var diagnostics []analysis.Diagnostic

			pass := &analysis.Pass{
				Fset: fset,
				Report: func(d analysis.Diagnostic) {
					diagnostics = append(diagnostics, d)
				},
			}

			checkFile(pass, astFile, tt.currentYear, tt.status)

			if tt.wantMsg == "" {
				if len(diagnostics) != 0 {
					t.Errorf("expected no diagnostics, got %d: %v", len(diagnostics), diagnostics)
				}

				return
			}

			if len(diagnostics) != 1 {
				t.Fatalf("expected 1 diagnostic, got %d", len(diagnostics))
			}

			if !strings.Contains(diagnostics[0].Message, tt.wantMsg) {
				t.Errorf("diagnostic message = %q, want it to contain %q", diagnostics[0].Message, tt.wantMsg)
			}
		})
	}
}

func TestCheckFileModifiedWithRange(t *testing.T) {
	fset := token.NewFileSet()
	f := fset.AddFile("test.go", -1, 1000)

	// Copyright with a year range: 2018-2023.
	comments := []*ast.CommentGroup{
		{
			List: []*ast.Comment{
				{
					Slash: token.Pos(f.Base()),
					Text:  "// Copyright © 2018-2023 Attestant Limited.",
				},
			},
		},
	}

	astFile := &ast.File{
		Package:  token.Pos(f.Base() + 100),
		Comments: comments,
	}

	var diagnostics []analysis.Diagnostic

	pass := &analysis.Pass{
		Fset: fset,
		Report: func(d analysis.Diagnostic) {
			diagnostics = append(diagnostics, d)
		},
	}

	checkFile(pass, astFile, 2026, fileStatusModified)

	if len(diagnostics) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diagnostics))
	}

	want := "copyright year 2023 is outdated; should be 2018-2026 for modified files"
	if diagnostics[0].Message != want {
		t.Errorf("diagnostic message = %q, want %q", diagnostics[0].Message, want)
	}
}
