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
	AccountResource     = "account"
	AccountsList        = "accounts"
	AccountUpdatedEvent = "account.updated"

	ApplicationFeeResource      = "application_fee"
	ApplicationFeesList         = "application_fees"
	ApplicationFeeCreatedEvent  = "application_fee.created"
	ApplicationFeeRefundedEvent = "application_fee.refunded"

	TopUpResource       = "topup"
	TopUpsList          = "topups"
	TopupCanceledEvent  = "topup.canceled"
	TopupCreatedEvent   = "topup.created"
	TopupFailedEvent    = "topup.failed"
	TopupReversedEvent  = "topup.reversed"
	TopupSucceededEvent = "topup.succeeded"

	TransferResource      = "transfer"
	TransfersList         = "transfers"
	TransferCreatedEvent  = "transfer.created"
	TransferFailedEvent   = "transfer.failed"
	TransferPaidEvent     = "transfer.paid"
	TransferReversedEvent = "transfer.reversed"
	TransferUpdatedEvent  = "transfer.updated"
)

var (
	AccountEvents = []string{
		AccountUpdatedEvent,
	}

	ApplicationFeeEvents = []string{
		ApplicationFeeCreatedEvent,
		ApplicationFeeRefundedEvent,
	}

	TopUpEvents = []string{
		TopupCanceledEvent,
		TopupCreatedEvent,
		TopupFailedEvent,
		TopupReversedEvent,
		TopupSucceededEvent,
	}

	TransferEvents = []string{
		TransferCreatedEvent,
		TransferFailedEvent,
		TransferPaidEvent,
		TransferReversedEvent,
		TransferUpdatedEvent,
	}
)
