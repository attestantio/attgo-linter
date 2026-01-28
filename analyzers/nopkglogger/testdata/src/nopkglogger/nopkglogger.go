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

package nopkglogger

import "nopkglogger/zerolog"

// Bad: package-level logger variable.
var log zerolog.Logger // want `package-level logger "log" detected; loggers should be struct fields for better dependency injection and testability`

// Bad: package-level logger pointer.
var logPtr *zerolog.Logger // want `package-level logger "logPtr" detected; loggers should be struct fields for better dependency injection and testability`

// Bad: multiple loggers in one declaration.
var (
	logger1 zerolog.Logger  // want `package-level logger "logger1" detected; loggers should be struct fields for better dependency injection and testability`
	logger2 *zerolog.Logger // want `package-level logger "logger2" detected; loggers should be struct fields for better dependency injection and testability`
)

// Good: non-logger package variables are fine.
var (
	version = "1.0.0"
	counter int
)

// Good: logger as struct field.
type Service struct {
	log zerolog.Logger
}

// Good: function-local logger is fine.
func doSomething() {
	var localLog zerolog.Logger
	_ = localLog
}
