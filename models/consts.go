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

const (
	ActionKey = "action"

	InsertAction = "insert"
	UpdateAction = "update"
	DeleteAction = "delete"

	SubscriptionResource                  = "subscription"
	SubscriptionsList                     = "subscriptions"
	SubscriptionCreatedEvent              = "customer.subscription.created"
	SubscriptionDeletedEvent              = "customer.subscription.deleted"
	SubscriptionPendingUpdateAppliedEvent = "customer.subscription.pending_update_applied"
	SubscriptionPendingUpdateExpiredEvent = "customer.subscription.pending_update_expired"
	SubscriptionTrialWillEndEvent         = "customer.subscription.trial_will_end"
	SubscriptionUpdatedEvent              = "customer.subscription.updated"

	PlanResource     = "plan"
	PlanList         = "plans"
	PlanCreatedEvent = "plan.created"
	PlanDeletedEvent = "plan.deleted"
	PlanUpdatedEvent = "plan.updated"
)
