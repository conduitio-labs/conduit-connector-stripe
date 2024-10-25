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
	"time"

	"github.com/conduitio-labs/conduit-connector-stripe/models"
	"github.com/conduitio-labs/conduit-connector-stripe/models/resources"
	"github.com/conduitio-labs/conduit-connector-stripe/source/iterator/mock"
	"github.com/conduitio/conduit-commons/opencdc"
	"github.com/golang/mock/gomock"
)

const (
	cursor = "some_id"
)

func TestCDCIterator_Next(t *testing.T) {
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
						models.KeyID:      "price_1651153850",
						models.KeyObject:  "plan",
						models.KeyAmount:  1499,
						models.KeyCreated: float64(1651153850),
					}},
					Type: resources.PlanDeletedEvent,
				},
				// 3rd event, update amount of a plan from 12.99 to 14.99
				{
					ID:      "evt_1652447186",
					Created: 1652447186,
					Data: models.EventDataObject{Object: map[string]interface{}{
						models.KeyID:      "price_1651153850",
						models.KeyObject:  "plan",
						models.KeyAmount:  1499,
						models.KeyCreated: float64(1651153850),
					}},
					Type: resources.PlanUpdatedEvent,
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
						models.KeyID:      "price_1651153850",
						models.KeyObject:  "plan",
						models.KeyAmount:  1299,
						models.KeyCreated: float64(1651153850),
					}},
					Type: resources.PlanUpdatedEvent,
				},
				// 1st event, create a new plan
				{
					ID:      "evt_1652447136",
					Created: 1652447136,
					Data: models.EventDataObject{Object: map[string]interface{}{
						models.KeyID:      "price_1651153850",
						models.KeyObject:  "plan",
						models.KeyAmount:  1099,
						models.KeyCreated: float64(1651153850),
					}},
					Type: resources.PlanCreatedEvent,
				},
			},
			HasMore: false,
		}

		result.Data = append(result.Data, responseFirst.Data...)
		result.Data = append(result.Data, responseSecond.Data...)

		pos := &Position{
			IteratorMode: modeCDC,
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

			rp, err := pos.marshalPosition()
			if err != nil {
				t.Errorf("format sdk position error = \"%s\"", err.Error())
			}

			err = compareResult(record, rp, result.Data[i])
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

		responseFirst := models.EventResponse{
			Data: []models.EventData{
				// 2nd event, update amount of a plan from 10.99 to 12.99
				{
					ID:      "evt_1652447179",
					Created: 1652447179,
					Data: models.EventDataObject{Object: map[string]interface{}{
						models.KeyID:      "price_1651153850",
						models.KeyObject:  "plan",
						models.KeyAmount:  1299,
						models.KeyCreated: float64(1651153850),
					}},
					Type: resources.PlanUpdatedEvent,
				},
				// 1st event, create a new plan
				{
					ID:      "evt_1652447136",
					Created: 1652447136,
					Data: models.EventDataObject{Object: map[string]interface{}{
						models.KeyID:      "price_1651153850",
						models.KeyObject:  "plan",
						models.KeyAmount:  1099,
						models.KeyCreated: float64(1651153850),
					}},
					Type: resources.PlanCreatedEvent,
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
						models.KeyID:      "price_1651153850",
						models.KeyObject:  "plan",
						models.KeyAmount:  1499,
						models.KeyCreated: float64(1651153850),
					}},
					Type: resources.PlanDeletedEvent,
				},
				// 3rd event, update amount of a plan from 12.99 to 14.99
				{
					ID:      "evt_1652447186",
					Created: 1652447186,
					Data: models.EventDataObject{Object: map[string]interface{}{
						models.KeyID:      "price_1651153850",
						models.KeyObject:  "plan",
						models.KeyAmount:  1499,
						models.KeyCreated: float64(1651153850),
					}},
					Type: resources.PlanUpdatedEvent,
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

		pos := &Position{
			IteratorMode: modeCDC,
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

			rp, err := pos.marshalPosition()
			if err != nil {
				t.Errorf("format sdk position error = \"%s\"", err.Error())
			}

			err = compareResult(record, rp, result.Data[i])
			if err != nil {
				t.Error(err)
			}
		}
	})
}

func compareResult(record opencdc.Record, position opencdc.Position, data models.EventData) error {
	if !reflect.DeepEqual(record.Key, opencdc.StructuredData{models.KeyID: data.Data.Object[models.KeyID]}) {
		return fmt.Errorf("key: got = %v, want %v", string(record.Key.Bytes()), data.Data.Object[models.KeyID])
	}

	createdAt, err := record.Metadata.GetCreatedAt()
	if err != nil {
		return fmt.Errorf("get created_at error = \"%s\"", err.Error())
	}

	if createdAt != time.Unix(data.Created, 0) {
		return fmt.Errorf("created at: got = %v, want %v", createdAt, time.Unix(data.Created, 0))
	}

	operation := models.EventsOperation[data.Type]
	if record.Operation != operation {
		return fmt.Errorf("operation: got = %v, want %v", record.Operation, operation)
	}

	if !reflect.DeepEqual(record.Position, position) {
		return fmt.Errorf("position: got = %v, want %v", string(record.Position), string(position))
	}

	if operation == opencdc.OperationDelete {
		return nil
	}

	payload, err := json.Marshal(data.Data.Object)
	if err != nil {
		return fmt.Errorf("marshal payload error = \"%s\"", err.Error())
	}

	if !reflect.DeepEqual(record.Payload.After.Bytes(), payload) {
		return fmt.Errorf("payload: got = %v, want %v", string(record.Payload.After.Bytes()), string(payload))
	}

	return nil
}
