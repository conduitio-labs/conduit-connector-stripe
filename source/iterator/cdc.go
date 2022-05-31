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

// A CDCIterator represents a struct of cdc iterator.
type CDCIterator struct {
	stripeSvc Stripe
	position  *position.Position

	// eventData is a slice of the event data from the Stripe response.
	eventData []models.EventData
}

// NewCDCIterator initializes cdc iterator.
func NewCDCIterator(stripeSvc Stripe, pos *position.Position) *CDCIterator {
	return &CDCIterator{
		stripeSvc: stripeSvc,
		position:  pos,
	}
}

// Next returns the next record.
func (iter *CDCIterator) Next() (sdk.Record, error) {
	if iter.eventData == nil || iter.position.Index == 0 {
		if err := iter.getData(); err != nil {
			return sdk.Record{}, fmt.Errorf("get event data: %w", err)
		}

		if len(iter.eventData) == 0 {
			return sdk.Record{}, sdk.ErrBackoffRetry
		}
	}

	index := iter.position.Index

	iter.position.Index++

	// update `Cursor` in the position if it is the last element of the resulting slice
	if len(iter.eventData) == iter.position.Index {
		iter.position.Index = 0
		iter.position.Cursor = iter.eventData[len(iter.eventData)-1].ID
	}

	payload, err := json.Marshal(iter.eventData[index].Data.Object)
	if err != nil {
		return sdk.Record{}, fmt.Errorf("marshal payload: %w", err)
	}

	rp, err := iter.position.FormatSDKPosition()
	if err != nil {
		return sdk.Record{}, fmt.Errorf("format sdk position: %w", err)
	}

	output := sdk.Record{
		Position: rp,
		Metadata: map[string]string{
			models.ActionKey: models.EventsAction[iter.eventData[index].Type],
		},
		CreatedAt: time.Unix(iter.eventData[index].Created, 0),
		Key: sdk.StructuredData{
			idKey: iter.eventData[index].Data.Object[idKey].(string),
		},
		Payload: sdk.RawData(payload),
	}

	return output, nil
}

// getData calls methods to assign Stripe event data to the iterator.
// Two methods are used because the data sorted by date of creation in descending order and the `ending_before` shift is unknown.
func (iter *CDCIterator) getData() error {
	if iter.position.Cursor == "" {
		// because the data is sorted by date of creation in descending order
		// and the shift `ending_before` is not known, it takes all the data and reverses it
		return iter.getDataWithStartingAfter()
	}

	return iter.getDataWithEndingBefore()
}

// getDataWithStartingAfter makes requests with `starting_after` parameter
// to receive all the event data, and assigns Stripe event data to the iterator.
func (iter *CDCIterator) getDataWithStartingAfter() error {
	var (
		eventData []models.EventData

		startingAfter string
	)

	// get all the event data
	for {
		// receive the data with `starting_after` parameter
		resp, err := iter.stripeSvc.GetEvent(iter.position.CreatedAt, startingAfter, "")
		if err != nil {
			return fmt.Errorf("get list of event objects: %w", err)
		}

		if len(resp.Data) > 0 {
			// update startingAfter parameter for the next request
			startingAfter = resp.Data[len(resp.Data)-1].ID

			eventData = append(eventData, resp.Data...)
		}

		// break the loop if there is no more data
		if !resp.HasMore {
			break
		}
	}

	if len(eventData) > 0 {
		// do reverse, because Stripe receives data sorted by date of creation in descending order
		for i, j := 0, len(eventData)-1; i < j; i, j = i+1, j-1 {
			eventData[i], eventData[j] = eventData[j], eventData[i]
		}
	}

	iter.eventData = eventData

	return nil
}

// getDataWithEndingBefore makes requests with `ending_before` parameter
// and assigns Stripe event data to the iterator.
func (iter *CDCIterator) getDataWithEndingBefore() error {
	var eventData []models.EventData

	// receive the data with `ending_before` parameter
	resp, err := iter.stripeSvc.GetEvent(iter.position.CreatedAt, "", iter.position.Cursor)
	if err != nil {
		return fmt.Errorf("get list of event objects: %w", err)
	}

	if len(resp.Data) > 0 {
		// do reverse, because Stripe receives data sorted by date of creation in descending order
		for i, j := 0, len(resp.Data)-1; i < j; i, j = i+1, j-1 {
			resp.Data[i], resp.Data[j] = resp.Data[j], resp.Data[i]
		}

		eventData = append(eventData, resp.Data...)
	}

	iter.eventData = eventData

	return nil
}
