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
	"strings"
	"testing"
)

func TestResolveFileStatus(t *testing.T) {
	tests := []struct {
		name         string
		filePath     string
		repoRoot     string
		changedFiles map[string]fileStatus
		expected     fileStatus
	}{
		{
			name:     "new file found",
			filePath: "/home/user/project/pkg/foo.go",
			repoRoot: "/home/user/project",
			changedFiles: map[string]fileStatus{
				"pkg/foo.go": fileStatusNew,
			},
			expected: fileStatusNew,
		},
		{
			name:     "modified file found",
			filePath: "/home/user/project/pkg/bar.go",
			repoRoot: "/home/user/project",
			changedFiles: map[string]fileStatus{
				"pkg/bar.go": fileStatusModified,
			},
			expected: fileStatusModified,
		},
		{
			name:     "file not in changed files",
			filePath: "/home/user/project/pkg/baz.go",
			repoRoot: "/home/user/project",
			changedFiles: map[string]fileStatus{
				"pkg/foo.go": fileStatusNew,
			},
			expected: "",
		},
		{
			name:     "file at repo root",
			filePath: "/home/user/project/main.go",
			repoRoot: "/home/user/project",
			changedFiles: map[string]fileStatus{
				"main.go": fileStatusNew,
			},
			expected: fileStatusNew,
		},
		{
			name:         "empty changed files map",
			filePath:     "/home/user/project/pkg/foo.go",
			repoRoot:     "/home/user/project",
			changedFiles: map[string]fileStatus{},
			expected:     "",
		},
		{
			name:         "nil changed files map",
			filePath:     "/home/user/project/pkg/foo.go",
			repoRoot:     "/home/user/project",
			changedFiles: nil,
			expected:     "",
		},
		{
			name:     "repo root with trailing slash",
			filePath: "/home/user/project/pkg/foo.go",
			repoRoot: "/home/user/project/",
			changedFiles: map[string]fileStatus{
				"pkg/foo.go": fileStatusNew,
			},
			expected: fileStatusNew,
		},
		{
			name:     "deeply nested file",
			filePath: "/home/user/project/a/b/c/d/file.go",
			repoRoot: "/home/user/project",
			changedFiles: map[string]fileStatus{
				"a/b/c/d/file.go": fileStatusModified,
			},
			expected: fileStatusModified,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveFileStatus(tt.filePath, tt.repoRoot, tt.changedFiles)
			if got != tt.expected {
				t.Errorf("resolveFileStatus() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestGitDefaultBranch(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test requiring git repo in short mode")
	}

	branch, err := gitDefaultBranch()
	if err != nil {
		t.Skipf("git symbolic-ref failed (expected in CI without origin): %v", err)
	}

	// Should be origin/main or origin/master.
	if branch != "origin/main" && branch != "origin/master" {
		t.Errorf("gitDefaultBranch() = %q, want origin/main or origin/master", branch)
	}
}

func TestResolveChangedFiles_InvalidBaseRef(t *testing.T) {
	tests := []struct {
		name    string
		baseRef string
		wantErr string
	}{
		{
			name:    "baseRef starting with dash",
			baseRef: "--exec=malicious",
			wantErr: "invalid base ref",
		},
		{
			name:    "baseRef starting with single dash",
			baseRef: "-v",
			wantErr: "invalid base ref",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := resolveChangedFiles(tt.baseRef)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error = %q, want it to contain %q", err.Error(), tt.wantErr)
			}
		})
	}
}
