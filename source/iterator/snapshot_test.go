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
	"reflect"
	"testing"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/golang/mock/gomock"

	"github.com/conduitio/conduit-connector-stripe/models"
	"github.com/conduitio/conduit-connector-stripe/source/iterator/mock"
	"github.com/conduitio/conduit-connector-stripe/source/position"
)

func TestIterator_Next(t *testing.T) {
	pos := position.Position{
		IteratorType: position.SnapshotType,
	}
	underTestSnapshot := Snapshot{
		position: &pos,
	}

	type wantData struct {
		position string
		action   string
		key      sdk.StructuredData
		payload  string
	}

	tests := []struct {
		name        string
		in          []byte
		len         int
		want        []wantData
		wantErr     bool
		expectedErr string
	}{
		{
			name: "valid data",
			in: []byte(`{
    "data": [
        {
            "id": "prod_La50",
            "created": 1651153850
        },
		{
            "id": "prod_La49",
            "created": 1651153849
        },
		{
            "id": "prod_La48",
            "created": 1651153848
        }
    ],
    "has_more": false
}`),
			len: 3,
			want: []wantData{
				{
					position: "false..s.0",
					action:   "insert",
					key: sdk.StructuredData{
						idKey: "prod_La50",
					},
					payload: `{"created":1651153850,"id":"prod_La50"}`,
				},
				{
					position: "false.prod_La50.s.0",
					action:   "insert",
					key: sdk.StructuredData{
						idKey: "prod_La49",
					},
					payload: `{"created":1651153849,"id":"prod_La49"}`,
				},
				{
					position: "false.prod_La49.s.0",
					action:   "insert",
					key: sdk.StructuredData{
						idKey: "prod_La48",
					},
					payload: `{"created":1651153848,"id":"prod_La48"}`,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			underTestSnapshot.index = 0

			if err := json.Unmarshal(tt.in, &underTestSnapshot.response); err != nil {
				t.Fatalf("%s: %v", tt.in, err)
			}

			for i := 0; i < tt.len; i++ {
				got, err := underTestSnapshot.Next()
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

				if string(got.Position) != tt.want[i].position {
					t.Errorf("position: got \"%s\", want \"%s\"", got.Position, tt.want[i].position)

					return
				}

				if got.Metadata["action"] != tt.want[i].action {
					t.Errorf("action: got \"%s\", want \"%s\"", got.Metadata["action"], tt.want[i].action)

					return
				}

				if !reflect.DeepEqual(got.Key, tt.want[i].key) {
					t.Errorf("key: got \"%s\", want \"%s\"", got.Key.Bytes(), tt.want[i].key)

					return
				}

				if string(got.Payload.Bytes()) != tt.want[i].payload {
					t.Errorf("payload: got \"%s\", want \"%s\"", got.Payload.Bytes(), tt.want[i].payload)

					return
				}
			}
		})
	}
}

func TestIterator_Integration_Next(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		result := models.ResourceResponse{
			Data: []map[string]interface{}{
				{
					"id":      "cus_LY6gsj",
					"object":  "customer",
					"created": float64(1651153903),
				},
				{
					"id":      "prod_LajCFS",
					"object":  "product",
					"created": float64(1651153850),
				},
			},
			HasMore: false,
		}

		pos := position.Position{
			IteratorType: position.SnapshotType,
		}

		m := mock.NewMockStripe(ctrl)
		m.EXPECT().GetResource(pos.Cursor).Return(result, nil)

		iter := NewSnapshot(m, &pos)

		for i := 0; i < len(result.Data); i++ {
			positionForCheck := pos.FormatSDKPosition()

			record, err := iter.Next()
			if err != nil {
				t.Errorf("next error = \"%s\"", err.Error())
			}

			payload, err := json.Marshal(result.Data[i])
			if err != nil {
				t.Errorf("marshal payload error = \"%s\"", err.Error())
			}

			if !reflect.DeepEqual(record.Payload.Bytes(), payload) {
				t.Errorf("payload: got = %v, want %v", string(record.Payload.Bytes()), string(payload))
			}

			if !reflect.DeepEqual(record.Key, sdk.StructuredData{idKey: result.Data[i]["id"]}) {
				t.Errorf("key: got = %v, want %v", string(record.Key.Bytes()), result.Data[i]["id"])
			}

			if record.CreatedAt.Unix() != int64(result.Data[i]["created"].(float64)) {
				t.Errorf("created: got = %v, want %v", record.CreatedAt.Unix(), result.Data[i]["created"])
			}

			if record.Metadata[models.ActionKey] != models.InsertAction {
				t.Errorf("action: got = %v, want %v", record.Metadata[models.ActionKey], models.InsertAction)
			}

			if !reflect.DeepEqual(record.Position, positionForCheck) {
				t.Errorf("position: got = %v, want %v", string(record.Position), string(pos.FormatSDKPosition()))
			}
		}
	})
}
