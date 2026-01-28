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

// Package attgolinter provides a golangci-lint plugin for Attestant organization
// Go style enforcement.
package attgolinter

import (
	"encoding/json"

	"github.com/attestantio/attgo-linter/analyzers/capitalcomment"
	"github.com/attestantio/attgo-linter/analyzers/currentyear"
	"github.com/attestantio/attgo-linter/analyzers/enumiota"
	"github.com/attestantio/attgo-linter/analyzers/funcopts"
	"github.com/attestantio/attgo-linter/analyzers/interfacecheck"
	"github.com/attestantio/attgo-linter/analyzers/nopkglogger"
	"github.com/attestantio/attgo-linter/analyzers/rawstring"
	"github.com/attestantio/attgo-linter/analyzers/structfieldorder"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("attgo", New)
}

// Plugin implements the golangci-lint module plugin interface.
type Plugin struct {
	cfg *Config
}

// New creates a new attgo linter plugin with the given settings.
func New(settings any) (register.LinterPlugin, error) {
	cfg := DefaultConfig()

	if settings != nil {
		// Settings come as map[string]any, marshal/unmarshal to apply.
		data, err := json.Marshal(settings)
		if err != nil {
			return nil, err
		}

		var userCfg Config
		if err := json.Unmarshal(data, &userCfg); err != nil {
			return nil, err
		}

		// Apply user configuration.
		// For booleans, we need to handle the explicit false case.
		// Re-unmarshal to check which fields were explicitly set.
		var rawSettings map[string]any
		if err := json.Unmarshal(data, &rawSettings); err != nil {
			return nil, err
		}

		if _, ok := rawSettings["enable_no_pkg_logger"]; ok {
			cfg.EnableNoPkgLogger = userCfg.EnableNoPkgLogger
		}
		if _, ok := rawSettings["enable_enum_iota"]; ok {
			cfg.EnableEnumIota = userCfg.EnableEnumIota
		}
		if _, ok := rawSettings["enable_current_year"]; ok {
			cfg.EnableCurrentYear = userCfg.EnableCurrentYear
		}
		if _, ok := rawSettings["enable_capital_comment"]; ok {
			cfg.EnableCapitalComment = userCfg.EnableCapitalComment
		}
		if _, ok := rawSettings["enable_func_opts"]; ok {
			cfg.EnableFuncOpts = userCfg.EnableFuncOpts
		}
		if _, ok := rawSettings["enable_raw_string"]; ok {
			cfg.EnableRawString = userCfg.EnableRawString
		}
		if _, ok := rawSettings["enable_struct_field_order"]; ok {
			cfg.EnableStructFieldOrder = userCfg.EnableStructFieldOrder
		}
		if _, ok := rawSettings["enable_interface_check"]; ok {
			cfg.EnableInterfaceCheck = userCfg.EnableInterfaceCheck
		}

		cfg.Merge(&userCfg)
	}

	return &Plugin{cfg: cfg}, nil
}

// BuildAnalyzers returns the analyzers to run based on configuration.
func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	var analyzers []*analysis.Analyzer

	// HIGH PRIORITY (enabled by default)
	if p.cfg.EnableNoPkgLogger {
		analyzers = append(analyzers, nopkglogger.NewAnalyzer(p.cfg.LoggerTypePatterns))
	}
	if p.cfg.EnableEnumIota {
		analyzers = append(analyzers, enumiota.NewAnalyzer(p.cfg.EnumTypeSuffixes))
	}
	if p.cfg.EnableCurrentYear {
		analyzers = append(analyzers, currentyear.Analyzer)
	}

	// MEDIUM PRIORITY (disabled by default)
	if p.cfg.EnableCapitalComment {
		analyzers = append(analyzers, capitalcomment.Analyzer)
	}
	if p.cfg.EnableFuncOpts {
		analyzers = append(analyzers, funcopts.Analyzer)
	}
	if p.cfg.EnableRawString {
		analyzers = append(analyzers, rawstring.Analyzer)
	}

	// LOW PRIORITY (disabled by default)
	if p.cfg.EnableStructFieldOrder {
		analyzers = append(analyzers, structfieldorder.Analyzer)
	}
	if p.cfg.EnableInterfaceCheck {
		analyzers = append(analyzers, interfacecheck.Analyzer)
	}

	return analyzers, nil
}

// GetLoadMode returns the load mode required by the plugin.
// LoadModeTypesInfo is needed for type-aware analysis (logger detection, enum types).
func (p *Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
