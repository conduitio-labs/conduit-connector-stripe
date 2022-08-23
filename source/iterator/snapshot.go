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

	"github.com/conduitio-labs/conduit-connector-stripe/models"
	sdk "github.com/conduitio/conduit-connector-sdk"
)

// A Snapshot represents a struct of snapshot iterator.
type Snapshot struct {
	stripeSvc Stripe
	position  *Position
	response  *models.ResourceResponse
	index     int
}

// NewSnapshot initializes snapshot iterator.
func NewSnapshot(stripeSvc Stripe, pos *Position) *Snapshot {
	return &Snapshot{
		stripeSvc: stripeSvc,
		position:  pos,
	}
}

// Next returns the next record.
// Note: The `Snapshot` iterator creates a copy of the data, which is sorted by date of creation in descending order.
func (i *Snapshot) Next() (sdk.Record, error) {
	if i.response == nil || len(i.response.Data) == i.index {
		if err := i.refreshData(); err != nil {
			return sdk.Record{}, fmt.Errorf("populate with the resource: %w", err)
		}

		// if there is no data - go to `CDC` iterator
		if len(i.response.Data) == 0 {
			i.position.IteratorMode = modeCDC
			i.position.Cursor = ""

			return sdk.Record{}, nil
		}
	}

	i.position.Cursor = i.response.Data[i.index][models.KeyID].(string)

	position, err := i.position.marshalPosition()
	if err != nil {
		return sdk.Record{}, fmt.Errorf("build record position: %w", err)
	}

	payload, err := i.buildRecordPayload()
	if err != nil {
		return sdk.Record{}, fmt.Errorf("build record payload: %w", err)
	}

	record := sdk.Util.Source.NewRecordSnapshot(
		position,
		i.buildRecordMetadata(),
		i.buildRecordKey(),
		payload,
	)

	i.index++

	return record, nil
}

// refreshData receives the resource data from Stripe, and assigns them to the iterator.
func (i *Snapshot) refreshData() error {
	resp, err := i.stripeSvc.GetResource(i.position.Cursor)
	if err != nil {
		return fmt.Errorf("get list of resource objects: %w", err)
	}

	i.response = &resp
	i.index = 0

	return nil
}

// buildRecordMetadata returns the metadata for the record.
func (i *Snapshot) buildRecordMetadata() map[string]string {
	metadata := make(sdk.Metadata, 1)

	createdAt := time.Now()
	if c, ok := i.response.Data[i.index][models.KeyCreated].(float64); ok {
		createdAt = time.Unix(int64(c), 0)
	}
	metadata.SetCreatedAt(createdAt)

	return metadata
}

// buildRecordKey returns the key for the record.
func (i *Snapshot) buildRecordKey() sdk.Data {
	return sdk.StructuredData{
		models.KeyID: i.response.Data[i.index][models.KeyID].(string),
	}
}

// buildRecordPayload returns the payload for the record.
func (i *Snapshot) buildRecordPayload() (sdk.Data, error) {
	payload, err := json.Marshal(i.response.Data[i.index])
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	return sdk.RawData(payload), nil
}
