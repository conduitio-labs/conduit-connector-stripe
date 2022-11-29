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
	"fmt"
	"strconv"

	"github.com/conduitio-labs/conduit-connector-stripe/validator"
)

type SnapshotMode string

const (
	// SecretKey is the configuration name for Stripe secret key.
	SecretKey = "secretKey"
	// ResourceName is the configuration name for Stripe resource.
	ResourceName = "resourceName"
	// Snapshot is the configuration name for the Snapshot field.
	Snapshot = "snapshot"
	// BatchSize is the configuration name for the number of objects in the batch returned from Stripe.
	BatchSize = "batchSize"

	// BatchSizeDefault is the default value of the batch size.
	BatchSizeDefault = 10
	// SnapshotDefault is a default value for the Snapshot field.
	SnapshotDefault = true
)

// A Config represents the configuration needed for Stripe.
type Config struct {
	SecretKey    string `validate:"required"`
	ResourceName string `validate:"required,resource_name"`
	BatchSize    int    `validate:"gte=1,lte=100,omitempty"`
	Snapshot     bool
}

// Parse parses Stripe configuration into a Config struct.
func Parse(cfg map[string]string) (Config, error) {
	config := Config{
		SecretKey:    cfg[SecretKey],
		ResourceName: cfg[ResourceName],
		Snapshot:     SnapshotDefault,
		BatchSize:    BatchSizeDefault,
	}

	if cfg[Snapshot] != "" {
		snapshot, err := strconv.ParseBool(cfg[Snapshot])
		if err != nil {
			return Config{}, fmt.Errorf("parse %q: %w", Snapshot, err)
		}

		config.Snapshot = snapshot
	}

	if cfg[BatchSize] != "" {
		batchSize, err := strconv.Atoi(cfg[BatchSize])
		if err != nil {
			return Config{}, validator.IntegerTypeConfigErr(BatchSize)
		}

		config.BatchSize = batchSize
	}

	err := config.Validate()
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
