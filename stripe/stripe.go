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
	"encoding/json"
	"fmt"

	"github.com/conduitio/conduit-connector-stripe/clients/http"
	"github.com/conduitio/conduit-connector-stripe/config"
	"github.com/conduitio/conduit-connector-stripe/models"
)

const stripeAPIURL = "https://api.stripe.com/v1"

type Stripe struct {
	cfg     *config.Config
	httpCli http.Client
}

// New initialises a new Stripe client.
func New(cfg *config.Config) Stripe {
	return Stripe{
		cfg:     cfg,
		httpCli: http.NewClient(cfg),
	}
}

// GetResource returns a list of resource objects from Stripe.
func (s Stripe) GetResource(startingAfter string) (models.StripeResponse, error) {
	var resp models.StripeResponse

	if startingAfter != "" {
		startingAfter = fmt.Sprintf("&starting_after=%s", startingAfter)
	}

	url := fmt.Sprintf("%s/%s?limit=%d%s",
		stripeAPIURL, config.ResourceNamesMap[s.cfg.ResourceName], s.cfg.Limit, startingAfter)

	header := make(map[string]string, 1)
	header["Authorization"] = fmt.Sprintf("Bearer %s", s.cfg.SecretKey)

	data, err := s.httpCli.Get(url, header)
	if err != nil {
		return resp, fmt.Errorf("get data from stripe, by url and header: %w", err)
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, fmt.Errorf("unmarshal response data: %w", err)
	}

	return resp, nil
}
