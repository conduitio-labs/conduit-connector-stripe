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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/conduitio/conduit-connector-stripe/models"
	"github.com/hashicorp/go-retryablehttp"
)

// A Client represents retryable http client.
type Client struct {
	httpClient *retryablehttp.Client
}

// NewClient returns a new retryable http client.
func NewClient(ctx context.Context) Client {
	retryClient := retryablehttp.NewClient()
	retryClient.Logger = sdk.Logger(ctx)

	return Client{
		httpClient: retryClient,
	}
}

// Get makes a GET http-request to the URL with headers.
func (cli Client) Get(url string, header ...map[string]string) ([]byte, error) {
	req, err := retryablehttp.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create new request: %w", err)
	}

	for i := range header {
		for k, v := range header[i] {
			req.Header.Add(k, v)
		}
	}

	resp, err := cli.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read all response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		errResp := models.ErrorResponse{}

		err = json.Unmarshal(data, &errResp)
		if err != nil {
			return nil, fmt.Errorf("unmarshal response: %w", err)
		}

		if errResp.Error.Message != "" {
			return nil, errors.New(errResp.Error.Message)
		}

		return nil, fmt.Errorf(models.UnexpectedErrorWithStatusCode, resp.StatusCode)
	}

	return data, nil
}

// Close closes any connections which were previously connected from previous requests.
func (cli Client) Close() {
	if cli.httpClient != nil {
		cli.httpClient.HTTPClient.CloseIdleConnections()
	}
}
