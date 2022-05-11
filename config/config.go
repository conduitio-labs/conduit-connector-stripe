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

package config

import (
	"strconv"
	"time"
)

const (
	RetryMaxDefault                    = 3
	LimitDefault                       = 50
	PollingPeriodDefault time.Duration = 1000000000 // 1 second in nanoseconds
)

const (
	// SecretKey is the configuration name for Stripe secret key.
	SecretKey = "stripe.secretKey"

	// ResourceName is the configuration name for Stripe resource.
	ResourceName = "stripe.resourceName"

	// HTTPClientRetryMax is the configuration name for the maximum number of retries in the HTTP client.
	HTTPClientRetryMax = "stripe.http_client_retry_max"

	// Limit is the configuration name for the number of objects returned by the query to Stripe.
	Limit = "stripe.limit"

	// PollingPeriod is the configuration name for the period between requests to Stripe in the CDC iterator.
	PollingPeriod = "stripe.polling_period"
)

// A Config represents the configuration needed for Stripe.
type Config struct {
	SecretKey          string        `validate:"required"`
	ResourceName       string        `validate:"required,resource_name"`
	HTTPClientRetryMax int           `validate:"gte=1,lte=10"`
	Limit              int           `validate:"gte=1,lte=100"`
	PollingPeriod      time.Duration `validate:"gte=1s,lte=1h"`
}

// Parse parses Stripe configuration into a Config struct.
func Parse(cfg map[string]string) (Config, error) {
	config := Config{
		SecretKey:    cfg[SecretKey],
		ResourceName: cfg[ResourceName],
	}

	config.HTTPClientRetryMax = RetryMaxDefault
	if cfg[HTTPClientRetryMax] != "" {
		retryMax, err := strconv.Atoi(cfg[HTTPClientRetryMax])
		if err != nil {
			return Config{}, config.IntegerTypeConfigErr(HTTPClientRetryMax)
		}

		config.HTTPClientRetryMax = retryMax
	}

	config.Limit = LimitDefault
	if cfg[Limit] != "" {
		limit, err := strconv.Atoi(cfg[Limit])
		if err != nil {
			return Config{}, config.IntegerTypeConfigErr(Limit)
		}

		config.Limit = limit
	}

	config.PollingPeriod = PollingPeriodDefault
	if cfg[PollingPeriod] != "" {
		pollingPeriod, err := time.ParseDuration(cfg[PollingPeriod])
		if err != nil {
			return Config{}, config.PollingPeriodIsNotDurationErr(PollingPeriod)
		}

		if pollingPeriod <= 0 {
			return Config{}, config.PollingPeriodPositiveErr(PollingPeriod)
		}

		config.PollingPeriod = pollingPeriod
	}

	return config, config.Validate()
}
