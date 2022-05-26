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

package iterator

import (
	"fmt"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/conduitio/conduit-connector-stripe/models"
	"github.com/conduitio/conduit-connector-stripe/source/position"
)

const idKey = "id"

// An Iterator represents a struct of iterator.
type Iterator struct {
	snapshot *SnapshotIterator
	cdc      *CDCIterator
	position *position.Position
}

// NewIterator initializes an iterator.
func NewIterator(stripeSvc Stripe, pos *position.Position) *Iterator {
	iterator := &Iterator{
		position: pos,
		cdc:      NewCDCIterator(stripeSvc, pos),
	}

	if pos.IteratorType == models.SnapshotIterator {
		iterator.snapshot = NewSnapshotIterator(stripeSvc, pos)
	}

	return iterator
}

// Next returns the next record.
func (iter *Iterator) Next() (sdk.Record, error) {
	switch iter.position.IteratorType {
	case models.SnapshotIterator:
		return iter.snapshot.Next()
	case models.CDCIterator:
		return iter.cdc.Next()
	}

	return sdk.Record{}, fmt.Errorf("unexpected iterator type: %s", iter.position.IteratorType)
}
