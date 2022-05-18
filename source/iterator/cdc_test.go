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
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	sdk "github.com/conduitio/conduit-connector-sdk"

	"github.com/conduitio/conduit-connector-stripe/models"
	"github.com/conduitio/conduit-connector-stripe/source/iterator/mock"
	"github.com/conduitio/conduit-connector-stripe/source/position"
)

func TestCDC_Next(t *testing.T) {
	t.Run("success starting_after case", func(t *testing.T) {
		var (
			ctrl = gomock.NewController(t)
			m    = mock.NewMockStripe(ctrl)

			result = models.EventResponse{}
		)

		responseFirst := models.EventResponse{
			Data: []models.EventData{
				// 4th event, delete a plan
				{
					ID:      "evt_1652447199",
					Created: 1652447199,
					Data: models.EventDataObject{Object: map[string]interface{}{
						"id":      "price_1651153850",
						"object":  "plan",
						"amount":  1499,
						"created": float64(1651153850),
					}},
					Type: "plan.deleted",
				},
				// 3rd event, update amount of a plan from 12.99 to 14.99
				{
					ID:      "evt_1652447186",
					Created: 1652447186,
					Data: models.EventDataObject{Object: map[string]interface{}{
						"id":      "price_1651153850",
						"object":  "plan",
						"amount":  1499,
						"created": float64(1651153850),
					}},
					Type: "plan.updated",
				},
			},
			HasMore: true,
		}

		responseSecond := models.EventResponse{
			Data: []models.EventData{
				// 2nd event, update amount of a plan from 10.99 to 12.99
				{
					ID:      "evt_1652447179",
					Created: 1652447179,
					Data: models.EventDataObject{Object: map[string]interface{}{
						"id":      "price_1651153850",
						"object":  "plan",
						"amount":  1299,
						"created": float64(1651153850),
					}},
					Type: "plan.updated",
				},
				// 1st event, create a new plan
				{
					ID:      "evt_1652447136",
					Created: 1652447136,
					Data: models.EventDataObject{Object: map[string]interface{}{
						"id":      "price_1651153850",
						"object":  "plan",
						"amount":  1099,
						"created": float64(1651153850),
					}},
					Type: "plan.created",
				},
			},
			HasMore: false,
		}

		result.Data = append(result.Data, responseFirst.Data...)
		result.Data = append(result.Data, responseSecond.Data...)

		pos := &position.Position{
			IteratorType: position.CDCType,
			CreatedAt:    1652790765,
		}

		m.EXPECT().GetEvent(pos.CreatedAt, "", "").Return(responseFirst, nil)
		m.EXPECT().GetEvent(pos.CreatedAt, responseFirst.Data[len(responseFirst.Data)-1].ID, "").Return(responseSecond, nil)

		iter := NewCDC(m, pos)

		// reverse loop due to starting_after case
		for i := len(result.Data) - 1; i >= 0; i-- {
			record, err := iter.Next()
			if err != nil {
				t.Errorf("next error = \"%s\"", err.Error())
			}

			err = compareResult(record, pos.FormatSDKPosition(), result.Data[i])
			if err != nil {
				t.Error(err)
			}
		}
	})

	t.Run("success ending_before case", func(t *testing.T) {
		var (
			ctrl = gomock.NewController(t)
			m    = mock.NewMockStripe(ctrl)

			result = models.EventResponse{}
		)

		const cursor = "some_id"

		responseFirst := models.EventResponse{
			Data: []models.EventData{
				// 2nd event, update amount of a plan from 10.99 to 12.99
				{
					ID:      "evt_1652447179",
					Created: 1652447179,
					Data: models.EventDataObject{Object: map[string]interface{}{
						"id":      "price_1651153850",
						"object":  "plan",
						"amount":  1299,
						"created": float64(1651153850),
					}},
					Type: "plan.updated",
				},
				// 1st event, create a new plan
				{
					ID:      "evt_1652447136",
					Created: 1652447136,
					Data: models.EventDataObject{Object: map[string]interface{}{
						"id":      "price_1651153850",
						"object":  "plan",
						"amount":  1099,
						"created": float64(1651153850),
					}},
					Type: "plan.created",
				},
			},
			HasMore: true,
		}

		responseSecond := models.EventResponse{
			Data: []models.EventData{
				// 4th event, delete a plan
				{
					ID:      "evt_1652447199",
					Created: 1652447199,
					Data: models.EventDataObject{Object: map[string]interface{}{
						"id":      "price_1651153850",
						"object":  "plan",
						"amount":  1499,
						"created": float64(1651153850),
					}},
					Type: "plan.deleted",
				},
				// 3rd event, update amount of a plan from 12.99 to 14.99
				{
					ID:      "evt_1652447186",
					Created: 1652447186,
					Data: models.EventDataObject{Object: map[string]interface{}{
						"id":      "price_1651153850",
						"object":  "plan",
						"amount":  1499,
						"created": float64(1651153850),
					}},
					Type: "plan.updated",
				},
			},
			HasMore: false,
		}

		// reverse result before adding due to ending_before case
		for i := len(responseFirst.Data) - 1; i >= 0; i-- {
			result.Data = append(result.Data, responseFirst.Data[i])
		}

		for i := len(responseSecond.Data) - 1; i >= 0; i-- {
			result.Data = append(result.Data, responseSecond.Data[i])
		}

		pos := &position.Position{
			IteratorType: position.CDCType,
			Cursor:       cursor,
			CreatedAt:    1652790765,
		}

		m.EXPECT().GetEvent(pos.CreatedAt, "", cursor).Return(responseFirst, nil)
		m.EXPECT().GetEvent(pos.CreatedAt, "", responseFirst.Data[0].ID).Return(responseSecond, nil)

		iter := NewCDC(m, pos)

		for i := range result.Data {
			record, err := iter.Next()
			if err != nil {
				t.Errorf("next error = \"%s\"", err.Error())
			}

			err = compareResult(record, pos.FormatSDKPosition(), result.Data[i])
			if err != nil {
				t.Error(err)
			}
		}
	})
}

func compareResult(record sdk.Record, position sdk.Position, data models.EventData) error {
	payload, err := json.Marshal(data.Data.Object)
	if err != nil {
		return fmt.Errorf("marshal payload error = \"%s\"", err.Error())
	}

	if !reflect.DeepEqual(record.Payload.Bytes(), payload) {
		return fmt.Errorf("payload: got = %v, want %v", string(record.Payload.Bytes()), string(payload))
	}

	if !reflect.DeepEqual(record.Key, sdk.StructuredData{idKey: data.Data.Object["id"]}) {
		return fmt.Errorf("key: got = %v, want %v", string(record.Key.Bytes()), data.Data.Object["id"])
	}

	if record.CreatedAt.Unix() != data.Created {
		return fmt.Errorf("created: got = %v, want %v", record.CreatedAt.Unix(), data.Created)
	}

	action := models.EventsAction[data.Type]
	if record.Metadata[models.ActionKey] != action {
		return fmt.Errorf("action: got = %v, want %v", record.Metadata[models.ActionKey], action)
	}

	if !reflect.DeepEqual(record.Position, position) {
		return fmt.Errorf("position: got = %v, want %v", string(record.Position), string(position))
	}

	return nil
}
