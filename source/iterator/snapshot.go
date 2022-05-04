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

	"github.com/conduitio/conduit-connector-stripe/clients/http"
	"github.com/conduitio/conduit-connector-stripe/source/position"
)

const actionInsert = "insert"

// A SnapshotIterator represents iteration over a slice of Stripe data.
type SnapshotIterator struct {
	httpClientSvc http.HTTP
	position      position.Position
	response      *http.StripeResponse
	index         int
}

// NewSnapshotIterator returns SnapshotIterator.
func NewSnapshotIterator(cliSvc http.HTTP, pos position.Position) *SnapshotIterator {
	return &SnapshotIterator{
		httpClientSvc: cliSvc,
		position:      pos,
	}
}

// Next returns the next record.
func (i *SnapshotIterator) Next() (sdk.Record, error) {
	if i.response == nil || len(i.response.Data) == i.index {
		if i.response != nil && !i.position.HasMore {
			return sdk.Record{}, nil
		}

		if err := i.getResources(); err != nil {
			return sdk.Record{}, fmt.Errorf("get response data: %w", err)
		}
	}

	payload, err := json.Marshal(i.response.Data[i.index])
	if err != nil {
		return sdk.Record{}, fmt.Errorf("marshal payload: %w", err)
	}

	output := sdk.Record{
		Position: i.position.FormatSDKPosition(),
		Metadata: map[string]string{
			"action": actionInsert,
		},
		CreatedAt: time.Unix(int64(i.response.Data[i.index]["created"].(float64)), 0),
		Key:       sdk.RawData(i.response.Data[i.index]["id"].(string)),
		Payload:   sdk.RawData(payload),
	}

	i.position.StartingAfter = i.response.Data[i.index]["id"].(string)
	i.index++

	return output, nil
}

func (i *SnapshotIterator) getResources() error {
	resp, err := i.httpClientSvc.GetResources(i.position.StartingAfter)
	if err != nil {
		return fmt.Errorf("get stripe resources: %w", err)
	}

	i.response = &resp
	i.position.HasMore = i.response.HasMore
	i.index = 0

	return nil
}
