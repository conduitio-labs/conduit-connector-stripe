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

	"github.com/conduitio/conduit-connector-stripe/validator"
)

const (
	// SecretKey is the configuration name for Stripe secret key.
	SecretKey = "secretKey"

	// ResourceName is the configuration name for Stripe resource.
	ResourceName = "resourceName"

	// BatchSize is the configuration name for the number of objects in the batch returned from Stripe.
	BatchSize = "batchSize"
)

// A Config represents the configuration needed for Stripe.
type Config struct {
	SecretKey    string `validate:"required"`
	ResourceName string `validate:"required,resource_name"`
	BatchSize    int
}

// Parse parses Stripe configuration into a Config struct.
func Parse(cfg map[string]string) (Config, error) {
	config := Config{
		SecretKey:    cfg[SecretKey],
		ResourceName: cfg[ResourceName],
	}

	if cfg[BatchSize] != "" {
		batchSize, err := strconv.Atoi(cfg[BatchSize])
		if err != nil {
			return Config{}, validator.IntegerTypeConfigErr(BatchSize)
		}

		if batchSize <= 0 || batchSize > 100 {
			return Config{}, validator.InvalidBatchSizeErr(BatchSize)
		}

		config.BatchSize = batchSize
	}

	err := config.Validate()
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
