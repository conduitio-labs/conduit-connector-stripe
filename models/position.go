// Copyright © 2022 Meroxa, Inc.
//
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

package models

// An IteratorType represents a type of iterators.
type IteratorType string

const (
	// A SnapshotIterator represents a snapshot iterator type.
	SnapshotIterator IteratorType = "snapshot"

	// A CDCIterator represents a cdc iterator type.
	CDCIterator IteratorType = "cdc"
)

// IteratorTypeMap contains valid iterators.
var IteratorTypeMap = map[IteratorType]struct{}{
	SnapshotIterator: {},
	CDCIterator:      {},
}
