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
	"os/exec"
	"path/filepath"
	"strings"
)

// fileStatus represents the git status of a file.
type fileStatus string

const (
	// fileStatusNew indicates a file was added.
	fileStatusNew fileStatus = "new"
	// fileStatusModified indicates a file was modified.
	fileStatusModified fileStatus = "modified"
)

// resolveChangedFiles performs a single LookPath check, resolves the repo root,
// and returns the changed files map. This avoids redundant LookPath calls.
func resolveChangedFiles(baseRef string) (string, map[string]fileStatus, error) {
	if strings.HasPrefix(baseRef, "-") {
		return "", nil, fmt.Errorf("invalid base ref %q: must not start with '-'", baseRef)
	}

	if _, err := exec.LookPath("git"); err != nil {
		return "", nil, fmt.Errorf("git not found: %w", err)
	}

	repoRoot, err := gitRepoRoot()
	if err != nil {
		return "", nil, fmt.Errorf("failed to get git repo root: %w", err)
	}

	changedFiles, err := gitChangedFiles(baseRef)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get git changed files: %w", err)
	}

	return repoRoot, changedFiles, nil
}

// gitChangedFiles returns a map of changed files relative to the given base ref.
// The keys are git-root-relative paths and the values indicate whether the file
// is new or modified.
// Callers must ensure git is available (via exec.LookPath) before calling.
func gitChangedFiles(baseRef string) (map[string]fileStatus, error) {
	//nolint:gosec // baseRef is validated by resolveChangedFiles before reaching here.
	cmd := exec.Command("git", "diff", "--name-status", baseRef)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git diff failed: %w", err)
	}

	result := make(map[string]fileStatus)

	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "\t")
		if len(parts) < 2 {
			continue
		}

		status := parts[0]
		switch status {
		case "A":
			result[parts[1]] = fileStatusNew
		case "M":
			result[parts[1]] = fileStatusModified
		default:
			if strings.HasPrefix(status, "R") {
				// Renamed files have format: R<score>\told\tnew.
				if len(parts) >= 3 {
					result[parts[2]] = fileStatusModified
				}
			}
			// Ignore other statuses (D, C, etc.).
		}
	}

	return result, nil
}

// gitRepoRoot returns the absolute path to the root of the git repository.
// Callers must ensure git is available (via exec.LookPath) before calling.
func gitRepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git rev-parse failed: %w", err)
	}

	return strings.TrimSpace(string(out)), nil
}

// resolveFileStatus converts an absolute file path to a repo-relative path
// and looks it up in the changed files map.
// Returns the file's status, or empty string if the file is not in the map.
func resolveFileStatus(filePath string, repoRoot string, changedFiles map[string]fileStatus) fileStatus {
	if changedFiles == nil {
		return ""
	}

	rel, err := filepath.Rel(repoRoot, filePath)
	if err != nil {
		return ""
	}

	return changedFiles[rel]
}
