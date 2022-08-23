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

	"github.com/conduitio-labs/conduit-connector-stripe/config"
	"github.com/conduitio-labs/conduit-connector-stripe/models"
	r "github.com/conduitio-labs/conduit-connector-stripe/models/resources"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/google/uuid"
	"github.com/hashicorp/go-retryablehttp"
)

const (
	resourceName = r.CustomerResource

	clientNameFmt        = "client-%s"
	clientDescriptionFmt = "info about the %s"
)

var (
	cfg map[string]string

	clients = make(map[string]interface{})
)

// AcceptanceTestDriver driver for the test.
type AcceptanceTestDriver struct {
	sdk.ConfigurableAcceptanceTestDriver
}

// WriteToSource returns a slice of records that should be prepared in the Stripe so that the source will read them.
func (d AcceptanceTestDriver) WriteToSource(t *testing.T, records []sdk.Record) []sdk.Record {
	ctx := context.Background()

	cli := retryablehttp.NewClient()
	cli.Logger = sdk.Logger(ctx)
	defer cli.HTTPClient.CloseIdleConnections()

	for i := range records {
		m := make(map[string]string)

		err := json.Unmarshal(records[i].Payload.After.Bytes(), &m)
		if err != nil {
			t.Error(err)
		}

		resource, err := addResource(ctx, cli, cfg, m)
		if err != nil {
			t.Error(err)
		}

		payload, err := json.Marshal(resource)
		if err != nil {
			t.Error(err)
		}

		records[i].Key = sdk.StructuredData{
			models.KeyID: resource[models.KeyID],
		}

		records[i].Payload.After = sdk.RawData(payload)
	}

	return records
}

// GenerateRecord generates a new Stripe record.
func (d AcceptanceTestDriver) GenerateRecord(t *testing.T, operation sdk.Operation) sdk.Record {
	var (
		name        = fmt.Sprintf(clientNameFmt, uuid.New().String())
		description = fmt.Sprintf(clientDescriptionFmt, name)
	)

	payload, _ := json.Marshal(map[string]string{
		models.KeyName:        name,
		models.KeyDescription: description,
	})

	return sdk.Record{
		Operation: operation,
		Payload:   sdk.Change{After: sdk.RawData(payload)},
	}
}

func TestAcceptance(t *testing.T) {
	ctx := context.Background()

	secretKey := os.Getenv("STRIPE_SECRET_KEY")
	if secretKey == "" {
		t.Skip("STRIPE_SECRET_KEY env var must be set")
	}

	cfg = map[string]string{
		config.SecretKey:    secretKey,
		config.ResourceName: resourceName,
	}

	sdk.AcceptanceTest(t, AcceptanceTestDriver{sdk.ConfigurableAcceptanceTestDriver{
		Config: sdk.ConfigurableAcceptanceTestDriverConfig{
			Connector:         Connector,
			SourceConfig:      cfg,
			DestinationConfig: nil,
			BeforeTest: func(t *testing.T) {
				cli := retryablehttp.NewClient()
				cli.Logger = sdk.Logger(ctx)
				defer cli.HTTPClient.CloseIdleConnections()

				if err := isEmpty(ctx, cli, cfg); err != nil {
					t.Error(err)
				}
			},
			AfterTest: func(t *testing.T) {
				cli := retryablehttp.NewClient()
				cli.Logger = sdk.Logger(ctx)
				defer cli.HTTPClient.CloseIdleConnections()

				if err := clearResources(ctx, cli, cfg); err != nil {
					t.Error(err)
				}
			},
		},
	}},
	)
}

func addResource(ctx context.Context, cli *retryablehttp.Client, cfg, params map[string]string,
) (map[string]interface{}, error) {
	var resource map[string]interface{}

	data, err := makeRequest(ctx, cli, http.MethodPost, "", cfg, params)
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

	clients[resource[models.KeyID].(string)] = nil

	return resource, nil
}

func isEmpty(ctx context.Context, cli *retryablehttp.Client, cfg map[string]string) error {
	var resource models.ResourceResponse

	data, err := makeRequest(ctx, cli, http.MethodGet, "", cfg, nil)
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

func clearResources(ctx context.Context, cli *retryablehttp.Client, cfg map[string]string) error {
	var err error

	for k := range clients {
		_, err = makeRequest(ctx, cli, http.MethodDelete, k, cfg, nil)
		if err != nil {
			return fmt.Errorf("make delete request: %w", err)
		}
	}

	clients = make(map[string]interface{})

	return nil
}

func makeRequest(ctx context.Context, cli *retryablehttp.Client, method, path string, cfg, params map[string]string,
) ([]byte, error) {
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
