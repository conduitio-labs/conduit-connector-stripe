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

	"github.com/conduitio-labs/conduit-connector-stripe/models"
	sdk "github.com/conduitio/conduit-connector-sdk"
)

// A Stripe defines the interface of methods.
type Stripe interface {
	GetResource(string) (models.ResourceResponse, error)
	GetEvent(createdAt int64, startingAfter, endingBefore string) (models.EventResponse, error)
}

// An Iterator represents a struct of iterator.
type Iterator struct {
	snapshot *Snapshot
	cdc      *CDC
	position *Position
}

// New initializes an iterator.
func New(stripeSvc Stripe, pos *Position, snapshot bool) *Iterator {
	iterator := &Iterator{
		position: pos,
		cdc:      NewCDC(stripeSvc, pos),
	}

	if !snapshot {
		pos.IteratorMode = modeCDC
	}

	if pos.IteratorMode == modeSnapshot {
		iterator.snapshot = NewSnapshot(stripeSvc, pos)
	}

	return iterator
}

// Next returns the next record.
func (iter *Iterator) Next() (sdk.Record, error) {
	switch iter.position.IteratorMode {
	case modeSnapshot:
		record, err := iter.snapshot.Next()
		if err != nil {
			return sdk.Record{}, err
		}

		if record.Key != nil {
			return record, nil
		}

		fallthrough
	case modeCDC:
		return iter.cdc.Next()
	}

	return sdk.Record{}, fmt.Errorf("unexpected iterator mode: %s", iter.position.IteratorMode)
}
