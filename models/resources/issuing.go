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
	IssuingAuthorizationResource     = "issuing.authorization"
	IssuingAuthorizationsList        = "issuing/authorizations"
	IssuingAuthorizationCreatedEvent = "issuing_authorization.created"
	IssuingAuthorizationRequestEvent = "issuing_authorization.request"
	IssuingAuthorizationUpdatedEvent = "issuing_authorization.updated"

	IssuingCardholderResource     = "issuing.cardholder"
	IssuingCardholdersList        = "issuing/cardholders"
	IssuingCardholderCreatedEvent = "issuing_cardholder.created"
	IssuingCardholderUpdatedEvent = "issuing_cardholder.updated"

	IssuingCardResource     = "issuing.card"
	IssuingCardsList        = "issuing/cards"
	IssuingCardCreatedEvent = "issuing_card.created"
	IssuingCardUpdatedEvent = "issuing_card.updated"

	IssuingDisputeResource             = "issuing.dispute"
	IssuingDisputesList                = "issuing/disputes"
	IssuingDisputeClosedEvent          = "issuing_dispute.closed"
	IssuingDisputeCreatedEvent         = "issuing_dispute.created"
	IssuingDisputeFundsReinstatedEvent = "issuing_dispute.funds_reinstated"
	IssuingDisputeSubmittedEvent       = "issuing_dispute.submitted"
	IssuingDisputeUpdatedEvent         = "issuing_dispute.updated"

	FundingInstructionResource = "funding_instruction"
	FundingInstructionsList    = "issuing/funding_instructions"

	IssuingTransactionResource     = "issuing.transaction"
	IssuingTransactionsList        = "issuing/transactions"
	IssuingTransactionCreatedEvent = "issuing_transaction.created"
	IssuingTransactionUpdatedEvent = "issuing_transaction.updated"
)
