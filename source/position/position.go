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
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	sdk "github.com/conduitio/conduit-connector-sdk"

	"github.com/conduitio/conduit-connector-stripe/config"
)

var (
	errParseCreatedAt = errors.New("the third part of position must be an int64")
	errParseIndex     = errors.New("the fourth part of position must be an int")
)

const (
	SnapshotType = "s"
	CDCType      = "c"
)

const (
	positionFormat = "%s.%s.%s.%d.%d"

	resourceSliceIndex     = 0
	iteratorTypeSliceIndex = 1
	cursorSliceIndex       = 2
	createdAtSliceIndex    = 3
	indexSliceIndex        = 4
)

// A Position represents a Stripe position.
type Position struct {
	ResourceName string
	IteratorType string
	Cursor       string
	CreatedAt    int64
	Index        int
}

// ParseSDKPosition parses SDK position and returns Position.
func ParseSDKPosition(p sdk.Position, cfg *config.Config) (Position, error) {
	if p == nil {
		return Position{
			ResourceName: cfg.ResourceName,
			IteratorType: SnapshotType,
			CreatedAt:    time.Now().Unix(),
		}, nil
	}

	parts := strings.Split(string(p), ".")

	if len(parts) != reflect.TypeOf(Position{}).NumField() {
		return Position{}, fmt.Errorf("the number of position elements must be equal to %d, now it is %d",
			reflect.TypeOf(Position{}).NumField(), len(parts))
	}

	if parts[resourceSliceIndex] != cfg.ResourceName {
		return Position{
			ResourceName: cfg.ResourceName,
			IteratorType: SnapshotType,
			CreatedAt:    time.Now().Unix(),
		}, nil
	}

	createdAt, err := strconv.ParseInt(parts[createdAtSliceIndex], 10, 64)
	if err != nil {
		return Position{}, errParseCreatedAt
	}

	index, err := strconv.Atoi(parts[indexSliceIndex])
	if err != nil {
		return Position{}, errParseIndex
	}

	return Position{
		ResourceName: parts[resourceSliceIndex],
		IteratorType: parts[iteratorTypeSliceIndex],
		Cursor:       parts[cursorSliceIndex],
		CreatedAt:    createdAt,
		Index:        index,
	}, nil
}

// FormatSDKPosition formats and returns sdk.Position.
func (p Position) FormatSDKPosition() sdk.Position {
	return sdk.Position(fmt.Sprintf(positionFormat, p.ResourceName, p.IteratorType, p.Cursor, p.CreatedAt, p.Index))
}
