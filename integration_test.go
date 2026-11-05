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

package attgolinter_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// TestIntegration_PerNameNolint builds a custom golangci-lint binary from this
// checkout, with two sample analyzers registered under their own names next to
// the attgo plugin, and checks which linters still report once a //nolint
// directive names one of them. It lints with the attgo settings of each
// consumer (vouch and dirk share theirs), so it also checks that those
// configurations keep working next to per-name plugins.
func TestIntegration_PerNameNolint(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping custom binary build in short mode")
	}

	if _, err := exec.LookPath("golangci-lint"); err != nil {
		if os.Getenv("CI") != "" {
			t.Fatal("golangci-lint is not on PATH: the test workflow installs it")
		}

		t.Skip("skipping: `golangci-lint custom` needs golangci-lint on PATH")
	}

	// Tests run in the package directory, which is the module root.
	checkout, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	work := t.TempDir()
	project := filepath.Join(work, "project")

	err = os.CopyFS(project, os.DirFS(filepath.Join(checkout, "testdata", "project")))
	if err != nil {
		t.Fatalf("failed to copy the project to lint: %v", err)
	}

	// The first plugin is the module as consumers use it. The second adds the
	// sample analyzers from the same checkout.
	customConfig := fmt.Sprintf(`version: %[1]s
name: custom-gcl
destination: %[2]q
plugins:
  - module: github.com/attestantio/attgo-linter
    path: %[3]q
  - module: github.com/attestantio/attgo-linter
    import: github.com/attestantio/attgo-linter/testdata/sampleplugin
    path: %[3]q
`, golangciLintVersion(t, checkout), work, checkout)

	err = os.WriteFile(filepath.Join(work, ".custom-gcl.yml"), []byte(customConfig), 0o600)
	if err != nil {
		t.Fatalf("failed to write .custom-gcl.yml: %v", err)
	}

	out, err := runCommand(work, "golangci-lint", "custom")
	if err != nil {
		t.Fatalf("failed to build the custom binary: %v\n%s", err, out)
	}

	lineOf := variableLines(t, filepath.Join(project, "project.go"))

	tests := []struct {
		name     string
		variable string
		expected []string
	}{
		{
			name:     "no directive leaves every linter reporting",
			variable: "plain",
			expected: []string{"attgo", "sample_first", "sample_second"},
		},
		{
			name:     "a per-name directive silences only that analyzer",
			variable: "byFirst",
			expected: []string{"attgo", "sample_second"},
		},
		{
			name:     "the same holds for the other analyzer",
			variable: "bySecond",
			expected: []string{"attgo", "sample_first"},
		},
		{
			name:     "the attgo directive silences only attgo",
			variable: "byAttgo",
			expected: []string{"sample_first", "sample_second"},
		},
		{
			name:     "a directive can name several linters",
			variable: "byBoth",
			expected: []string{"attgo"},
		},
		{
			name:     "an attgo rule name is not a linter name",
			variable: "byRuleName",
			expected: []string{"attgo", "sample_first", "sample_second"},
		},
	}

	for _, config := range []string{"vouch.yml", "go-eth2-client.yml"} {
		t.Run(config, func(t *testing.T) {
			lintersByLine := runLint(t, filepath.Join(work, "custom-gcl"), project, config)

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					line, found := lineOf[tt.variable]
					if !found {
						t.Fatalf("variable %s not found in project.go", tt.variable)
					}

					got := lintersByLine[line]
					slices.Sort(got)

					if !slices.Equal(got, tt.expected) {
						t.Errorf("linters reporting on line %d (%s) = %v, want %v", line, tt.variable, got, tt.expected)
					}
				})
			}
		})
	}
}

// golangciLintVersion returns the golangci-lint version of the repository's
// .custom-gcl.yml, so the test builds the version the repository itself lints
// with.
func golangciLintVersion(t *testing.T, checkout string) string {
	t.Helper()

	data, err := os.ReadFile(filepath.Join(checkout, ".custom-gcl.yml"))
	if err != nil {
		t.Fatalf("failed to read .custom-gcl.yml: %v", err)
	}

	for _, line := range strings.Split(string(data), "\n") {
		version, found := strings.CutPrefix(line, "version:")
		if found {
			return strings.TrimSpace(version)
		}
	}

	t.Fatal("no version found in .custom-gcl.yml")

	return ""
}

// runCommand runs a command in dir and returns its combined output.
func runCommand(dir, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()

	return string(out), err
}

// runLint runs the custom binary in dir with the given configuration file and
// returns the names of the linters that reported on each line.
func runLint(t *testing.T, binary, dir, config string) map[int][]string {
	t.Helper()

	reportPath := filepath.Join(dir, config+".json")

	out, err := runCommand(dir, binary, "run",
		"--config="+config,
		"--output.json.path="+reportPath,
		// Report every linter on a line, not only the first.
		"--uniq-by-line=false", "--max-issues-per-linter=0", "--max-same-issues=0",
	)

	// Exit status 1 means findings. Anything else is a load failure or a pass.
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 {
		t.Fatalf("expected the lint run to report issues (exit status 1), got %v\n%s", err, out)
	}

	data, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("failed to read the JSON report: %v", err)
	}

	var report struct {
		Issues []struct {
			FromLinter string
			Pos        struct{ Line int }
		}
	}

	err = json.Unmarshal(data, &report)
	if err != nil {
		t.Fatalf("failed to parse the JSON report: %v", err)
	}

	res := map[int][]string{}
	for _, issue := range report.Issues {
		res[issue.Pos.Line] = append(res[issue.Pos.Line], issue.FromLinter)
	}

	return res
}

// variableLines maps each variable declared in a one-per-line `name = value`
// block of the file to its line number.
func variableLines(t *testing.T, path string) map[string]int {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}

	res := map[string]int{}

	for i, line := range strings.Split(string(data), "\n") {
		name, _, found := strings.Cut(strings.TrimSpace(line), " = ")
		if found {
			res[name] = i + 1
		}
	}

	return res
}
