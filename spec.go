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

package stripe

import (
	sdk "github.com/conduitio/conduit-connector-sdk"

	"github.com/ConduitIO/conduit-connector-stripe/config"
)

type Spec struct{}

// Specification returns the connector's specification.
func Specification() sdk.Specification {
	return sdk.Specification{
		Name:    "stripe",
		Summary: "A Stripe source plugin for Conduit, written in Go.",
		Version: "v0.1.0",
		Author:  "Meroxa, Inc.",
		SourceParams: map[string]sdk.Parameter{
			config.SecretKey: {
				Default:     "",
				Required:    true,
				Description: "Stripe secret key.",
			},
			config.ResourceName: {
				Default:     "",
				Required:    true,
				Description: "Stripe resource name.",
			},
		},
	}
}
