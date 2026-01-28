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

package enumiota

// Bad: string-based enum type.
type SANType string

const (
	SANTypeDNS   SANType = "dns"   // want `enum constant "SANTypeDNS" uses string value; consider using uint64 with iota pattern instead`
	SANTypeEmail SANType = "email" // want `enum constant "SANTypeEmail" uses string value; consider using uint64 with iota pattern instead`
)

// Bad: another string-based enum.
type RequestStatus string

const (
	RequestStatusPending  RequestStatus = "pending"  // want `enum constant "RequestStatusPending" uses string value; consider using uint64 with iota pattern instead`
	RequestStatusApproved RequestStatus = "approved" // want `enum constant "RequestStatusApproved" uses string value; consider using uint64 with iota pattern instead`
)

// Good: uint64-based enum with iota.
type DataKind uint64

const (
	DataKindUnknown DataKind = iota
	DataKindJSON
	DataKindXML
)

// Good: int-based enum with iota.
type ProcessState int

const (
	ProcessStateIdle ProcessState = iota
	ProcessStateRunning
	ProcessStateStopped
)

// Good: enum with explicit values (still integer).
type Priority uint64

const (
	PriorityLow    Priority = 10
	PriorityMedium Priority = 50
	PriorityHigh   Priority = 100
)

// Not an enum: regular string type without enum suffix.
type Name string

const (
	DefaultName Name = "default" // No warning - not an enum type suffix.
)

// Not an enum: type without enum suffix.
type Color string

const (
	ColorRed  Color = "red"  // No warning - Color doesn't have enum suffix.
	ColorBlue Color = "blue" // No warning.
)
