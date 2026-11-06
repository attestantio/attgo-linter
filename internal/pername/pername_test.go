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

package pername_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/attestantio/attgo-linter/internal/pername"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

// newTestAnalyzer returns an analyzer with the given name that reports nothing.
func newTestAnalyzer(name string) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: name,
		Doc:  "analyzer registered by a test",
		Run:  func(*analysis.Pass) (any, error) { return nil, nil },
	}
}

func TestRegister_Plugin(t *testing.T) {
	analyzer := newTestAnalyzer("attgo_registered_plugin")
	pername.Register(analyzer)

	newPlugin, err := register.GetPlugin("attgo_registered_plugin")
	if err != nil {
		t.Fatalf("expected the plugin to be registered under the analyzer name: %v", err)
	}

	plugin, err := newPlugin(nil)
	if err != nil {
		t.Fatalf("failed to create the plugin: %v", err)
	}

	analyzers, err := plugin.BuildAnalyzers()
	if err != nil {
		t.Fatalf("failed to build the analyzers: %v", err)
	}

	if len(analyzers) != 1 || analyzers[0] != analyzer {
		t.Errorf("expected the plugin to build exactly the registered analyzer, got %v", analyzers)
	}

	// Without type information, analyzers that read pass.TypesInfo would see nil.
	if plugin.GetLoadMode() != register.LoadModeTypesInfo {
		t.Errorf("load mode = %q, want %q", plugin.GetLoadMode(), register.LoadModeTypesInfo)
	}
}

func TestRegister_Settings(t *testing.T) {
	tests := []struct {
		name     string
		settings any
		wantErr  string
		// wantCause is the rejected setting that the wrapped error names.
		wantCause string
	}{
		{
			name:     "no settings",
			settings: nil,
		},
		{
			name:     "empty settings",
			settings: map[string]any{},
		},
		{
			name:      "enable flag",
			settings:  map[string]any{"enable": false},
			wantErr:   `attgo_registered_settings takes no settings`,
			wantCause: `enable`,
		},
	}

	pername.Register(newTestAnalyzer("attgo_registered_settings"))

	newPlugin, err := register.GetPlugin("attgo_registered_settings")
	if err != nil {
		t.Fatalf("expected the plugin to be registered under the analyzer name: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := newPlugin(tt.settings)

			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				return
			}

			if err == nil {
				t.Fatalf("expected an error containing %q, got none", tt.wantErr)
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error = %q, want it to contain %q", err.Error(), tt.wantErr)
			}

			if !strings.Contains(err.Error(), tt.wantCause) {
				t.Errorf("error = %q, want it to name the rejected setting %q", err.Error(), tt.wantCause)
			}

			if errors.Unwrap(err) == nil {
				t.Error("expected the settings decoding error to be wrapped")
			}
		})
	}
}
