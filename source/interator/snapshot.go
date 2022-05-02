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

package interator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/hashicorp/go-retryablehttp"

	sdk "github.com/conduitio/conduit-connector-sdk"

	"github.com/ConduitIO/conduit-connector-stripe/clients/http"
	"github.com/ConduitIO/conduit-connector-stripe/source/position"
)

// A SnapshotIterator represents iteration over a slice of Stripe data.
type SnapshotIterator struct {
	httpClient *http.Client
	position   position.Position
	response   *Response
	index      int
}

// A Response represents a response from Stripe.
type Response struct {
	Data    []map[string]interface{} `json:"data"`
	HasMore bool                     `json:"has_more"`
}

// NewSnapshotIterator returns SnapshotIterator.
func NewSnapshotIterator(cli *http.Client, pos position.Position) *SnapshotIterator {
	return &SnapshotIterator{
		httpClient: cli,
		position:   pos,
	}
}

// Next returns the next record.
func (i *SnapshotIterator) Next() (sdk.Record, error) {
	if i.response == nil || len(i.response.Data) == i.index {
		if i.response != nil && !i.position.HasMore {
			i.index = 0

			return sdk.Record{}, nil
		}

		if err := i.getData(); err != nil {
			return sdk.Record{}, fmt.Errorf("get response data: %w", err)
		}
	}

	payload, err := json.Marshal(i.response.Data[i.index])
	if err != nil {
		return sdk.Record{}, fmt.Errorf("marshal payload: %w", err)
	}

	output := sdk.Record{
		Position: i.position.FormatSDKPosition(),
		Metadata: map[string]string{
			"action": "insert",
		},
		CreatedAt: time.Unix(int64(i.response.Data[i.index]["created"].(float64)), 0),
		Key:       sdk.RawData(i.response.Data[i.index]["id"].(string)),
		Payload:   sdk.RawData(payload),
	}

	i.position.StartingAfter = i.response.Data[i.index]["id"].(string)
	i.index++

	return output, nil
}

func (i *SnapshotIterator) getData() error {
	startingAfter := ""
	if i.position.StartingAfter != "" {
		startingAfter = fmt.Sprintf("&starting_after=%s", i.position.StartingAfter)
	}

	req, err := retryablehttp.NewRequest("GET",
		fmt.Sprintf("https://api.stripe.com/v1/%s?limit=%d%s",
			i.httpClient.Config.ResourceName, i.httpClient.Config.Limit, startingAfter), nil)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", i.httpClient.Config.SecretKey))

	resp, err := i.httpClient.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read all response data: %w", err)
	}

	err = json.Unmarshal(body, &i.response)
	if err != nil {
		return fmt.Errorf("unmarshal response data: %w", err)
	}

	i.position.HasMore = i.response.HasMore
	i.index = 0

	return nil
}
