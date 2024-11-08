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

	"github.com/conduitio-labs/conduit-connector-stripe/models/resources"
	"github.com/conduitio/conduit-commons/opencdc"
)

// A ResourceResponse represents a response resource data from Stripe.
type ResourceResponse struct {
	Data []map[string]interface{} `json:"data"`
}

// A EventResponse represents a response event data from Stripe.
type EventResponse struct {
	Data    EventsData `json:"data"`
	HasMore bool       `json:"has_more"`
}

// An EventsData represents a slice of EventData.
type EventsData []EventData

// A EventData represents a data of event response.
type EventData struct {
	ID      string          `json:"id"`
	Created int64           `json:"created"`
	Data    EventDataObject `json:"data"`
	Type    string          `json:"type"`
}

// An EventDataObject represents a full object of event data.
type EventDataObject struct {
	Object map[string]interface{} `json:"object"`
}

// An ErrorResponse represents a response error from Stripe.
type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
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
	resources.SubscriptionScheduleResource:        resources.SubscriptionSchedulesList,
	resources.CheckoutSessionResource:             resources.CheckoutSessionsList,
	resources.AccountResource:                     resources.AccountsList,
	resources.ApplicationFeeResource:              resources.ApplicationFeesList,
	resources.TopUpResource:                       resources.TopUpsList,
	resources.TransferResource:                    resources.TransfersList,
	resources.ChargeResource:                      resources.ChargesList,
	resources.CustomerResource:                    resources.CustomersList,
	resources.DisputeResource:                     resources.DisputesList,
	resources.FileResource:                        resources.FilesList,
	resources.PaymentIntentResource:               resources.PaymentIntentsList,
	resources.SetupIntentResource:                 resources.SetupIntentsList,
	resources.PayoutResource:                      resources.PayoutsList,
	resources.RefundResource:                      resources.RefundsList,
	resources.RadarEarlyFraudWarningResource:      resources.RadarEarlyFraudWarningsList,
	resources.ReviewResource:                      resources.ReviewsList,
	resources.IdentityVerificationSessionResource: resources.IdentityVerificationSessionsList,
	resources.IssuingAuthorizationResource:        resources.IssuingAuthorizationsList,
	resources.IssuingCardholderResource:           resources.IssuingCardholdersList,
	resources.IssuingCardResource:                 resources.IssuingCardsList,
	resources.IssuingDisputeResource:              resources.IssuingDisputesList,
	resources.IssuingTransactionResource:          resources.IssuingTransactionsList,
	resources.OrderResource:                       resources.OrdersList,
	resources.PaymentLinkResource:                 resources.PaymentLinksList,
	resources.PaymentMethodResource:               resources.PaymentMethodsList,
	resources.ProductResource:                     resources.ProductsList,
	resources.PriceResource:                       resources.PricesList,
	resources.CouponResource:                      resources.CouponsList,
	resources.PromotionCodeResource:               resources.PromotionCodesList,
	resources.TaxRateResource:                     resources.TaxRatesList,
	resources.ReportingReportRunResource:          resources.ReportingReportRunsList,
	resources.ReportingReportTypeResource:         resources.ReportingReportTypesList,
	resources.ScheduledQueryRunResource:           resources.ScheduledQueryRunsList,
	resources.TerminalReaderResource:              resources.TerminalReadersList,
}

// EventsMap represents a dictionary with all events in each resource,
// where the key is the resource and the value is a slice of events.
var EventsMap = map[string][]string{
	resources.CreditNoteResource:                  resources.CreditNoteEvents,
	resources.InvoiceResource:                     resources.InvoiceEvents,
	resources.InvoiceItemResource:                 resources.InvoiceItemEvents,
	resources.PlanResource:                        resources.PlanEvents,
	resources.QuoteResource:                       resources.QuoteEvents,
	resources.SubscriptionResource:                resources.SubscriptionEvents,
	resources.SubscriptionScheduleResource:        resources.SubscriptionScheduleEvents,
	resources.CheckoutSessionResource:             resources.CheckoutSessionEvents,
	resources.AccountResource:                     resources.AccountEvents,
	resources.ApplicationFeeResource:              resources.ApplicationFeeEvents,
	resources.TopUpResource:                       resources.TopUpEvents,
	resources.TransferResource:                    resources.TransferEvents,
	resources.ChargeResource:                      resources.ChargeEvents,
	resources.CustomerResource:                    resources.CustomerEvents,
	resources.DisputeResource:                     resources.DisputeEvents,
	resources.FileResource:                        resources.FileEvents,
	resources.PaymentIntentResource:               resources.PaymentIntentEvents,
	resources.SetupIntentResource:                 resources.SetupIntentEvents,
	resources.PayoutResource:                      resources.PayoutEvents,
	resources.RefundResource:                      resources.RefundEvents,
	resources.RadarEarlyFraudWarningResource:      resources.RadarEarlyFraudWarningEvents,
	resources.ReviewResource:                      resources.ReviewEvents,
	resources.IdentityVerificationSessionResource: resources.IdentityVerificationSessionEvents,
	resources.IssuingAuthorizationResource:        resources.IssuingAuthorizationEvents,
	resources.IssuingCardholderResource:           resources.IssuingCardholderEvents,
	resources.IssuingCardResource:                 resources.IssuingCardEvents,
	resources.IssuingDisputeResource:              resources.IssuingDisputeEvents,
	resources.IssuingTransactionResource:          resources.IssuingTransactionEvents,
	resources.OrderResource:                       resources.OrderEvents,
	resources.PaymentLinkResource:                 resources.PaymentLinkEvents,
	resources.PaymentMethodResource:               resources.PaymentMethodEvents,
	resources.ProductResource:                     resources.ProductEvents,
	resources.PriceResource:                       resources.PriceEvents,
	resources.CouponResource:                      resources.CouponEvents,
	resources.PromotionCodeResource:               resources.PromotionCodeEvents,
	resources.TaxRateResource:                     resources.TaxRateEvents,
	resources.ReportingReportRunResource:          resources.ReportingReportRunEvents,
	resources.ReportingReportTypeResource:         resources.ReportingReportTypeEvents,
	resources.ScheduledQueryRunResource:           resources.ScheduledQueryRunEvents,
	resources.TerminalReaderResource:              resources.TerminalReaderEvents,
}

// EventsOperation represents a dictionary with operations of events,
// where the key is an event and the value is an action.
var EventsOperation = (func() map[string]opencdc.Operation {
	eventsOperation := make(map[string]opencdc.Operation)

	for _, events := range EventsMap {
		for _, event := range events {
			switch e := event; {
			case strings.Contains(e, KeyCreated):
				eventsOperation[e] = opencdc.OperationCreate
			case strings.Contains(e, KeyDeleted):
				eventsOperation[e] = opencdc.OperationDelete
			default:
				eventsOperation[e] = opencdc.OperationUpdate
			}
		}
	}

	return eventsOperation
})()

// Reverse reverses an EventsData.
func (e EventsData) Reverse() {
	for i, j := 0, len(e)-1; i < j; i, j = i+1, j-1 {
		e[i], e[j] = e[j], e[i]
	}
}
