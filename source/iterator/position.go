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

	"github.com/conduitio/conduit-commons/opencdc"
)

// Mode defines an iterator mode.
type mode string

const (
	modeSnapshot mode = "snapshot"
	modeCDC      mode = "cdc"
)

// Position represents Oracle position.
type Position struct {
	// Mode represents current iterator mode.
	IteratorMode mode `json:"mode"`

	// CreatedAt is the Unix time from which the system should receive events of the resource in the CDC iterator.
	CreatedAt int64 `json:"created_at"`

	// Cursor is the resource or event identifier for receiving shifted data in the following requests.
	Cursor string `json:"cursor"`

	// Index is the current index of the returning record from the batch of previously received resources.
	Index int `json:"index"`
}

// ParseSDKPosition parses opencdc.Position and returns Position.
func ParseSDKPosition(position opencdc.Position) (*Position, error) {
	if position == nil {
		return &Position{
			IteratorMode: modeSnapshot,
			CreatedAt:    time.Now().Unix(),
		}, nil
	}

	pos := Position{}

	err := json.Unmarshal(position, &pos)
	if err != nil {
		return nil, fmt.Errorf("unmarshal opencdc.Position into Position: %w", err)
	}

	return &pos, nil
}

// marshalPosition marshals Position and returns opencdc.Position or an error.
func (p Position) marshalPosition() (opencdc.Position, error) {
	positionBytes, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("marshal position: %w", err)
	}

	return positionBytes, nil
}
