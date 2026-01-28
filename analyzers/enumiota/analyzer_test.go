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

package enumiota_test

import (
	"testing"

	"github.com/attestantio/attgo-linter/analyzers/enumiota"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()

	enumSuffixes := []string{
		"Type",
		"Status",
		"State",
		"Kind",
		"Mode",
	}

	analyzer := enumiota.NewAnalyzer(enumSuffixes)

	analysistest.Run(t, testdata, analyzer, "enumiota")
}
