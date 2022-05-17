### Supported Stripe resources

| resource            | events  |
|-----------------|---------|
| [`credit_note`](https://stripe.com/docs/api/credit_notes) | `credit_note.created`, `credit_note.updated`, `credit_note.voided` |
| [`billing_portal.configuration`](https://stripe.com/docs/api/customer_portal/configuration) | `billing_portal.configuration.created`, `billing_portal.configuration.updated` |
| [`invoice`](https://stripe.com/docs/api/invoices) | `invoice.created`, `invoice.deleted`, `invoice.finalization_failed`, `invoice.finalized`, `invoice.marked_uncollectible`, `invoice.paid`, `invoice.payment_action_required`, `invoice.payment_failed`, `invoice.payment_succeeded`, `invoice.sent`, `invoice.upcoming`, `invoice.updated`, `invoice.voided` |
| [`invoiceitem`](https://stripe.com/docs/api/invoiceitems) | `invoiceitem.created`, `invoiceitem.deleted`, `invoiceitem.updated` |
| [`plan`](https://stripe.com/docs/api/plans) | `plan.created`, `plan.deleted`, `plan.updated` |
| [`quote`](https://stripe.com/docs/api/quotes) | `quote.accepted`, `quote.canceled`, `quote.created`, `quote.finalized` |
| [`subscription`](https://stripe.com/docs/api/subscriptions) | `customer.subscription.created`, `customer.subscription.deleted`, `customer.subscription.pending_update_applied`, `customer.subscription.pending_update_expired`, `customer.subscription.trial_will_end`, `customer.subscription.updated` |
| [`subscription_schedule`](https://stripe.com/docs/api/subscription_schedules) | `subscription_schedule.aborted`, `subscription_schedule.canceled`, `subscription_schedule.completed`, `subscription_schedule.created`, `subscription_schedule.expiring`, `subscription_schedule.released`, `subscription_schedule.updated` |
| [`checkout.session`](https://stripe.com/docs/api/checkout/sessions) | `checkout.session.async_payment_failed`, `checkout.session.async_payment_succeeded`, `checkout.session.completed`, `checkout.session.expired` |
| [`account`](https://stripe.com/docs/api/accounts) | `account.updated` |
| [`application_fee`](https://stripe.com/docs/api/application_fees) | `application_fee.created`, `application_fee.refunded` |
| [`country_spec`](https://stripe.com/docs/api/country_specs) | |
| [`topup`](https://stripe.com/docs/api/topups) | `topup.canceled`, `topup.created`, `topup.failed`, `topup.reversed`, `topup.succeeded` |
| [`transfer`](https://stripe.com/docs/api/transfers) | `transfer.created`, `transfer.failed`, `transfer.paid`, `transfer.reversed`, `transfer.updated` |
| [`balance_transaction`](https://stripe.com/docs/api/balance_transactions) | |
| [`charge`](https://stripe.com/docs/api/charges) | `charge.captured`, `charge.expired`, `charge.failed`, `charge.pending`, `charge.refunded`, `charge.succeeded`, `charge.updated` |
| [`customer`](https://stripe.com/docs/api/customers) | `customer.created`, `customer.deleted`, `customer.updated` |
| [`dispute`](https://stripe.com/docs/api/disputes) | `charge.dispute.closed`, `charge.dispute.created`, `charge.dispute.funds_reinstated`, `charge.dispute.funds_withdrawn`, `charge.dispute.updated` |
| [`event`](https://stripe.com/docs/api/events) | |
| [`file`](https://stripe.com/docs/api/files) | `file.created` |
| [`file_link`](https://stripe.com/docs/api/file_links) | |
| [`payment_intent`](https://stripe.com/docs/api/payment_intents) | `payment_intent.amount_capturable_updated`, `payment_intent.canceled`, `payment_intent.created`, `payment_intent.partially_funded`, `payment_intent.payment_failed`, `payment_intent.processing`, `payment_intent.requires_action`, `payment_intent.succeeded` |
| [`setup_intent`](https://stripe.com/docs/api/setup_intents) | `setup_intent.canceled`, `setup_intent.created`, `setup_intent.requires_action`, `setup_intent.setup_failed`, `setup_intent.succeeded` |
| [`setup_attempt`](https://stripe.com/docs/api/setup_attempts) | |
| [`payout`](https://stripe.com/docs/api/payouts) | `payout.canceled`, `payout.created`, `payout.failed`, `payout.paid`, `payout.updated` |
| [`refund`](https://stripe.com/docs/api/refunds) | `charge.refund.updated` |
| [`linked_account`](https://stripe.com/docs/api/financial_connections) | |
| [`radar.early_fraud_warning`](https://stripe.com/docs/api/radar/early_fraud_warnings) | `radar.early_fraud_warning.created`, `radar.early_fraud_warning.updated` |
| [`review`](https://stripe.com/docs/api/radar/reviews) | `review.closed`, `review.opened` |
| [`radar.value_list`](https://stripe.com/docs/api/radar/value_lists) | |
| [`radar.value_list_item`](https://stripe.com/docs/api/radar/value_list_items) | |
| [`identity.verification_session`](https://stripe.com/docs/api/identity/verification_sessions) | `identity.verification_session.canceled`, `identity.verification_session.created`, `identity.verification_session.processing`, `identity.verification_session.redacted`, `identity.verification_session.requires_input`, `identity.verification_session.verified` |
| [`identity.verification_report`](https://stripe.com/docs/api/identity/verification_reports) | |
| [`issuing.authorization`](https://stripe.com/docs/api/issuing/authorizations) | `issuing_authorization.created`, `issuing_authorization.request`, `issuing_authorization.updated` |
| [`issuing.cardholder`](https://stripe.com/docs/api/issuing/cardholders) | `issuing_cardholder.created`, `issuing_cardholder.updated` |
| [`issuing.card`](https://stripe.com/docs/api/issuing/cards) | `issuing_card.created`, `issuing_card.updated` |
| [`issuing.dispute`](https://stripe.com/docs/api/issuing/disputes) | `issuing_dispute.closed`, `issuing_dispute.created`, `issuing_dispute.funds_reinstated`, `issuing_dispute.submitted`, `issuing_dispute.updated` |
| [`funding_instruction`](https://stripe.com/docs/api/issuing/funding_instructions) | |
| [`issuing.transaction`](https://stripe.com/docs/api/issuing/transactions) | `issuing_transaction.created`, `issuing_transaction.updated` |
| [`order`](https://stripe.com/docs/api/orders_v2) | `order.created`, `order.payment_failed`, `order.payment_succeeded`, `order.updated` |
| [`payment_link`](https://stripe.com/docs/api/payment_links) | `payment_link.created`, `payment_link.updated` |
| [`payment_method`](https://stripe.com/docs/api/payment_methods) | `payment_method.attached`, `payment_method.automatically_updated`, `payment_method.detached`, `payment_method.updated` |
| [`product`](https://stripe.com/docs/api/products) | `product.created`, `product.deleted`, `product.updated` |
| [`price`](https://stripe.com/docs/api/prices) | `price.created`, `price.deleted`, `price.updated` |
| [`coupon`](https://stripe.com/docs/api/coupons) | `coupon.created`, `coupon.deleted`, `coupon.updated` |
| [`promotion_code`](https://stripe.com/docs/api/promotion_codes) | `promotion_code.created`, `promotion_code.updated` |
| [`tax_code`](https://stripe.com/docs/api/tax_codes) | |
| [`tax_rate`](https://stripe.com/docs/api/tax_rates) | `tax_rate.created`, `tax_rate.updated` |
| [`shipping_rate`](https://stripe.com/docs/api/shipping_rates) | |
| [`reporting.report_run`](https://stripe.com/docs/api/reporting/report_run) | `reporting.report_run.failed`, `reporting.report_run.succeeded` |
| [`reporting.report_type`](https://stripe.com/docs/api/reporting/report_type) | `reporting.report_type.updated` |
| [`scheduled_query_run`](https://stripe.com/docs/api/sigma/scheduled_queries) | `sigma.scheduled_query_run.created` |
| [`terminal.location`](https://stripe.com/docs/api/terminal/locations) | |
| [`terminal.reader`](https://stripe.com/docs/api/terminal/readers) | `terminal.reader.action_failed`, `terminal.reader.action_succeeded` |
| [`terminal.configuration`](https://stripe.com/docs/api/terminal/configuration) | |
| [`financial_account`](https://stripe.com/docs/api/treasury/financial_accounts) | |
| [`transaction`](https://stripe.com/docs/api/treasury/transactions) | |
| [`transaction_entry`](https://stripe.com/docs/api/treasury/transaction_entries) | |
| [`outbound_transfer`](https://stripe.com/docs/api/treasury/outbound_transfers) | |
| [`outbound_payment`](https://stripe.com/docs/api/treasury/outbound_payments) | |
| [`inbound_transfer`](https://stripe.com/docs/api/treasury/inbound_transfers) | |
| [`received_credit`](https://stripe.com/docs/api/treasury/received_credits) | |
| [`received_debit`](https://stripe.com/docs/api/treasury/received_debits) | |
| [`credit_reversal`](https://stripe.com/docs/api/treasury/credit_reversals) | |
| [`debit_reversal`](https://stripe.com/docs/api/treasury/debit_reversals) | |
| [`webhook_endpoint`](https://stripe.com/docs/api/webhook_endpoints) | |