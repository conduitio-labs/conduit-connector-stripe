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

func TestCDC_Next(t *testing.T) {
	t.Run("starting_after case", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		result := models.EventResponse{
			Data: []models.EventData{
				{
					ID:      "evt_1KyyCy3dG",
					Created: 1652447136,
					Data: models.EventDataObject{Object: map[string]interface{}{
						"id":      "cus_LY6gsj",
						"object":  "customer",
						"created": float64(1651153903),
					}},
					Type: "customer.created",
				},
				{
					ID:      "evt_Jit566F2Y",
					Created: 1652447179,
					Data: models.EventDataObject{Object: map[string]interface{}{
						"id":      "prod_LajCFS",
						"object":  "product",
						"created": float64(1651153850),
					}},
					Type: "product.deleted",
				},
			},
			HasMore: false,
		}

		pos := position.Position{
			IteratorType: position.SnapshotType,
			CreatedAt:    1652790765,
		}

		m := mock.NewMockStripe(ctrl)
		m.EXPECT().GetEvent(pos.CreatedAt, "", "").Return(result, nil)

		iter := NewCDC(m, &pos)

		// reverse loop due to starting_after case
		for i := len(result.Data) - 1; i >= 0; i-- {
			record, err := iter.Next()
			if err != nil {
				t.Errorf("next error = \"%s\"", err.Error())
			}

			payload, err := json.Marshal(result.Data[i].Data.Object)
			if err != nil {
				t.Errorf("marshal payload error = \"%s\"", err.Error())
			}

			if !reflect.DeepEqual(record.Payload.Bytes(), payload) {
				t.Errorf("payload: got = %v, want %v", string(record.Payload.Bytes()), string(payload))
			}

			if !reflect.DeepEqual(record.Key, sdk.StructuredData{idKey: result.Data[i].Data.Object["id"]}) {
				t.Errorf("key: got = %v, want %v", string(record.Key.Bytes()), result.Data[i].Data.Object["id"])
			}

			if record.CreatedAt.Unix() != result.Data[i].Created {
				t.Errorf("created: got = %v, want %v", record.CreatedAt.Unix(), result.Data[i].Created)
			}

			action := models.EventsAction[result.Data[i].Type]
			if record.Metadata[models.ActionKey] != action {
				t.Errorf("action: got = %v, want %v", record.Metadata[models.ActionKey], action)
			}

			if !reflect.DeepEqual(record.Position, pos.FormatSDKPosition()) {
				t.Errorf("position: got = %v, want %v", string(record.Position), string(pos.FormatSDKPosition()))
			}
		}
	})
}
