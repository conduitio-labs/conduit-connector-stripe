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

// A CDC represents a struct of cdc iterator.
type CDC struct {
	stripeSvc Stripe
	position  *Position

	// eventData is a slice of the event data from the Stripe response.
	eventData []models.EventData
}

// NewCDC initializes cdc iterator.
func NewCDC(stripeSvc Stripe, pos *Position) *CDC {
	return &CDC{
		stripeSvc: stripeSvc,
		position:  pos,
	}
}

// Next returns the next record.
func (i *CDC) Next() (sdk.Record, error) {
	if i.eventData == nil || i.position.Index == 0 {
		if err := i.getData(); err != nil {
			return sdk.Record{}, fmt.Errorf("get event data: %w", err)
		}

		if len(i.eventData) == 0 {
			return sdk.Record{}, sdk.ErrBackoffRetry
		}
	}

	metadata := i.buildRecordMetadata()

	key := i.buildRecordKey()

	payload, err := i.buildRecordPayload()
	if err != nil {
		return sdk.Record{}, fmt.Errorf("build record payload: %w", err)
	}

	index := i.position.Index

	i.position.Index++

	// update `Cursor` in the position if it is the last element of the resulting slice
	if len(i.eventData) == i.position.Index {
		i.position.Index = 0
		i.position.Cursor = i.eventData[len(i.eventData)-1].ID
	}

	position, err := i.position.marshalPosition()
	if err != nil {
		return sdk.Record{}, fmt.Errorf("build record position: %w", err)
	}

	// there is no default case, because sdk.OperationUpdate is the default operation
	switch models.EventsOperation[i.eventData[index].Type] {
	case sdk.OperationCreate:
		return sdk.Util.Source.NewRecordCreate(
			position,
			metadata,
			key,
			payload,
		), nil
	case sdk.OperationUpdate:
		return sdk.Util.Source.NewRecordUpdate(
			position,
			metadata,
			key,
			nil,
			payload,
		), nil
	case sdk.OperationDelete:
		return sdk.Util.Source.NewRecordDelete(
			position,
			metadata,
			key,
		), nil
	}

	return sdk.Record{}, nil
}

// getData calls methods to assign Stripe event data to the iterator.
func (i *CDC) getData() error {
	if i.position.Cursor == "" {
		// because the data is sorted by date of creation in descending order
		// and the shift `ending_before` is not known, it takes all the data and reverses it
		return i.getDataWithStartingAfter()
	}

	return i.getDataWithEndingBefore()
}

// getDataWithStartingAfter makes requests with `starting_after` parameter
// to receive all the event data, and assigns Stripe event data to the iterator.
func (i *CDC) getDataWithStartingAfter() error {
	var (
		eventsData models.EventsData

		startingAfter string
	)

	// get all the event data
	for {
		// receive the data with `starting_after` parameter
		resp, err := i.stripeSvc.GetEvent(i.position.CreatedAt, startingAfter, "")
		if err != nil {
			return fmt.Errorf("get list of event objects: %w", err)
		}

		if len(resp.Data) > 0 {
			// update startingAfter parameter for the next request
			startingAfter = resp.Data[len(resp.Data)-1].ID

			eventsData = append(eventsData, resp.Data...)
		}

		// break the loop if there is no more data
		if !resp.HasMore {
			break
		}
	}

	if len(eventsData) > 0 {
		// do reverse, because Stripe receives data sorted by date of creation in descending order
		eventsData.Reverse()
	}

	i.eventData = eventsData

	return nil
}

// getDataWithEndingBefore makes requests with `ending_before` parameter
// and assigns Stripe event data to the iterator.
func (i *CDC) getDataWithEndingBefore() error {
	var eventsData models.EventsData

	// receive the data with `ending_before` parameter
	resp, err := i.stripeSvc.GetEvent(i.position.CreatedAt, "", i.position.Cursor)
	if err != nil {
		return fmt.Errorf("get list of event objects: %w", err)
	}

	if len(resp.Data) > 0 {
		// do reverse, because Stripe receives data sorted by date of creation in descending order
		resp.Data.Reverse()

		eventsData = append(eventsData, resp.Data...)
	}

	i.eventData = eventsData

	return nil
}

// buildRecordMetadata returns the metadata for the record.
func (i *CDC) buildRecordMetadata() map[string]string {
	metadata := sdk.Metadata{}

	metadata.SetCreatedAt(time.Unix(i.eventData[i.position.Index].Created, 0))

	return metadata
}

// buildRecordKey returns the key for the record.
func (i *CDC) buildRecordKey() sdk.Data {
	return sdk.StructuredData{
		models.KeyID: i.eventData[i.position.Index].Data.Object[models.KeyID].(string),
	}
}

// buildRecordPayload returns the payload for the record.
func (i *CDC) buildRecordPayload() (sdk.Data, error) {
	payload, err := json.Marshal(i.eventData[i.position.Index].Data.Object)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	return sdk.RawData(payload), nil
}
