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

//go:generate paramgen -output=paramgen.go Config

package config

import (
	"fmt"

	"github.com/conduitio-labs/conduit-connector-stripe/models"
)

type Config struct {
	// SecretKey is the configuration name for Stripe secret key.
	SecretKey string `json:"secretKey" validate:"required"`
	// ResourceName is the configuration name for Stripe resource.
	ResourceName string `json:"resourceName" validate:"required"`
	// BatchSize is the configuration name for the number of objects in the batch returned from Stripe.
	BatchSize int `json:"batchSize" default:"10" validate:"gt=0,lt=100001"`
	// Snapshot is the configuration name for the Snapshot field.
	Snapshot bool `json:"snapshot" default:"true"`
}

// Validate executes manual validations beyond what is defined in struct tags.
func (c *Config) Validate() error {
	// c.SecretKey has required validation handled in struct tag

	// c.ResourceName required validation is handled in stuct tag
	// handling "resource_name" validation
	_, ok := models.ResourcesMap[c.ResourceName]
	if !ok {
		return fmt.Errorf("%q wrong resource name", c.ResourceName)
	}

	return nil
}
