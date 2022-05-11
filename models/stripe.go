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

package models

// A ResourceResponse represents a response resource data from Stripe.
type ResourceResponse struct {
	Data    []map[string]interface{} `json:"data"`
	HasMore bool                     `json:"has_more"`
}

// A EventData represents a data of event response.
type EventData struct {
	ID      string `json:"id"`
	Created int64  `json:"created"`
	Data    struct {
		Object map[string]interface{} `json:"object"`
	} `json:"data"`
	Type string `json:"type"`
}

// A EventResponse represents a response event data from Stripe.
type EventResponse struct {
	Data    []EventData `json:"data"`
	HasMore bool        `json:"has_more"`
}

// ResourcesMap represents a dictionary with valid resources,
// where the key is the object type and the value is the name
// of the API endpoints of that object type.
var ResourcesMap = map[string]string{
	SubscriptionResource: SubscriptionsList,
	PlanResource:         PlanList,
}

// EventsMap represents a dictionary with all events in each resource,
// where the key is the resource and the value is a slice of events.
var EventsMap = map[string][]string{
	SubscriptionResource: {
		SubscriptionCreatedEvent,
		SubscriptionDeletedEvent,
		SubscriptionPendingUpdateAppliedEvent,
		SubscriptionPendingUpdateExpiredEvent,
		SubscriptionTrialWillEndEvent,
		SubscriptionUpdatedEvent,
	},
	PlanResource: {
		PlanCreatedEvent,
		PlanDeletedEvent,
		PlanUpdatedEvent,
	},
}

// EventsAction represents a dictionary with actions of events,
// where the key is an event and the value is an action.
var EventsAction = map[string]string{
	SubscriptionCreatedEvent:              InsertAction,
	SubscriptionDeletedEvent:              DeleteAction,
	SubscriptionPendingUpdateAppliedEvent: UpdateAction,
	SubscriptionPendingUpdateExpiredEvent: UpdateAction,
	SubscriptionTrialWillEndEvent:         UpdateAction,
	SubscriptionUpdatedEvent:              UpdateAction,
	PlanCreatedEvent:                      InsertAction,
	PlanDeletedEvent:                      DeleteAction,
	PlanUpdatedEvent:                      UpdateAction,
}
