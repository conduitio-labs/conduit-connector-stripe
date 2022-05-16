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

import (
	"strings"

	"github.com/conduitio/conduit-connector-stripe/models/resources"
)

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
	resources.CreditNoteResource:                  resources.CreditNotesList,
	resources.BillingPortalConfigurationResource:  resources.BillingPortalConfigurationsList,
	resources.InvoiceResource:                     resources.InvoicesList,
	resources.InvoiceItemResource:                 resources.InvoiceItemsList,
	resources.PlanResource:                        resources.PlanList,
	resources.QuoteResource:                       resources.QuotesList,
	resources.SubscriptionResource:                resources.SubscriptionsList,
	resources.SubscriptionItemResource:            resources.SubscriptionItemsList,
	resources.SubscriptionScheduleResource:        resources.SubscriptionSchedulesList,
	resources.CheckoutSessionResource:             resources.CheckoutSessionsList,
	resources.AccountResource:                     resources.AccountsList,
	resources.ApplicationFeeResource:              resources.ApplicationFeesList,
	resources.CountrySpecResource:                 resources.CountrySpecsList,
	resources.TopUpResource:                       resources.TopUpsList,
	resources.TransferResource:                    resources.TransfersList,
	resources.BalanceTransactionResource:          resources.BalanceTransactionsList,
	resources.ChargeResource:                      resources.ChargesList,
	resources.CustomerResource:                    resources.CustomersList,
	resources.DisputeResource:                     resources.DisputesList,
	resources.EventResource:                       resources.EventsList,
	resources.FileResource:                        resources.FilesList,
	resources.FileLinkResource:                    resources.FileLinksList,
	resources.PaymentIntentResource:               resources.PaymentIntentsList,
	resources.SetupIntentResource:                 resources.SetupIntentsList,
	resources.SetupAttemptResource:                resources.SetupAttemptsList,
	resources.PayoutResource:                      resources.PayoutsList,
	resources.RefundResource:                      resources.RefundsList,
	resources.FinancialConnectionsAccountResource: resources.FinancialConnectionsAccountsResource,
	resources.RadarEarlyFraudWarningResource:      resources.RadarEarlyFraudWarningsList,
	resources.ReviewResource:                      resources.ReviewsList,
	resources.RadarValueListResource:              resources.RadarValueListsList,
	resources.RadarValueListItemResource:          resources.RadarValueListItemsList,
	resources.IdentityVerificationSessionResource: resources.IdentityVerificationSessionsList,
	resources.IdentityVerificationReportResource:  resources.IdentityVerificationReportsList,
	resources.IssuingAuthorizationResource:        resources.IssuingAuthorizationsList,
	resources.IssuingCardholderResource:           resources.IssuingCardholdersList,
	resources.IssuingCardResource:                 resources.IssuingCardsList,
	resources.IssuingDisputeResource:              resources.IssuingDisputesList,
	resources.FundingInstructionResource:          resources.FundingInstructionsList,
	resources.IssuingTransactionResource:          resources.IssuingTransactionsList,
	resources.OrderResource:                       resources.OrdersList,
	resources.PaymentLinkResource:                 resources.PaymentLinksList,
	resources.PaymentMethodResource:               resources.PaymentMethodsList,
	resources.ProductResource:                     resources.ProductsList,
	resources.PriceResource:                       resources.PricesList,
	resources.CouponResource:                      resources.CouponsList,
	resources.PromotionCode:                       resources.PromotionCodesList,
	resources.TaxCodeResource:                     resources.TaxCodesList,
	resources.TaxRateResource:                     resources.TaxRatesList,
	resources.ShippingRateResource:                resources.ShippingRatesList,
	resources.ReportingReportRunResource:          resources.ReportingReportRunsList,
	resources.ReportingReportTypeResource:         resources.ReportingReportTypesList,
	resources.ScheduledQueryRunResource:           resources.ScheduledQueryRunsList,
	resources.TerminalLocationResource:            resources.TerminalLocationsList,
	resources.TerminalReaderResource:              resources.TerminalReadersList,
	resources.TerminalConfigurationResource:       resources.TerminalConfigurationsList,
	resources.FinancialAccountResource:            resources.FinancialAccountsList,
	resources.TransactionResource:                 resources.TransactionsList,
	resources.TransactionEntryResource:            resources.TransactionEntriesList,
	resources.OutboundTransferResource:            resources.OutboundTransfersList,
	resources.OutboundPaymentResource:             resources.OutboundPaymentsList,
	resources.InboundTransferResource:             resources.InboundTransfersList,
	resources.ReceivedCreditResource:              resources.ReceivedCreditsList,
	resources.ReceivedDebitResource:               resources.ReceivedDebitsList,
	resources.CreditReversalResource:              resources.CreditReversalsList,
	resources.DebitReversalResource:               resources.DebitReversalsList,
	resources.WebhookEndpointResource:             resources.WebhookEndpointsList,
}

// EventsMap represents a dictionary with all events in each resource,
// where the key is the resource and the value is a slice of events.
var EventsMap = map[string][]string{
	resources.CreditNoteResource: {
		resources.CreditNoteCreatedEvent,
		resources.CreditNoteUpdatedEvent,
		resources.CreditNoteVoidedEvent,
	},
	resources.InvoiceResource: {
		resources.InvoiceCreatedEvent,
		resources.InvoiceDeletedEvent,
		resources.InvoiceFinalizationFailedEvent,
		resources.InvoiceFinalizedEvent,
		resources.InvoiceMarkedUncollectibleEvent,
		resources.InvoicePaidEvent,
		resources.InvoicePaymentActionRequiredEvent,
		resources.InvoicePaymentFailedEvent,
		resources.InvoicePaymentSucceededEvent,
		resources.InvoiceSentEvent,
		resources.InvoiceUpcomingEvent,
		resources.InvoiceUpdatedEvent,
		resources.InvoiceVoidedEvent,
	},
	resources.InvoiceItemResource: {
		resources.InvoiceItemCreatedEvent,
		resources.InvoiceItemDeletedEvent,
		resources.InvoiceItemUpdatedEvent,
	},
	resources.PlanResource: {
		resources.PlanCreatedEvent,
		resources.PlanDeletedEvent,
		resources.PlanUpdatedEvent,
	},
	resources.QuoteResource: {
		resources.QuoteAcceptedEvent,
		resources.QuoteCanceledEvent,
		resources.QuoteCreatedEvent,
		resources.QuoteFinalizedEvent,
	},
	resources.SubscriptionResource: {
		resources.SubscriptionCreatedEvent,
		resources.SubscriptionDeletedEvent,
		resources.SubscriptionPendingUpdateAppliedEvent,
		resources.SubscriptionPendingUpdateExpiredEvent,
		resources.SubscriptionTrialWillEndEvent,
		resources.SubscriptionUpdatedEvent,
	},
	resources.SubscriptionScheduleResource: {
		resources.SubscriptionScheduleAbortedEvent,
		resources.SubscriptionScheduleCanceledEvent,
		resources.SubscriptionScheduleCompletedEvent,
		resources.SubscriptionScheduleCreatedEvent,
		resources.SubscriptionScheduleExpiringEvent,
		resources.SubscriptionScheduleReleasedEvent,
		resources.SubscriptionScheduleUpdatedEvent,
	},
	resources.CheckoutSessionResource: {
		resources.CheckoutSessionAsyncPaymentFailedEvent,
		resources.CheckoutSessionAsyncPaymentSucceededEvent,
		resources.CheckoutSessionCompletedEvent,
		resources.CheckoutSessionExpiredEvent,
	},
	resources.AccountResource: {
		resources.AccountUpdatedEvent,
	},
	resources.ApplicationFeeResource: {
		resources.ApplicationFeeCreatedEvent,
		resources.ApplicationFeeRefundedEvent,
	},
	resources.TopUpResource: {
		resources.TopupCanceledEvent,
		resources.TopupCreatedEvent,
		resources.TopupFailedEvent,
		resources.TopupReversedEvent,
		resources.TopupSucceededEvent,
	},
	resources.TransferResource: {
		resources.TransferCreatedEvent,
		resources.TransferFailedEvent,
		resources.TransferPaidEvent,
		resources.TransferReversedEvent,
		resources.TransferUpdatedEvent,
	},
	resources.ChargeResource: {
		resources.ChargeCapturedEvent,
		resources.ChargeExpiredEvent,
		resources.ChargeFailedEvent,
		resources.ChargePendingEvent,
		resources.ChargeRefundedEvent,
		resources.ChargeSucceededEvent,
		resources.ChargeUpdatedEvent,
	},
	resources.CustomerResource: {
		resources.CustomerCreatedEvent,
		resources.CustomerDeletedEvent,
		resources.CustomerUpdatedEvent,
	},
	resources.DisputeResource: {
		resources.DisputeClosedEvent,
		resources.DisputeCreatedEvent,
		resources.DisputeFundsReinstatedEvent,
		resources.DisputeFundsWithdrawnEvent,
		resources.DisputeupdatedEvent,
	},
	resources.PaymentIntentResource: {
		resources.PaymentIntentAmountCapturableUpdatedEvent,
		resources.PaymentIntentCanceledEvent,
		resources.PaymentIntentCreatedEvent,
		resources.PaymentIntentPartiallyFundedEvent,
		resources.PaymentIntentPaymentFailedEvent,
		resources.PaymentIntentProcessingEvent,
		resources.PaymentIntentRequiresActionEvent,
		resources.PaymentIntentSucceededEvent,
	},
	resources.SetupIntentResource: {
		resources.SetupIntentCanceledEvent,
		resources.SetupIntentCreatedEvent,
		resources.SetupIntentRequiresActionEvent,
		resources.SetupIntentSetupFailedEvent,
		resources.SetupIntentSucceededEvent,
	},
	resources.PayoutResource: {
		resources.PayoutsCanceledEvent,
		resources.PayoutsCreatedEvent,
		resources.PayoutsFailedEvent,
		resources.PayoutsPaidEvent,
		resources.PayoutsUpdatedEvent,
	},
	resources.RefundResource: {
		resources.RefundUpdatedEvent,
	},
	resources.RadarEarlyFraudWarningResource: {
		resources.RadarEarlyFraudWarningCreatedEvent,
		resources.RadarEarlyFraudWarningUpdatedEvent,
	},
	resources.ReviewResource: {
		resources.ReviewClosedEvent,
		resources.ReviewOpenedEvent,
	},
	resources.IdentityVerificationSessionResource: {
		resources.IdentityVerificationSessionCanceledEvent,
		resources.IdentityVerificationSessionCreatedEvent,
		resources.IdentityVerificationSessionProcessingEvent,
		resources.IdentityVerificationSessionRedactedEvent,
		resources.IdentityVerificationSessionRequiresInputEvent,
		resources.IdentityVerificationSessionVerifiedEvent,
	},
	resources.IssuingAuthorizationResource: {
		resources.IssuingAuthorizationCreatedEvent,
		resources.IssuingAuthorizationRequestEvent,
		resources.IssuingAuthorizationUpdatedEvent,
	},
	resources.IssuingCardholderResource: {
		resources.IssuingCardholderCreatedEvent,
		resources.IssuingCardholderUpdatedEvent,
	},
	resources.IssuingCardResource: {
		resources.IssuingCardCreatedEvent,
		resources.IssuingCardUpdatedEvent,
	},
	resources.IssuingDisputeResource: {
		resources.IssuingDisputeClosedEvent,
		resources.IssuingDisputeCreatedEvent,
		resources.IssuingDisputeFundsReinstatedEvent,
		resources.IssuingDisputeSubmittedEvent,
		resources.IssuingDisputeUpdatedEvent,
	},
	resources.IssuingTransactionResource: {
		resources.IssuingTransactionCreatedEvent,
		resources.IssuingTransactionUpdatedEvent,
	},
	resources.OrderResource: {
		resources.OrderCreatedEvent,
		resources.OrderPaymentFailedEvent,
		resources.OrderPaymentSucceededEvent,
		resources.OrderUpdatedEvent,
	},
	resources.PaymentLinkResource: {
		resources.PaymentLinkCreatedEvent,
		resources.PaymentLinkUpdatedEvent,
	},
	resources.PaymentMethodResource: {
		resources.PaymentMethodAttachedEvent,
		resources.PaymentMethodAutomaticallyUpdatedEvent,
		resources.PaymentMethodDetachedEvent,
		resources.PaymentMethodUpdatedEvent,
	},
	resources.ProductResource: {
		resources.ProductCreatedEvent,
		resources.ProductDeletedEvent,
		resources.ProductUpdatedEvent,
	},
	resources.PriceResource: {
		resources.PriceCreatedEvent,
		resources.PriceDeletedEvent,
		resources.PriceUpdatedEvent,
	},
	resources.CouponResource: {
		resources.CouponCreatedEvent,
		resources.CouponDeletedEvent,
		resources.CouponUpdatedEvent,
	},
	resources.PromotionCode: {
		resources.PromotionCodeCreatedEvent,
		resources.PromotionCodeUpdatedEvent,
	},
	resources.TaxRateResource: {
		resources.TaxRateCreatedEvent,
		resources.TaxRateUpdatedEvent,
	},
	resources.ReportingReportRunResource: {
		resources.ReportingReportRunFailedEvent,
		resources.ReportingReportRunSucceededEvent,
	},
	resources.ReportingReportTypeResource: {
		resources.ReportingReportTypeUpdatedEvent,
	},
	resources.ScheduledQueryRunResource: {
		resources.ScheduledQueryRunCreatedEvent,
	},
	resources.TerminalReaderResource: {
		resources.TerminalReaderActionFailedEvent,
		resources.TerminalReaderActionSucceededEvent,
	},
}

// EventsAction represents a dictionary with actions of events,
// where the key is an event and the value is an action.
var EventsAction = (func() map[string]string {
	eventsAction := make(map[string]string)

	for _, events := range EventsMap {
		for _, event := range events {
			switch e := event; {
			case strings.Contains(e, eventKeyCreated):
				eventsAction[e] = InsertAction
			case strings.Contains(e, eventKeyDeleted):
				eventsAction[e] = DeleteAction
			default:
				eventsAction[e] = UpdateAction
			}
		}
	}

	return eventsAction
})()
