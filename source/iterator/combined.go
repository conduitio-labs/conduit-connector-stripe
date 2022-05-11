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
	"time"

	sdk "github.com/conduitio/conduit-connector-sdk"

	"github.com/conduitio/conduit-connector-stripe/source/position"
)

const idKey = "id"

// A Combined represents a struct of combined iterator.
type Combined struct {
	snapshot *Snapshot
	cdc      *CDC
	position *position.Position
}

// New initializes a combined iterator.
func New(stripeSvc Stripe, pos *position.Position, pollingPeriod time.Duration) *Combined {
	combined := &Combined{
		snapshot: NewSnapshot(stripeSvc, pos),
		cdc:      NewCDC(stripeSvc, pos, pollingPeriod),
		position: pos,
	}

	return combined
}

// Next returns the next record.
func (iter *Combined) Next() (sdk.Record, error) {
	switch iter.position.IteratorType {
	case position.SnapshotType:
		r, err := iter.snapshot.Next()
		if err != nil {
			return sdk.Record{}, err
		}

		if r.Payload == nil {
			iter.position.IteratorType = position.CDCType
			iter.position.Cursor = ""

			r, err = iter.cdc.Next()
		}

		return r, err
	case position.CDCType:
		return iter.cdc.Next()
	}

	return sdk.Record{}, nil
}

// Stop stops the iterator.
func (iter *Combined) Stop() error {
	if iter.position.IteratorType == position.CDCType {
		return iter.cdc.Stop()
	}

	return nil
}
