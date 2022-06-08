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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/conduitio/conduit-connector-stripe/config"
	"github.com/conduitio/conduit-connector-stripe/models"
	r "github.com/conduitio/conduit-connector-stripe/models/resources"
	"github.com/conduitio/conduit-connector-stripe/source"
	"github.com/google/uuid"
	"github.com/hashicorp/go-retryablehttp"
	"go.uber.org/goleak"
)

var (
	ctx context.Context
	cfg map[string]string

	clients = make(map[string]interface{})
)

// ConfigurableAcceptanceTestDriver driver for the test.
type ConfigurableAcceptanceTestDriver struct {
	sdk.ConfigurableAcceptanceTestDriver
}

// WriteToSource returns a slice of records that should be prepared in the Stripe so that the source will read them.
func (d ConfigurableAcceptanceTestDriver) WriteToSource(_ *testing.T, records []sdk.Record) []sdk.Record {
	return records
}

// GenerateRecord generates a new Stripe record.
func (d ConfigurableAcceptanceTestDriver) GenerateRecord(t *testing.T) sdk.Record {
	const (
		id               = "id"
		nameValue        = "client-%s"
		descriptionValue = "info about the %s"
	)

	var (
		name        = fmt.Sprintf(nameValue, uuid.New().String())
		description = fmt.Sprintf(descriptionValue, name)
	)

	resource, err := addResource(ctx, cfg, map[string]string{
		"name":        name,
		"description": description,
	})
	if err != nil {
		t.Error(err)
	}

	payload, _ := json.Marshal(resource)

	return sdk.Record{
		Position:  sdk.Position(uuid.New().String()),
		Metadata:  nil,
		CreatedAt: time.Now(),
		Key: sdk.StructuredData{
			id: resource[id],
		},
		Payload: sdk.RawData(payload),
	}
}

func TestAcceptance(t *testing.T) {
	const resourceName = r.CustomerResource

	ctx = context.Background()

	secretKey := os.Getenv("STRIPE_SECRET_KEY")
	if secretKey == "" {
		t.Skip("STRIPE_SECRET_KEY env var must be set")
	}

	cfg = map[string]string{
		config.SecretKey:    secretKey,
		config.ResourceName: resourceName,
	}

	sdk.AcceptanceTest(t, ConfigurableAcceptanceTestDriver{sdk.ConfigurableAcceptanceTestDriver{
		Config: sdk.ConfigurableAcceptanceTestDriverConfig{
			Connector: sdk.Connector{
				NewSpecification: Specification,
				NewSource:        source.NewSource,
				NewDestination:   nil,
			},
			SourceConfig:      cfg,
			DestinationConfig: nil,
			BeforeTest: func(t *testing.T) {
				if err := isEmpty(ctx, cfg); err != nil {
					t.Error(err)
				}
			},
			AfterTest: func(t *testing.T) {
				if err := clearResources(ctx, cfg); err != nil {
					t.Error(err)
				}
			},
			GoleakOptions: []goleak.Option{
				goleak.IgnoreTopFunction("net/http.(*persistConn).writeLoop"),
				goleak.IgnoreTopFunction("internal/poll.runtime_pollWait"),
			},
		},
	}},
	)
}

func addResource(ctx context.Context, cfg, params map[string]string) (map[string]interface{}, error) {
	var resource map[string]interface{}

	data, err := makeRequest(ctx, http.MethodPost, "", cfg, params)
	if err != nil {
		return nil, fmt.Errorf("make post request: %w", err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if len(resource) == 0 {
		return nil, errors.New("response is empty")
	}

	clients[resource["id"].(string)] = nil

	return resource, nil
}

func isEmpty(ctx context.Context, cfg map[string]string) error {
	var resource models.ResourceResponse

	data, err := makeRequest(ctx, http.MethodGet, "", cfg, nil)
	if err != nil {
		return fmt.Errorf("make get request: %w", err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	if len(resource.Data) > 0 {
		return errors.New("stripe already has a resource")
	}

	return nil
}

func clearResources(ctx context.Context, cfg map[string]string) error {
	var err error

	for k := range clients {
		_, err = makeRequest(ctx, http.MethodDelete, k, cfg, nil)
		if err != nil {
			return fmt.Errorf("make delete request: %w", err)
		}
	}

	clients = make(map[string]interface{})

	return nil
}

func makeRequest(ctx context.Context, method, path string, cfg, params map[string]string) ([]byte, error) {
	reqURL, err := url.Parse(models.APIURL)
	if err != nil {
		return nil, fmt.Errorf("parse api url: %w", err)
	}

	reqURL.Path += fmt.Sprintf(models.PathFmt, models.ResourcesMap[cfg[config.ResourceName]])

	if path != "" {
		reqURL.Path += fmt.Sprintf(models.PathFmt, path)
	}

	values := reqURL.Query()
	for k, v := range params {
		values.Add(k, v)
	}

	reqURL.RawQuery = values.Encode()

	cli := retryablehttp.NewClient()

	req, err := retryablehttp.NewRequestWithContext(ctx, method, reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create new request: %w", err)
	}
	req.Header.Add(models.HeaderAuthKey, fmt.Sprintf(models.HeaderAuthValueFormat, cfg[config.SecretKey]))

	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read all response body: %w", err)
	}

	return data, nil
}
