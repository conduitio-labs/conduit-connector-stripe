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

// A CDC represents a struct of cdc iterator.
type CDC struct {
	stripeSvc Stripe
	position  *position.Position
	eventData []models.EventData
	index     int
}

// NewCDC initializes cdc iterator.
func NewCDC(stripeSvc Stripe, pos *position.Position) *CDC {
	return &CDC{
		stripeSvc: stripeSvc,
		position:  pos,
	}
}

// Next returns the next record.
func (iter *CDC) Next() (sdk.Record, error) {
	if iter.eventData == nil || len(iter.eventData) == iter.index {
		if err := iter.getEventData(); err != nil {
			return sdk.Record{}, fmt.Errorf("get event data: %w", err)
		}

		if len(iter.eventData) == 0 {
			return sdk.Record{}, sdk.ErrBackoffRetry
		}
	}

	payload, err := json.Marshal(iter.eventData[iter.index].Data.Object)
	if err != nil {
		return sdk.Record{}, fmt.Errorf("marshal payload: %w", err)
	}

	output := sdk.Record{
		Position: iter.position.FormatSDKPosition(),
		Metadata: map[string]string{
			models.ActionKey: models.EventsAction[iter.eventData[iter.index].Type],
		},
		CreatedAt: time.Unix(iter.eventData[iter.index].Created, 0),
		Key: sdk.StructuredData{
			idKey: iter.eventData[iter.index].Data.Object[idKey].(string),
		},
		Payload: sdk.RawData(payload),
	}

	iter.index++

	return output, nil
}

func (iter *CDC) getEventData() error {
	var (
		eventData []models.EventData

		hasMore = true
	)

	iter.index = 0

	for hasMore {
		resp, err := iter.stripeSvc.GetEvent(iter.position)
		if err != nil {
			return fmt.Errorf("get list of event objects: %w", err)
		}

		hasMore = resp.HasMore

		if len(resp.Data) > 0 {
			if iter.position.CreatedAt != 0 {
				iter.position.CreatedAt = 0
			}

			iter.position.Cursor = resp.Data[len(resp.Data)-1].ID

			eventData = append(eventData, resp.Data...)

			break
		}
	}

	if len(eventData) > 0 {
		for i, j := 0, len(eventData)-1; i < j; i, j = i+1, j-1 {
			eventData[i], eventData[j] = eventData[j], eventData[i]
		}
	}

	iter.eventData = eventData

	return nil
}
