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
	CreditNoteResource     = "credit_note"
	CreditNotesList        = "credit_notes"
	CreditNoteCreatedEvent = "credit_note.created"
	CreditNoteUpdatedEvent = "credit_note.updated"
	CreditNoteVoidedEvent  = "credit_note.voided"

	BillingPortalConfigurationResource     = "billing_portal.configuration"
	BillingPortalConfigurationsList        = "billing_portal/configurations"
	BillingPortalConfigurationCreatedEvent = "billing_portal.configuration.created"
	BillingPortalConfigurationUpdatedEvent = "billing_portal.configuration.updated"

	InvoiceResource                   = "invoice"
	InvoicesList                      = "invoices"
	InvoiceCreatedEvent               = "invoice.created"
	InvoiceDeletedEvent               = "invoice.deleted"
	InvoiceFinalizationFailedEvent    = "invoice.finalization_failed"
	InvoiceFinalizedEvent             = "invoice.finalized"
	InvoiceMarkedUncollectibleEvent   = "invoice.marked_uncollectible"
	InvoicePaidEvent                  = "invoice.paid"
	InvoicePaymentActionRequiredEvent = "invoice.payment_action_required"
	InvoicePaymentFailedEvent         = "invoice.payment_failed"
	InvoicePaymentSucceededEvent      = "invoice.payment_succeeded"
	InvoiceSentEvent                  = "invoice.sent"
	InvoiceUpcomingEvent              = "invoice.upcoming"
	InvoiceUpdatedEvent               = "invoice.updated"
	InvoiceVoidedEvent                = "invoice.voided"

	InvoiceItemResource     = "invoiceitem"
	InvoiceItemsList        = "invoiceitems"
	InvoiceItemCreatedEvent = "invoiceitem.created"
	InvoiceItemDeletedEvent = "invoiceitem.deleted"
	InvoiceItemUpdatedEvent = "invoiceitem.updated"

	PlanResource     = "plan"
	PlanList         = "plans"
	PlanCreatedEvent = "plan.created"
	PlanDeletedEvent = "plan.deleted"
	PlanUpdatedEvent = "plan.updated"

	QuoteResource       = "quote"
	QuotesList          = "quotes"
	QuoteAcceptedEvent  = "quote.accepted"
	QuoteCanceledEvent  = "quote.canceled"
	QuoteCreatedEvent   = "quote.created"
	QuoteFinalizedEvent = "quote.finalized"

	SubscriptionResource                  = "subscription"
	SubscriptionsList                     = "subscriptions"
	SubscriptionCreatedEvent              = "customer.subscription.created"
	SubscriptionDeletedEvent              = "customer.subscription.deleted"
	SubscriptionPendingUpdateAppliedEvent = "customer.subscription.pending_update_applied"
	SubscriptionPendingUpdateExpiredEvent = "customer.subscription.pending_update_expired"
	SubscriptionTrialWillEndEvent         = "customer.subscription.trial_will_end"
	SubscriptionUpdatedEvent              = "customer.subscription.updated"

	SubscriptionItemResource = "subscription_item"
	SubscriptionItemsList    = "subscription_items"

	SubscriptionScheduleResource       = "subscription_schedule"
	SubscriptionSchedulesList          = "subscription_schedules"
	SubscriptionScheduleAbortedEvent   = "subscription_schedule.aborted"
	SubscriptionScheduleCanceledEvent  = "subscription_schedule.canceled"
	SubscriptionScheduleCompletedEvent = "subscription_schedule.completed"
	SubscriptionScheduleCreatedEvent   = "subscription_schedule.created"
	SubscriptionScheduleExpiringEvent  = "subscription_schedule.expiring"
	SubscriptionScheduleReleasedEvent  = "subscription_schedule.released"
	SubscriptionScheduleUpdatedEvent   = "subscription_schedule.updated"
)

var (
	CreditNoteEvents = []string{
		CreditNoteCreatedEvent,
		CreditNoteUpdatedEvent,
		CreditNoteVoidedEvent,
	}

	BillingPortalConfigurationEvents = []string{
		BillingPortalConfigurationCreatedEvent,
		BillingPortalConfigurationUpdatedEvent,
	}

	InvoiceEvents = []string{
		InvoiceCreatedEvent,
		InvoiceDeletedEvent,
		InvoiceFinalizationFailedEvent,
		InvoiceFinalizedEvent,
		InvoiceMarkedUncollectibleEvent,
		InvoicePaidEvent,
		InvoicePaymentActionRequiredEvent,
		InvoicePaymentFailedEvent,
		InvoicePaymentSucceededEvent,
		InvoiceSentEvent,
		InvoiceUpcomingEvent,
		InvoiceUpdatedEvent,
		InvoiceVoidedEvent,
	}

	InvoiceItemEvents = []string{
		InvoiceItemCreatedEvent,
		InvoiceItemDeletedEvent,
		InvoiceItemUpdatedEvent,
	}

	PlanEvents = []string{
		PlanCreatedEvent,
		PlanDeletedEvent,
		PlanUpdatedEvent,
	}

	QuoteEvents = []string{
		QuoteAcceptedEvent,
		QuoteCanceledEvent,
		QuoteCreatedEvent,
		QuoteFinalizedEvent,
	}

	SubscriptionEvents = []string{
		SubscriptionCreatedEvent,
		SubscriptionDeletedEvent,
		SubscriptionPendingUpdateAppliedEvent,
		SubscriptionPendingUpdateExpiredEvent,
		SubscriptionTrialWillEndEvent,
		SubscriptionUpdatedEvent,
	}

	SubscriptionScheduleEvents = []string{
		SubscriptionScheduleAbortedEvent,
		SubscriptionScheduleCanceledEvent,
		SubscriptionScheduleCompletedEvent,
		SubscriptionScheduleCreatedEvent,
		SubscriptionScheduleExpiringEvent,
		SubscriptionScheduleReleasedEvent,
		SubscriptionScheduleUpdatedEvent,
	}
)
