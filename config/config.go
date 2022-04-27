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

const (
	// SecretKey is the config name for Stripe secret key.
	SecretKey = "stripe.secretKey"

	// ResourceName is the config name for Stripe resource.
	ResourceName = "stripe.resourceName"
)

// Config represents configuration needed for Stripe.
type Config struct {
	SecretKey    string `validate:"required"`
	ResourceName string `validate:"required"`
}

// Parse parses Stripe configuration into a Config struct.
func Parse(cfg map[string]string) (Config, error) {
	config := Config{
		SecretKey:    cfg[SecretKey],
		ResourceName: cfg[ResourceName],
	}

	return config, config.Validate()
}
