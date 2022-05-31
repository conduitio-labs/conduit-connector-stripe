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
	"encoding/json"
	"fmt"
	"time"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/conduitio/conduit-connector-stripe/models"
	"github.com/conduitio/conduit-connector-stripe/source/position"
)

// A SnapshotIterator represents a struct of snapshot iterator.
type SnapshotIterator struct {
	stripeSvc Stripe
	position  *position.Position
	response  *models.ResourceResponse
	index     int
}

// NewSnapshotIterator initializes snapshot iterator.
func NewSnapshotIterator(stripeSvc Stripe, pos *position.Position) *SnapshotIterator {
	return &SnapshotIterator{
		stripeSvc: stripeSvc,
		position:  pos,
	}
}

// Next returns the next record.
// Note: The `Snapshot` iterator creates a copy of the data, which is sorted by date of creation in descending order.
func (iter *SnapshotIterator) Next() (sdk.Record, error) {
	if iter.response == nil || len(iter.response.Data) == iter.index {
		if err := iter.refreshData(); err != nil {
			return sdk.Record{}, fmt.Errorf("populate with the resource: %w", err)
		}

		// if there is no data - go to `CDC` iterator
		if len(iter.response.Data) == 0 {
			iter.position.IteratorType = models.CDCIterator
			iter.position.Cursor = ""

			return sdk.Record{}, nil
		}
	}

	payload, err := json.Marshal(iter.response.Data[iter.index])
	if err != nil {
		return sdk.Record{}, fmt.Errorf("marshal payload: %w", err)
	}

	iter.position.Cursor = iter.response.Data[iter.index][idKey].(string)

	created, ok := iter.response.Data[iter.index]["created"].(float64)
	if !ok {
		created = float64(time.Now().Unix())
	}

	rp, err := iter.position.FormatSDKPosition()
	if err != nil {
		return sdk.Record{}, fmt.Errorf("format sdk position: %w", err)
	}

	output := sdk.Record{
		Position: rp,
		Metadata: map[string]string{
			models.ActionKey: models.InsertAction,
		},
		CreatedAt: time.Unix(int64(created), 0),
		Key: sdk.StructuredData{
			idKey: iter.response.Data[iter.index][idKey].(string),
		},
		Payload: sdk.RawData(payload),
	}

	iter.index++

	return output, nil
}

// refreshData receives the resource data from Stripe, and assigns them to the iterator.
func (iter *SnapshotIterator) refreshData() error {
	resp, err := iter.stripeSvc.GetResource(iter.position.Cursor)
	if err != nil {
		return fmt.Errorf("get list of resource objects: %w", err)
	}

	iter.response = &resp
	iter.index = 0

	return nil
}
