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

package http

import (
	"encoding/json"
	"fmt"

	"github.com/ConduitIO/conduit-connector-stripe/config"
)

const stripeAPIURL = "https://api.stripe.com/v1"

// A StripeResponse represents a response from Stripe.
type StripeResponse struct {
	Data    []map[string]interface{} `json:"data"`
	HasMore bool                     `json:"has_more"`
}

// GetResources returns Stripe resources.
func (h http) GetResources(startingAfter string) (StripeResponse, error) {
	var resp StripeResponse

	if startingAfter != "" {
		startingAfter = fmt.Sprintf("&starting_after=%s", startingAfter)
	}

	url := fmt.Sprintf("%s/%s?limit=%d%s",
		stripeAPIURL, config.ResourceNamesMap[h.cfg.ResourceName], h.cfg.Limit, startingAfter)

	header := make(map[string]string, 1)
	header["Authorization"] = fmt.Sprintf("Bearer %s", h.cfg.SecretKey)

	data, err := h.get(url, header)
	if err != nil {
		return resp, fmt.Errorf("get data from stripe, by url and header: %w", err)
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, fmt.Errorf("unmarshal response data: %w", err)
	}

	return resp, nil
}
