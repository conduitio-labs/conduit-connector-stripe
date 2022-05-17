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
	ProductResource     = "product"
	ProductsList        = "products"
	ProductCreatedEvent = "product.created"
	ProductDeletedEvent = "product.deleted"
	ProductUpdatedEvent = "product.updated"

	PriceResource     = "price"
	PricesList        = "prices"
	PriceCreatedEvent = "price.created"
	PriceDeletedEvent = "price.deleted"
	PriceUpdatedEvent = "price.updated"

	CouponResource     = "coupon"
	CouponsList        = "coupons"
	CouponCreatedEvent = "coupon.created"
	CouponDeletedEvent = "coupon.deleted"
	CouponUpdatedEvent = "coupon.updated"

	PromotionCodeResource     = "promotion_code"
	PromotionCodesList        = "promotion_codes"
	PromotionCodeCreatedEvent = "promotion_code.created"
	PromotionCodeUpdatedEvent = "promotion_code.updated"

	TaxCodeResource = "tax_code"
	TaxCodesList    = "tax_codes"

	TaxRateResource     = "tax_rate"
	TaxRatesList        = "tax_rates"
	TaxRateCreatedEvent = "tax_rate.created"
	TaxRateUpdatedEvent = "tax_rate.updated"

	ShippingRateResource = "shipping_rate"
	ShippingRatesList    = "shipping_rates"
)

var (
	ProductEvents = []string{
		ProductCreatedEvent,
		ProductDeletedEvent,
		ProductUpdatedEvent,
	}

	PriceEvents = []string{
		PriceCreatedEvent,
		PriceDeletedEvent,
		PriceUpdatedEvent,
	}

	CouponEvents = []string{
		CouponCreatedEvent,
		CouponDeletedEvent,
		CouponUpdatedEvent,
	}

	PromotionCodeEvents = []string{
		PromotionCodeCreatedEvent,
		PromotionCodeUpdatedEvent,
	}

	TaxRateEvents = []string{
		TaxRateCreatedEvent,
		TaxRateUpdatedEvent,
	}
)
