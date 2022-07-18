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
	"net/url"
	"strconv"

	"github.com/conduitio-labs/conduit-connector-stripe/clients/http"
	"github.com/conduitio-labs/conduit-connector-stripe/config"
	"github.com/conduitio-labs/conduit-connector-stripe/models"
)

const (
	pathEvents       = "/events"
	batchSize        = "limit"
	startingAfterKey = "starting_after"
	endingBeforeKey  = "ending_before"
	typesKey         = "types[]"
	createdKey       = "created[gt]"
)

// A Stripe represents Stripe client struct.
type Stripe struct {
	cfg     config.Config
	httpCli http.Client
}

// New initialises a new Stripe client.
func New(cfg config.Config, httpCli http.Client) Stripe {
	return Stripe{
		cfg:     cfg,
		httpCli: httpCli,
	}
}

// GetResource returns a list of resource objects.
func (s Stripe) GetResource(startingAfter string) (models.ResourceResponse, error) {
	var resp models.ResourceResponse

	reqURL, err := url.Parse(models.APIURL)
	if err != nil {
		return resp, fmt.Errorf("parse api url: %w", err)
	}

	reqURL.Path += fmt.Sprintf(models.PathFmt, models.ResourcesMap[s.cfg.ResourceName])

	values := reqURL.Query()
	values.Add(batchSize, strconv.Itoa(s.cfg.BatchSize))

	if startingAfter != "" {
		values.Add(startingAfterKey, startingAfter)
	}

	reqURL.RawQuery = values.Encode()

	header := make(map[string]string, 1)
	header[models.HeaderAuthKey] = fmt.Sprintf(models.HeaderAuthValueFormat, s.cfg.SecretKey)

	data, err := s.httpCli.Get(reqURL.String(), header)
	if err != nil {
		return resp, fmt.Errorf("get data from stripe, by url %s and header: %w", reqURL.String(), err)
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, fmt.Errorf("unmarshal response data: %w", err)
	}

	return resp, nil
}

// GetEvent returns a list of event objects.
func (s Stripe) GetEvent(createdAt int64, startingAfter, endingBefore string) (models.EventResponse, error) {
	var resp models.EventResponse

	reqURL, err := url.Parse(models.APIURL)
	if err != nil {
		return resp, fmt.Errorf("parse api url: %w", err)
	}

	reqURL.Path += pathEvents

	values := reqURL.Query()
	values.Add(createdKey, strconv.FormatInt(createdAt, 10))
	values.Add(batchSize, strconv.Itoa(s.cfg.BatchSize))

	if events, ok := models.EventsMap[s.cfg.ResourceName]; ok {
		for i := range events {
			values.Add(typesKey, events[i])
		}
	}

	if startingAfter != "" {
		values.Add(startingAfterKey, startingAfter)
	}

	if endingBefore != "" {
		values.Add(endingBeforeKey, endingBefore)
	}

	reqURL.RawQuery = values.Encode()

	header := make(map[string]string, 1)
	header[models.HeaderAuthKey] = fmt.Sprintf(models.HeaderAuthValueFormat, s.cfg.SecretKey)

	data, err := s.httpCli.Get(reqURL.String(), header)
	if err != nil {
		return resp, fmt.Errorf("get data from stripe, by url %s and header: %w", reqURL.String(), err)
	}

	err = json.Unmarshal(data, &resp)
	if err != nil {
		return resp, fmt.Errorf("unmarshal response data: %w", err)
	}

	return resp, nil
}
