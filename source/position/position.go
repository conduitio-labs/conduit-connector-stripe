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

package position

import (
	"encoding/json"
	"fmt"
	"time"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/conduitio/conduit-connector-stripe/models"
)

// A Position represents a Stripe position.
type Position struct {
	// IteratorType is the type of the current iterator.
	IteratorType models.IteratorType `json:"iterator_type" validate:"required,iterator_type"`

	// CreatedAt is the Unix time from which the system should receive events of the resource in the CDC iterator.
	CreatedAt int64 `json:"created_at" validate:"required"`

	// Cursor is the resource or event identifier for receiving shifted data in the following requests.
	Cursor string `json:"cursor"`

	// Index is the current index of the returning record from the bunch of previously received resources.
	Index int `json:"index"`
}

// ParseSDKPosition unmarshal sdk.Position and returns Position.
func ParseSDKPosition(rp sdk.Position) (Position, error) {
	if rp == nil {
		return Position{
			IteratorType: models.SnapshotIterator,
			CreatedAt:    time.Now().Unix(),
		}, nil
	}

	pos := Position{}

	err := json.Unmarshal(rp, &pos)
	if err != nil {
		return Position{}, fmt.Errorf("failed to unmarshal position: %w", err)
	}

	err = pos.Validate()
	if err != nil {
		return Position{}, err
	}

	return pos, nil
}

// FormatSDKPosition marshals Position and returns sdk.Position.
func (p Position) FormatSDKPosition() (sdk.Position, error) {
	err := p.Validate()
	if err != nil {
		return nil, err
	}

	rp, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal position: %w", err)
	}

	return rp, nil
}
