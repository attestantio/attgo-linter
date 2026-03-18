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
	"testing"
)

func TestDefaultConfig_CurrentYearBaseRef(t *testing.T) {
	c := DefaultConfig()
	if c.CurrentYearBaseRef != "origin/main" {
		t.Errorf("expected CurrentYearBaseRef to be %q, got %q", "origin/main", c.CurrentYearBaseRef)
	}
}

func TestConfig_Merge_CurrentYearBaseRef(t *testing.T) {
	tests := []struct {
		name     string
		base     string
		other    string
		nilOther bool
		expected string
	}{
		{
			name:     "propagates non-empty base ref",
			base:     "",
			other:    "origin/main",
			expected: "origin/main",
		},
		{
			name:     "does not overwrite with empty",
			base:     "origin/develop",
			other:    "",
			expected: "origin/develop",
		},
		{
			name:     "overwrites with different value",
			base:     "origin/develop",
			other:    "origin/main",
			expected: "origin/main",
		},
		{
			name:     "nil other does not panic",
			base:     "origin/main",
			nilOther: true,
			expected: "origin/main",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := DefaultConfig()
			c.CurrentYearBaseRef = tt.base

			if tt.nilOther {
				c.Merge(nil)
			} else {
				c.Merge(&Config{CurrentYearBaseRef: tt.other})
			}

			if c.CurrentYearBaseRef != tt.expected {
				t.Errorf("CurrentYearBaseRef = %q, want %q", c.CurrentYearBaseRef, tt.expected)
			}
		})
	}
}
