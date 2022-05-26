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
	"bytes"
	"reflect"
	"testing"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/conduitio/conduit-connector-stripe/models"
	"github.com/conduitio/conduit-connector-stripe/validator"
	"go.uber.org/multierr"
)

func TestParseSDKPosition(t *testing.T) {
	tests := []struct {
		name        string
		in          sdk.Position
		want        Position
		wantErr     bool
		expectedErr string
	}{
		{
			name: "valid sdk position of snapshot iterator",
			in: sdk.Position(`
{
	"iterator_type":"snapshot",
	"created_at":1652279623
}
`),
			want: Position{
				IteratorType: models.SnapshotIterator,
				CreatedAt:    1652279623,
				Cursor:       "",
				Index:        0,
			},
		},
		{
			name: "valid sdk position of snapshot iterator with cursor",
			in: sdk.Position(`
{
	"iterator_type":"snapshot",
	"created_at":1652279623,
	"cursor":"sub_1KtXkmJit567F2YtZzGSIrsh"
}
`),
			want: Position{
				IteratorType: models.SnapshotIterator,
				CreatedAt:    1652279623,
				Cursor:       "sub_1KtXkmJit567F2YtZzGSIrsh",
				Index:        0,
			},
		},
		{
			name: "valid sdk position of cdc iterator",
			in: sdk.Position(`
{
	"iterator_type":"cdc",
	"created_at":1652279623,
	"cursor":"evt_1KtXkmJit567F2YtZzGSIrsh",
	"index": 1
}
`),
			want: Position{
				IteratorType: models.CDCIterator,
				CreatedAt:    1652279623,
				Cursor:       "evt_1KtXkmJit567F2YtZzGSIrsh",
				Index:        1,
			},
		},
		{
			name: "invalid input data",
			in: sdk.Position(`
{
	"test":123
}
`),
			wantErr: true,
			expectedErr: multierr.Combine(validator.RequiredErr("IteratorType"),
				validator.RequiredErr("CreatedAt")).Error(),
		},
		{
			name: "IteratorType is required",
			in: sdk.Position(`
{
	"created_at":1652279623
}
`),
			wantErr:     true,
			expectedErr: validator.RequiredErr("IteratorType").Error(),
		},
		{
			name: "CreatedAt is required",
			in: sdk.Position(`
{
	"iterator_type":"snapshot"
}
`),
			wantErr:     true,
			expectedErr: validator.RequiredErr("CreatedAt").Error(),
		},
		{
			name: "unexpected iterator type",
			in: sdk.Position(`
{
	"iterator_type":"test",
	"created_at":1652279623
}
`),
			wantErr:     true,
			expectedErr: validator.UnexpectedIteratorTypeErr().Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSDKPosition(tt.in)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("parse error = \"%s\", wantErr %t", err.Error(), tt.wantErr)

					return
				}

				if err.Error() != tt.expectedErr {
					t.Errorf("expected error \"%s\", got \"%s\"", tt.expectedErr, err.Error())

					return
				}

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatSDKPosition(t *testing.T) {
	tests := []struct {
		name        string
		in          Position
		want        sdk.Position
		wantErr     bool
		expectedErr string
	}{
		{
			name: "valid position of snapshot iterator",
			in: Position{
				IteratorType: models.SnapshotIterator,
				CreatedAt:    1652279623,
			},
			want: sdk.Position(`{"iterator_type":"snapshot","created_at":1652279623,"cursor":"","index":0}`),
		},
		{
			name: "valid position of snapshot iterator with cursor",
			in: Position{
				IteratorType: models.SnapshotIterator,
				CreatedAt:    1652279623,
				Cursor:       "sub_1KtXkmJit567F2YtZzGSIrsh",
			},
			want: sdk.Position(
				`{"iterator_type":"snapshot","created_at":1652279623,"cursor":"sub_1KtXkmJit567F2YtZzGSIrsh","index":0}`),
		},
		{
			name: "valid sdk position of cdc iterator",
			in: Position{
				IteratorType: models.CDCIterator,
				CreatedAt:    1652279623,
				Cursor:       "evt_1KtXkmJit567F2YtZzGSIrsh",
				Index:        1,
			},
			want: sdk.Position(
				`{"iterator_type":"cdc","created_at":1652279623,"cursor":"evt_1KtXkmJit567F2YtZzGSIrsh","index":1}`),
		},
		{
			name: "IteratorType is required",
			in: Position{
				CreatedAt: 1652279623,
			},
			wantErr:     true,
			expectedErr: validator.RequiredErr("IteratorType").Error(),
		},
		{
			name: "CreatedAt is required",
			in: Position{
				IteratorType: models.SnapshotIterator,
			},
			wantErr:     true,
			expectedErr: validator.RequiredErr("CreatedAt").Error(),
		},
		{
			name: "unexpected iterator type",
			in: Position{
				IteratorType: "test",
				CreatedAt:    1652279623,
			},
			wantErr:     true,
			expectedErr: validator.UnexpectedIteratorTypeErr().Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.in.FormatSDKPosition()
			if err != nil {
				if !tt.wantErr {
					t.Errorf("parse error = \"%s\", wantErr %t", err.Error(), tt.wantErr)

					return
				}

				if err.Error() != tt.expectedErr {
					t.Errorf("expected error \"%s\", got \"%s\"", tt.expectedErr, err.Error())

					return
				}

				return
			}

			if !bytes.Equal(got, tt.want) {
				t.Errorf("parse = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
