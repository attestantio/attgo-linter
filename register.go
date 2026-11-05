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

package attgolinter

import (
	"fmt"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

// RegisterAnalyzer registers analyzer as a golangci-lint module plugin named
// after it. golangci-lint reports every analyzer of a plugin under the plugin's
// name, so //nolint:attgo silences all attgo rules. A plugin of its own lets a
// repository enable the analyzer, and //nolint:<analyzer.Name> silence it,
// alone.
//
// Call it from an init function. A repository opts in by adding a
// linters.settings.custom entry named analyzer.Name, which takes no settings,
// and, when it sets linters.default to none or fast, by listing the name under
// linters.enable. The name must be lowercase, because golangci-lint lowercases
// the names in //nolint directives.
//
// An analyzer value belongs to one plugin. golangci-lint maps each analyzer
// value to one linter name, so a value that the attgo plugin also builds
// reports under only one of the two names when both are enabled.
func RegisterAnalyzer(analyzer *analysis.Analyzer) {
	register.Plugin(analyzer.Name, func(settings any) (register.LinterPlugin, error) {
		_, err := register.DecodeSettings[struct{}](settings)
		if err != nil {
			return nil, fmt.Errorf("%s takes no settings: %w", analyzer.Name, err)
		}

		return &analyzerPlugin{analyzer: analyzer}, nil
	})
}

// analyzerPlugin is the plugin of a single analyzer.
type analyzerPlugin struct {
	analyzer *analysis.Analyzer
}

// BuildAnalyzers returns the analyzer.
func (p *analyzerPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{p.analyzer}, nil
}

// GetLoadMode returns the load mode required by the plugin.
func (p *analyzerPlugin) GetLoadMode() string {
	// Every analyzer loads type information. Add an option when one needs only syntax.
	return register.LoadModeTypesInfo
}
