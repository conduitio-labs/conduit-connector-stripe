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

	"github.com/conduitio-labs/conduit-connector-stripe/models"
	"github.com/conduitio-labs/conduit-connector-stripe/source/iterator/mock"
	"github.com/conduitio-labs/conduit-connector-stripe/source/position"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/golang/mock/gomock"
)

func TestSnapshotIterator_Next(t *testing.T) {
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
		}

		pos := position.Position{
			IteratorType: models.SnapshotIterator,
			CreatedAt:    1652790765,
		}

		m := mock.NewMockStripe(ctrl)
		m.EXPECT().GetResource(pos.Cursor).Return(result, nil)

		iter := NewSnapshotIterator(m, &pos)

		for i := 0; i < len(result.Data); i++ {
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
				t.Errorf("created: got = %v, want %v",
					record.CreatedAt.Unix(), int64(result.Data[i]["created"].(float64)))
			}

			if record.Metadata[models.ActionKey] != models.InsertAction {
				t.Errorf("action: got = %v, want %v", record.Metadata[models.ActionKey], models.InsertAction)
			}

			rp, err := pos.FormatSDKPosition()
			if err != nil {
				t.Errorf("format sdk position error = \"%s\"", err.Error())
			}

			if !reflect.DeepEqual(record.Position, rp) {
				t.Errorf("position: got = %v, want %v", string(record.Position), string(rp))
			}
		}
	})
}
