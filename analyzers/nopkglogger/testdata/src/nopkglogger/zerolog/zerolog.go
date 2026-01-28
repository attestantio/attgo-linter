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

// Package zerolog is a mock zerolog package for testing.
package zerolog

// Logger is a mock logger type.
type Logger struct{}

// Info returns a mock event.
func (l Logger) Info() *Event { return &Event{} }

// Event is a mock event type.
type Event struct{}

// Msg logs a message.
func (e *Event) Msg(string) {}
