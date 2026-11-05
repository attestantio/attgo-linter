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

package project

// Every literal below is reported by attgo_raw_string (six escape sequences)
// and by both sample analyzers. The directive on its line decides which of
// the three linters stay silent. The last variable sits under a comment that
// attgo_capital_comment reports only where the config turns that rule on.
//
// Keep this file free of imports, see MAINTENANCE.md.
var (
	plain = "flag-me \"a\" \"b\" \"c\""

	byFirst = "flag-me \"a\" \"b\" \"c\"" //nolint:sample_first

	bySecond = "flag-me \"a\" \"b\" \"c\"" //nolint:sample_second

	byAttgo = "flag-me \"a\" \"b\" \"c\"" //nolint:attgo

	byBoth = "flag-me \"a\" \"b\" \"c\"" //nolint:sample_first,sample_second

	byRuleName = "flag-me \"a\" \"b\" \"c\"" //nolint:attgo_raw_string

	// lowercase on purpose.
	capitalComment = 1
)
