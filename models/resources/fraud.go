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

package resources

const (
	RadarEarlyFraudWarningResource     = "radar.early_fraud_warning"
	RadarEarlyFraudWarningsList        = "radar/early_fraud_warnings"
	RadarEarlyFraudWarningCreatedEvent = "radar.early_fraud_warning.created"
	RadarEarlyFraudWarningUpdatedEvent = "radar.early_fraud_warning.updated"

	ReviewResource    = "review"
	ReviewsList       = "reviews"
	ReviewClosedEvent = "review.closed"
	ReviewOpenedEvent = "review.opened"

	RadarValueListResource = "radar.value_list"
	RadarValueListsList    = "radar/value_lists"

	RadarValueListItemResource = "radar.value_list_item"
	RadarValueListItemsList    = "radar/value_list_items"
)

var (
	RadarEarlyFraudWarningEvents = []string{
		RadarEarlyFraudWarningCreatedEvent,
		RadarEarlyFraudWarningUpdatedEvent,
	}

	ReviewEvents = []string{
		ReviewClosedEvent,
		ReviewOpenedEvent,
	}
)
