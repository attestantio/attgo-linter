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
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/tools/go/analysis"
)

func TestIntegration_VouchRepo(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	currentYear := time.Now().Year()

	// Setup: clone vouch into a temp directory.
	tmpDir := t.TempDir()

	// Resolve symlinks so paths match git's repo root (macOS /var -> /private/var).
	tmpDir, err := filepath.EvalSymlinks(tmpDir)
	if err != nil {
		t.Fatalf("failed to resolve tmpdir symlinks: %v", err)
	}

	vouchDir := filepath.Join(tmpDir, "vouch")

	runGit(t, tmpDir, "clone", "--depth=1", "https://github.com/attestantio/vouch.git", vouchDir)
	runGit(t, vouchDir, "checkout", "-b", "test-copyright-check")

	// --- File 1: Modified existing file, copyright NOT updated (should be FLAGGED). ---
	file1Rel := "util/scatter.go"
	file1Abs := filepath.Join(vouchDir, file1Rel)

	file1Content, err := os.ReadFile(file1Abs)
	if err != nil {
		t.Fatalf("failed to read %s: %v", file1Rel, err)
	}

	err = os.WriteFile(file1Abs, append(file1Content, []byte("\n// test modification\n")...), 0o644)
	if err != nil {
		t.Fatalf("failed to write %s: %v", file1Rel, err)
	}

	runGit(t, vouchDir, "add", file1Rel)

	// --- File 2: Modified existing file, copyright CORRECTLY updated to range (should PASS). ---
	file2Rel := "util/epoch.go"
	file2Abs := filepath.Join(vouchDir, file2Rel)

	file2Content, err := os.ReadFile(file2Abs)
	if err != nil {
		t.Fatalf("failed to read %s: %v", file2Rel, err)
	}

	// Update the copyright year to a range ending in current year.
	file2Str := string(file2Content)
	matches := copyrightYearPattern.FindStringSubmatch(file2Str)
	if len(matches) < 3 {
		t.Fatalf("could not find copyright year in %s", file2Rel)
	}

	// Extract the original year (either first year or last year if no range).
	origYear := matches[2]
	if matches[1] != "" {
		origYear = matches[1]
	}

	updatedCopyright := strings.Replace(
		file2Str,
		matches[0],
		fmt.Sprintf("Copyright © %s-%d", origYear, currentYear),
		1,
	)
	updatedCopyright += "\n// test modification\n"

	err = os.WriteFile(file2Abs, []byte(updatedCopyright), 0o644)
	if err != nil {
		t.Fatalf("failed to write %s: %v", file2Rel, err)
	}

	runGit(t, vouchDir, "add", file2Rel)

	// --- File 3: New file with OLD copyright year (should be FLAGGED). ---
	file3Rel := "util/testfile_new_old.go"
	file3Abs := filepath.Join(vouchDir, file3Rel)

	err = os.WriteFile(file3Abs, []byte(`// Copyright © 2020 Attestant Limited.
// Licensed under the Apache License, Version 2.0
package util

var TestNewOld = true
`), 0o644)
	if err != nil {
		t.Fatalf("failed to write %s: %v", file3Rel, err)
	}

	runGit(t, vouchDir, "add", file3Rel)

	// --- File 4: New file with CURRENT copyright year (should PASS). ---
	file4Rel := "util/testfile_new_current.go"
	file4Abs := filepath.Join(vouchDir, file4Rel)

	err = os.WriteFile(file4Abs, []byte(fmt.Sprintf(`// Copyright © %d Attestant Limited.
// Licensed under the Apache License, Version 2.0
package util

var TestNewCurrent = true
`, currentYear)), 0o644)
	if err != nil {
		t.Fatalf("failed to write %s: %v", file4Rel, err)
	}

	runGit(t, vouchDir, "add", file4Rel)

	// Commit all changes.
	runGit(t, vouchDir, "commit", "-m", "test changes")

	// --- Run the analyzer functions from within the vouch repo. ---

	// t.Chdir is goroutine-safe (Go 1.24+) and auto-restores on cleanup.
	t.Chdir(vouchDir)

	repoRoot, changedFiles, err := resolveChangedFiles("auto")
	if err != nil {
		t.Fatalf("resolveChangedFiles failed: %v", err)
	}

	// Verify changed files map contains expected files with correct statuses.
	assertFileStatus(t, changedFiles, file1Rel, fileStatusModified)
	assertFileStatus(t, changedFiles, file2Rel, fileStatusModified)
	assertFileStatus(t, changedFiles, file3Rel, fileStatusNew)
	assertFileStatus(t, changedFiles, file4Rel, fileStatusNew)

	// Verify resolveFileStatus works for each file.
	if got := resolveFileStatus(file1Abs, repoRoot, changedFiles); got != fileStatusModified {
		t.Errorf("resolveFileStatus(%s) = %q, want %q", file1Rel, got, fileStatusModified)
	}

	if got := resolveFileStatus(file2Abs, repoRoot, changedFiles); got != fileStatusModified {
		t.Errorf("resolveFileStatus(%s) = %q, want %q", file2Rel, got, fileStatusModified)
	}

	if got := resolveFileStatus(file3Abs, repoRoot, changedFiles); got != fileStatusNew {
		t.Errorf("resolveFileStatus(%s) = %q, want %q", file3Rel, got, fileStatusNew)
	}

	if got := resolveFileStatus(file4Abs, repoRoot, changedFiles); got != fileStatusNew {
		t.Errorf("resolveFileStatus(%s) = %q, want %q", file4Rel, got, fileStatusNew)
	}

	// Verify untouched file is NOT in the changed files map.
	untouchedFile := filepath.Join(vouchDir, "main.go")
	if got := resolveFileStatus(untouchedFile, repoRoot, changedFiles); got != fileStatusUnchanged {
		t.Errorf("resolveFileStatus(main.go) = %d, want fileStatusUnchanged", got)
	}

	// --- Run checkFile on each test file and verify diagnostics. ---

	// File 1: modified, copyright NOT updated -> should be flagged.
	diags := runCheckFile(t, file1Abs, currentYear, fileStatusModified)
	if len(diags) != 1 {
		t.Errorf("file1 (modified, old copyright): expected 1 diagnostic, got %d", len(diags))
	} else if !strings.Contains(diags[0].Message, "modified files") {
		t.Errorf("file1: expected 'modified files' in message, got %q", diags[0].Message)
	}

	// File 2: modified, copyright updated to current range -> should NOT be flagged.
	diags = runCheckFile(t, file2Abs, currentYear, fileStatusModified)
	if len(diags) != 0 {
		t.Errorf("file2 (modified, current copyright): expected 0 diagnostics, got %d: %v", len(diags), diags)
	}

	// File 3: new file, old copyright -> should be flagged.
	diags = runCheckFile(t, file3Abs, currentYear, fileStatusNew)
	if len(diags) != 1 {
		t.Errorf("file3 (new, old copyright): expected 1 diagnostic, got %d", len(diags))
	} else if !strings.Contains(diags[0].Message, "new files") {
		t.Errorf("file3: expected 'new files' in message, got %q", diags[0].Message)
	}

	// File 4: new file, current copyright -> should NOT be flagged.
	diags = runCheckFile(t, file4Abs, currentYear, fileStatusNew)
	if len(diags) != 0 {
		t.Errorf("file4 (new, current copyright): expected 0 diagnostics, got %d: %v", len(diags), diags)
	}
}

// runGit runs a git command in the given directory and fails the test on error.
func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()

	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=Test",
		"GIT_AUTHOR_EMAIL=test@test.com",
		"GIT_COMMITTER_NAME=Test",
		"GIT_COMMITTER_EMAIL=test@test.com",
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v\n%s", strings.Join(args, " "), err, out)
	}
}

// runCheckFile parses the given Go file and runs checkFile, returning any diagnostics.
func runCheckFile(t *testing.T, filePath string, currentYear int, status fileStatus) []analysis.Diagnostic {
	t.Helper()

	fset := token.NewFileSet()

	astFile, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		t.Fatalf("failed to parse %s: %v", filePath, err)
	}

	var diagnostics []analysis.Diagnostic

	pass := &analysis.Pass{
		Fset: fset,
		Report: func(d analysis.Diagnostic) {
			diagnostics = append(diagnostics, d)
		},
	}

	checkFile(pass, astFile, currentYear, status)

	return diagnostics
}

// assertFileStatus checks that a file appears in the changed files map with the expected status.
func assertFileStatus(t *testing.T, changedFiles map[string]fileStatus, relPath string, expected fileStatus) {
	t.Helper()

	got, ok := changedFiles[relPath]
	if !ok {
		t.Errorf("file %q not found in changed files map", relPath)

		return
	}

	if got != expected {
		t.Errorf("file %q status = %q, want %q", relPath, got, expected)
	}
}
