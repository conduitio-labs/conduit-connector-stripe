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
	BalanceTransactionResource = "balance_transaction"
	BalanceTransactionsList    = "balance_transactions"

	ChargeResource       = "charge"
	ChargesList          = "charges"
	ChargeCapturedEvent  = "charge.captured"
	ChargeExpiredEvent   = "charge.expired"
	ChargeFailedEvent    = "charge.failed"
	ChargePendingEvent   = "charge.pending"
	ChargeRefundedEvent  = "charge.refunded"
	ChargeSucceededEvent = "charge.succeeded"
	ChargeUpdatedEvent   = "charge.updated"

	CustomerResource     = "customer"
	CustomersList        = "customers"
	CustomerCreatedEvent = "customer.created"
	CustomerDeletedEvent = "customer.deleted"
	CustomerUpdatedEvent = "customer.updated"

	DisputeResource             = "dispute"
	DisputesList                = "disputes"
	DisputeClosedEvent          = "charge.dispute.closed"
	DisputeCreatedEvent         = "charge.dispute.created"
	DisputeFundsReinstatedEvent = "charge.dispute.funds_reinstated"
	DisputeFundsWithdrawnEvent  = "charge.dispute.funds_withdrawn"
	DisputeupdatedEvent         = "charge.dispute.updated"

	EventResource = "event"
	EventsList    = "events"

	FileResource     = "file"
	FilesList        = "files"
	FileCreatedEvent = "file.created"

	FileLinkResource = "file_link"
	FileLinksList    = "file_links"

	PaymentIntentResource                     = "payment_intent"
	PaymentIntentsList                        = "payment_intents"
	PaymentIntentAmountCapturableUpdatedEvent = "payment_intent.amount_capturable_updated"
	PaymentIntentCanceledEvent                = "payment_intent.canceled"
	PaymentIntentCreatedEvent                 = "payment_intent.created"
	PaymentIntentPartiallyFundedEvent         = "payment_intent.partially_funded"
	PaymentIntentPaymentFailedEvent           = "payment_intent.payment_failed"
	PaymentIntentProcessingEvent              = "payment_intent.processing"
	PaymentIntentRequiresActionEvent          = "payment_intent.requires_action"
	PaymentIntentSucceededEvent               = "payment_intent.succeeded"

	SetupIntentResource            = "setup_intent"
	SetupIntentsList               = "setup_intents"
	SetupIntentCanceledEvent       = "setup_intent.canceled"
	SetupIntentCreatedEvent        = "setup_intent.created"
	SetupIntentRequiresActionEvent = "setup_intent.requires_action"
	SetupIntentSetupFailedEvent    = "setup_intent.setup_failed"
	SetupIntentSucceededEvent      = "setup_intent.succeeded"

	SetupAttemptResource = "setup_attempt"
	SetupAttemptsList    = "setup_attempts"

	PayoutResource       = "payout"
	PayoutsList          = "payouts"
	PayoutsCanceledEvent = "payout.canceled"
	PayoutsCreatedEvent  = "payout.created"
	PayoutsFailedEvent   = "payout.failed"
	PayoutsPaidEvent     = "payout.paid"
	PayoutsUpdatedEvent  = "payout.updated"

	RefundResource     = "refund"
	RefundsList        = "refunds"
	RefundUpdatedEvent = "charge.refund.updated"
)
