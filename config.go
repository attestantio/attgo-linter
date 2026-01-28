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

// Config holds the configuration for the attgo linter plugin.
type Config struct {
	// HIGH PRIORITY - enabled by default
	EnableNoPkgLogger bool `json:"enable_no_pkg_logger"`
	EnableEnumIota    bool `json:"enable_enum_iota"`
	EnableCurrentYear bool `json:"enable_current_year"`

	// MEDIUM PRIORITY - disabled by default
	EnableCapitalComment bool `json:"enable_capital_comment"`
	EnableFuncOpts       bool `json:"enable_func_opts"`
	EnableRawString      bool `json:"enable_raw_string"`

	// LOW PRIORITY - disabled by default
	EnableStructFieldOrder bool `json:"enable_struct_field_order"`
	EnableInterfaceCheck   bool `json:"enable_interface_check"`

	// LoggerTypePatterns specifies the type patterns to detect as loggers.
	// Default patterns include common logging libraries.
	LoggerTypePatterns []string `json:"logger_type_patterns"`

	// EnumTypeSuffixes specifies the suffixes that identify enum types.
	// Default: ["Type", "Status", "State", "Kind", "Mode"]
	EnumTypeSuffixes []string `json:"enum_type_suffixes"`
}

// DefaultConfig returns a Config with sensible defaults.
// HIGH priority rules are enabled by default.
func DefaultConfig() *Config {
	return &Config{
		// HIGH PRIORITY - enabled by default
		EnableNoPkgLogger: true,
		EnableEnumIota:    true,
		EnableCurrentYear: true,

		// MEDIUM PRIORITY - disabled by default
		EnableCapitalComment: false,
		EnableFuncOpts:       false,
		EnableRawString:      false,

		// LOW PRIORITY - disabled by default
		EnableStructFieldOrder: false,
		EnableInterfaceCheck:   false,

		// Default logger patterns
		LoggerTypePatterns: []string{
			"zerolog.Logger",
			"*zerolog.Logger",
			"zap.Logger",
			"*zap.Logger",
			"zap.SugaredLogger",
			"*zap.SugaredLogger",
			"logrus.Logger",
			"*logrus.Logger",
			"logrus.Entry",
			"*logrus.Entry",
			"slog.Logger",
			"*slog.Logger",
			"log.Logger",
			"*log.Logger",
		},

		// Default enum suffixes
		EnumTypeSuffixes: []string{
			"Type",
			"Status",
			"State",
			"Kind",
			"Mode",
		},
	}
}

// Merge applies non-zero values from other to c.
func (c *Config) Merge(other *Config) {
	if other == nil {
		return
	}

	// For boolean fields, we need explicit handling since false is the zero value.
	// The golangci-lint plugin system passes only explicitly set values,
	// so we check if the value differs from what would be "unset".
	// This is handled by the plugin initialization.

	if len(other.LoggerTypePatterns) > 0 {
		c.LoggerTypePatterns = other.LoggerTypePatterns
	}

	if len(other.EnumTypeSuffixes) > 0 {
		c.EnumTypeSuffixes = other.EnumTypeSuffixes
	}
}
