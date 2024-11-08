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

package source

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/conduitio-labs/conduit-connector-stripe/config"
	"github.com/conduitio-labs/conduit-connector-stripe/models"
	r "github.com/conduitio-labs/conduit-connector-stripe/models/resources"
	"github.com/conduitio/conduit-commons/opencdc"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/google/uuid"
	"github.com/hashicorp/go-retryablehttp"
	"go.uber.org/goleak"
)

const resourceName = r.CustomerResource

var clients = make(map[string]interface{})

func TestSource_Read(t *testing.T) { // nolint:gocyclo,nolintlint
	t.Run("read nothing", func(t *testing.T) {
		var ctx = context.Background()

		// Ignore go-cache's background goroutine used for periodic cleanup operations.
		// It's not a true leak as it's part of the cache's design.
		defer goleak.VerifyNone(t, goleak.IgnoreTopFunction("github.com/twmb/go-cache/cache.New[...].func1"))

		cfg, err := prepareDefaultConfig()
		if err != nil {
			t.Log(err)
			t.Skip()
		}

		cli := retryablehttp.NewClient()
		cli.Logger = sdk.Logger(ctx)
		defer cli.HTTPClient.CloseIdleConnections()

		err = isEmpty(ctx, cli, cfg)
		if err != nil {
			t.Errorf("check is empty: %s", err.Error())

			return
		}

		source := NewSource()

		err = source.Configure(ctx, cfg)
		if err != nil {
			t.Errorf("configure: %s", err.Error())
		}

		err = source.Open(ctx, nil)
		if err != nil {
			t.Errorf("open: %s", err.Error())
		}

		_, err = source.Read(ctx)
		if err != sdk.ErrBackoffRetry {
			t.Errorf("read: %s", err)
		}

		err = source.Teardown(ctx)
		if err != nil {
			t.Errorf("teardown: %s", err.Error())
		}
	})

	t.Run("invalid secret key", func(t *testing.T) {
		const (
			invalidSecretKey = "invalid_secret_key"
			expectedErr      = "populate with the resource: " +
				"get list of resource objects: " +
				"get data from stripe, by url https://api.stripe.com/v1/customers?limit=10 and header: " +
				"Invalid API Key provided: invalid_******_key"
		)

		var ctx = context.Background()

		defer goleak.VerifyNone(t, goleak.IgnoreTopFunction("github.com/twmb/go-cache/cache.New[...].func1"))

		cfg, err := prepareConfig(invalidSecretKey, "")
		if err != nil {
			t.Log(err)
			t.Skip()
		}

		source := NewSource()

		err = source.Configure(ctx, cfg)
		if err != nil {
			t.Errorf("configure: %s", err.Error())
		}

		err = source.Open(ctx, nil)
		if err != nil {
			t.Errorf("open: %s", err.Error())
		}

		_, err = source.Read(ctx)
		if err != nil {
			if err.Error() != expectedErr {
				t.Errorf("expected error \"%s\", got \"%s\"", expectedErr, err.Error())

				return
			}
		}

		err = source.Teardown(ctx)
		if err != nil {
			t.Errorf("teardown: %s", err.Error())
		}
	})

	// add resCount resources
	// initialize source
	// read resCount - stop resources
	// reinitialize source
	// read rest resources
	// read when there is no resources
	t.Run("snapshot iterator", func(t *testing.T) {
		var (
			ctx = context.Background()

			record opencdc.Record
			pos    opencdc.Position
		)

		defer goleak.VerifyNone(t, goleak.IgnoreTopFunction("github.com/twmb/go-cache/cache.New[...].func1"))

		cfg, err := prepareDefaultConfig()
		if err != nil {
			t.Log(err)
			t.Skip()
		}

		cli := retryablehttp.NewClient()
		cli.Logger = sdk.Logger(ctx)
		defer cli.HTTPClient.CloseIdleConnections()

		err = isEmpty(ctx, cli, cfg)
		if err != nil {
			t.Errorf("check is empty: %s", err.Error())

			return
		}

		defer func(ctx context.Context, cfg map[string]string) {
			err = clearResources(ctx, cli, cfg)
			if err != nil {
				t.Errorf("clear resources: %s", err.Error())
			}
		}(ctx, cfg)

		// add resCount resources
		genResCount := 12
		stopReadIndex := 5

		resources, err := generateResources(ctx, cli, cfg, genResCount)
		if err != nil {
			t.Errorf("prepare data: %s", err.Error())
		}

		source := NewSource()

		err = source.Configure(ctx, cfg)
		if err != nil {
			t.Errorf("configure: %s", err.Error())
		}

		err = source.Open(ctx, nil)
		if err != nil {
			t.Errorf("open: %s", err.Error())
		}

		for i := len(resources) - 1; i >= stopReadIndex; i-- {
			record, err = source.Read(ctx)
			if err != nil {
				t.Errorf("read: %s", err.Error())
			}

			pos = record.Position

			err = compareResult(record, resources[i], opencdc.OperationSnapshot)
			if err != nil {
				t.Error(err.Error())
			}
		}

		err = source.Teardown(ctx)
		if err != nil {
			t.Errorf("teardown: %s", err.Error())
		}

		source = NewSource()

		err = source.Configure(ctx, cfg)
		if err != nil {
			t.Errorf("configure: %s", err.Error())
		}

		err = source.Open(ctx, pos)
		if err != nil {
			t.Errorf("open: %s", err.Error())
		}

		for i := stopReadIndex - 1; i >= 0; i-- {
			record, err = source.Read(ctx)
			if err != nil {
				t.Errorf("read: %s", err.Error())
			}

			err = compareResult(record, resources[i], opencdc.OperationSnapshot)
			if err != nil {
				t.Error(err.Error())
			}
		}

		// read empty source
		_, err = source.Read(ctx)
		if err != nil && err != sdk.ErrBackoffRetry {
			t.Errorf("read: %s", err)
		}

		err = source.Teardown(ctx)
		if err != nil {
			t.Errorf("teardown: %s", err.Error())
		}
	})

	// add resCount resources
	// initialize source
	// read all resources by `Snapshot` iterator
	// reinitialize source
	// add resource
	// read resource by `CDC` iterator
	// read when there is no events
	t.Run("reinitialize source between snapshot and cdc iterators", func(t *testing.T) {
		var (
			ctx = context.Background()

			record opencdc.Record
			rp     opencdc.Position
		)

		defer goleak.VerifyNone(t, goleak.IgnoreTopFunction("github.com/twmb/go-cache/cache.New[...].func1"))

		cfg, err := prepareDefaultConfig()
		if err != nil {
			t.Log(err)
			t.Skip()
		}

		cli := retryablehttp.NewClient()
		cli.Logger = sdk.Logger(ctx)
		defer cli.HTTPClient.CloseIdleConnections()

		err = isEmpty(ctx, cli, cfg)
		if err != nil {
			t.Errorf("check is empty: %s", err.Error())

			return
		}

		defer func(ctx context.Context, cfg map[string]string) {
			err = clearResources(ctx, cli, cfg)
			if err != nil {
				t.Errorf("clear resources: %s", err.Error())
			}
		}(ctx, cfg)

		// add resCount resources
		genResCount := 3

		resources, err := generateResources(ctx, cli, cfg, genResCount)
		if err != nil {
			t.Errorf("prepare data: %s", err.Error())
		}

		source := NewSource()

		err = source.Configure(ctx, cfg)
		if err != nil {
			t.Errorf("configure: %s", err.Error())
		}

		err = source.Open(ctx, nil)
		if err != nil {
			t.Errorf("open: %s", err.Error())
		}

		for i := len(resources) - 1; i >= 0; i-- {
			record, err = source.Read(ctx)
			if err != nil {
				t.Errorf("read: %s", err.Error())
			}

			rp = record.Position

			err = compareResult(record, resources[i], opencdc.OperationSnapshot)
			if err != nil {
				t.Error(err.Error())
			}
		}

		err = source.Teardown(ctx)
		if err != nil {
			t.Errorf("teardown: %s", err.Error())
		}

		source = NewSource()

		err = source.Configure(ctx, cfg)
		if err != nil {
			t.Errorf("configure: %s", err.Error())
		}

		err = source.Open(ctx, rp)
		if err != nil {
			t.Errorf("open: %s", err.Error())
		}

		// add resource
		resources, err = generateResources(ctx, cli, cfg, 1)
		if err != nil {
			t.Errorf("prepare data: %s", err.Error())
		}

		record, err = source.Read(ctx)
		if err != nil {
			t.Errorf("read: %s", err.Error())
		}

		err = compareResult(record, resources[0], opencdc.OperationCreate)
		if err != nil {
			t.Error(err.Error())
		}

		// read empty source
		_, err = source.Read(ctx)
		if err != nil && err != sdk.ErrBackoffRetry {
			t.Errorf("read: %s", err)
		}

		err = source.Teardown(ctx)
		if err != nil {
			t.Errorf("teardown: %s", err.Error())
		}
	})

	// add resCount resources
	// initialize source
	// read all resources by `Snapshot` iterator
	// read when there is no resources
	// add resource
	// update resource
	// delete resource
	// reinitialize source
	// update resource
	t.Run("reinitialize cdc", func(t *testing.T) {
		var (
			ctx = context.Background()

			record opencdc.Record
		)

		defer goleak.VerifyNone(t, goleak.IgnoreTopFunction("github.com/twmb/go-cache/cache.New[...].func1"))

		cfg, err := prepareConfigWithBatchSize("7")
		if err != nil {
			t.Log(err)
			t.Skip()
		}

		cli := retryablehttp.NewClient()
		cli.Logger = sdk.Logger(ctx)
		defer cli.HTTPClient.CloseIdleConnections()

		err = isEmpty(ctx, cli, cfg)
		if err != nil {
			t.Errorf("check is empty: %s", err.Error())

			return
		}

		defer func(ctx context.Context, cfg map[string]string) {
			err = clearResources(ctx, cli, cfg)
			if err != nil {
				t.Errorf("clear resources: %s", err.Error())
			}
		}(ctx, cfg)

		// generate resCount resources
		genResCount := 5
		updateIndex1 := 0
		updateIndex2 := 2
		deleteIndex := 1

		resources, err := generateResources(ctx, cli, cfg, genResCount)
		if err != nil {
			t.Errorf("prepare data: %s", err.Error())
		}

		source := NewSource()

		err = source.Configure(ctx, cfg)
		if err != nil {
			t.Errorf("configure: %s", err.Error())
		}

		err = source.Open(ctx, nil)
		if err != nil {
			t.Errorf("open: %s", err.Error())
		}

		for i := len(resources) - 1; i >= 0; i-- {
			record, err = source.Read(ctx)
			if err != nil {
				t.Errorf("read: %s", err.Error())
			}

			err = compareResult(record, resources[i], opencdc.OperationSnapshot)
			if err != nil {
				t.Error(err.Error())
			}
		}

		updateID1 := resources[updateIndex1][models.KeyID].(string)
		updateID2 := resources[updateIndex2][models.KeyID].(string)
		deleteID := resources[deleteIndex][models.KeyID].(string)

		// read empty source
		_, err = source.Read(ctx)
		if err != nil && err != sdk.ErrBackoffRetry {
			t.Errorf("read: %s", err)
		}

		// add resource
		resources, err = generateResources(ctx, cli, cfg, 1)
		if err != nil {
			t.Errorf("prepare data: %s", err.Error())
		}

		record, err = source.Read(ctx)
		if err != nil {
			t.Errorf("read: %s", err.Error())
		}

		err = compareResult(record, resources[0], opencdc.OperationCreate)
		if err != nil {
			t.Error(err.Error())
		}

		// update resource
		updatedResource1, err := updateDescription(ctx, cli, cfg, updateID1, "new description")
		if err != nil {
			t.Errorf("update resource: %s", err.Error())
		}

		record, err = source.Read(ctx)
		if err != nil {
			t.Errorf("read: %s", err.Error())
		}

		err = compareResult(record, updatedResource1, opencdc.OperationUpdate)
		if err != nil {
			t.Error(err.Error())
		}

		// delete resource
		deletedResource, err := deleteResource(ctx, cli, cfg, deleteID)
		if err != nil {
			t.Errorf("delete resource: %s", err.Error())
		}

		record, err = source.Read(ctx)
		if err != nil {
			t.Errorf("read: %s", err.Error())
		}

		err = compareResult(record, deletedResource, opencdc.OperationDelete)
		if err != nil {
			t.Error(err.Error())
		}

		err = source.Teardown(ctx)
		if err != nil {
			t.Errorf("teardown: %s", err.Error())
		}

		source = NewSource()

		err = source.Configure(ctx, cfg)
		if err != nil {
			t.Errorf("configure: %s", err.Error())
		}

		err = source.Open(ctx, record.Position)
		if err != nil {
			t.Errorf("open: %s", err.Error())
		}

		// update resource
		updatedResource2, err := updateDescription(ctx, cli, cfg, updateID2, "new description after stop")
		if err != nil {
			t.Errorf("update resource: %s", err.Error())
		}

		record, err = source.Read(ctx)
		if err != nil {
			t.Errorf("read: %s", err.Error())
		}

		err = compareResult(record, updatedResource2, opencdc.OperationUpdate)
		if err != nil {
			t.Error(err.Error())
		}

		// generate new clients for more than one page
		genResCount = 23

		resources, err = generateResources(ctx, cli, cfg, genResCount)
		if err != nil {
			t.Errorf("prepare data: %s", err.Error())
		}

		for i := 0; i < genResCount; i++ {
			record, err = source.Read(ctx)
			if err != nil {
				t.Errorf("read: %s", err.Error())
			}

			err = compareResult(record, resources[i], opencdc.OperationCreate)
			if err != nil {
				t.Error(err.Error())
			}
		}

		// read empty source
		_, err = source.Read(ctx)
		if err != nil && err != sdk.ErrBackoffRetry {
			t.Errorf("read: %s", err)
		}

		err = source.Teardown(ctx)
		if err != nil {
			t.Errorf("teardown: %s", err.Error())
		}
	})

	t.Run("teardown", func(t *testing.T) {
		var (
			ctx = context.Background()
		)

		source := NewSource()

		err := source.Teardown(ctx)
		if err != nil {
			t.Errorf("teardown: %s", err.Error())
		}

		cfg, err := prepareDefaultConfig()
		if err != nil {
			t.Log(err)
			t.Skip()
		}

		source = NewSource()

		err = source.Configure(ctx, cfg)
		if err != nil {
			t.Errorf("configure: %s", err.Error())
		}

		err = source.Teardown(ctx)
		if err != nil {
			t.Errorf("teardown: %s", err.Error())
		}
	})

	t.Run("snapshot mode", func(t *testing.T) {
		var (
			ctx = context.Background()

			record opencdc.Record
		)

		defer goleak.VerifyNone(t, goleak.IgnoreTopFunction("github.com/twmb/go-cache/cache.New[...].func1"))

		cfg, err := prepareDefaultConfig()
		if err != nil {
			t.Log(err)
			t.Skip()
		}

		// update snapshot field in the config
		cfg[config.ConfigSnapshot] = "false"

		cli := retryablehttp.NewClient()
		cli.Logger = sdk.Logger(ctx)
		defer cli.HTTPClient.CloseIdleConnections()

		err = isEmpty(ctx, cli, cfg)
		if err != nil {
			t.Errorf("check is empty: %s", err.Error())

			return
		}

		defer func(ctx context.Context, cfg map[string]string) {
			err = clearResources(ctx, cli, cfg)
			if err != nil {
				t.Errorf("clear resources: %s", err.Error())
			}
		}(ctx, cfg)

		_, err = generateResources(ctx, cli, cfg, 1)
		if err != nil {
			t.Errorf("prepare data: %s", err.Error())
		}

		source := NewSource()

		err = source.Configure(ctx, cfg)
		if err != nil {
			t.Errorf("configure: %s", err.Error())
		}

		err = source.Open(ctx, nil)
		if err != nil {
			t.Errorf("open: %s", err.Error())
		}

		_, err = source.Read(ctx)
		if err != sdk.ErrBackoffRetry {
			t.Errorf("read: %s", err)
		}

		resources, err := generateResources(ctx, cli, cfg, 1)
		if err != nil {
			t.Errorf("prepare data: %s", err.Error())
		}

		record, err = source.Read(ctx)
		if err != nil {
			t.Errorf("read: %s", err)
		}

		err = compareResult(record, resources[0], opencdc.OperationCreate)
		if err != nil {
			t.Error(err.Error())
		}

		_, err = source.Read(ctx)
		if err != sdk.ErrBackoffRetry {
			t.Errorf("read: %s", err)
		}

		err = source.Teardown(ctx)
		if err != nil {
			t.Errorf("teardown: %s", err.Error())
		}
	})
}

func prepareDefaultConfig() (map[string]string, error) {
	return prepareConfig("", "")
}

func prepareConfigWithBatchSize(batchSize string) (map[string]string, error) {
	return prepareConfig("", batchSize)
}

func prepareConfig(secretKey, batchSize string) (map[string]string, error) {
	if secretKey == "" {
		secretKey = os.Getenv("STRIPE_SECRET_KEY")
		if secretKey == "" {
			return map[string]string{}, errors.New("STRIPE_SECRET_KEY env var must be set")
		}
	}

	return map[string]string{
		config.ConfigSecretKey:    secretKey,
		config.ConfigResourceName: resourceName,
		config.ConfigBatchSize:    batchSize,
	}, nil
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

func generateResources(ctx context.Context, cli *retryablehttp.Client, cfg map[string]string, count int,
) ([]map[string]interface{}, error) {
	const (
		nameValue        = "client-%s"
		descriptionValue = "info about the %s"
	)

	var (
		resources []map[string]interface{}
		resource  map[string]interface{}
		err       error
	)

	for i := 0; i < count; i++ {
		clientName := fmt.Sprintf(nameValue, uuid.New().String())

		resource, err = addResource(ctx, cli, cfg, map[string]string{
			models.KeyName:        clientName,
			models.KeyDescription: fmt.Sprintf(descriptionValue, clientName),
		})
		if err != nil {
			return nil, fmt.Errorf("add resource: %w", err)
		}

		resources = append(resources, resource)
	}

	time.Sleep(5 * time.Second)

	return resources, nil
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

func updateDescription(ctx context.Context, cli *retryablehttp.Client, cfg map[string]string, id, description string,
) (map[string]interface{}, error) {
	var resource map[string]interface{}

	data, err := makeRequest(ctx, cli, http.MethodPost, id, cfg, map[string]string{
		models.KeyDescription: description,
	})
	if err != nil {
		return nil, fmt.Errorf("make put request: %w", err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if len(resource) == 0 {
		return nil, errors.New("response is empty")
	}

	return resource, nil
}

func deleteResource(ctx context.Context, cli *retryablehttp.Client, cfg map[string]string, id string,
) (map[string]interface{}, error) {
	var resource map[string]interface{}

	data, err := makeRequest(ctx, cli, http.MethodDelete, id, cfg, nil)
	if err != nil {
		return nil, fmt.Errorf("make delete request: %w", err)
	}

	delete(clients, id)

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if len(resource) == 0 {
		return nil, errors.New("response is empty")
	}

	return resource, nil
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

	reqURL.Path += fmt.Sprintf(models.PathFmt, models.ResourcesMap[cfg[config.ConfigResourceName]])

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
	req.Header.Add(models.HeaderAuthKey, fmt.Sprintf(models.HeaderAuthValueFormat, cfg[config.ConfigSecretKey]))

	resp, err := cli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read all response body: %w", err)
	}

	return data, nil
}

func compareResult(record opencdc.Record, resource map[string]interface{}, operation opencdc.Operation) error {
	if !reflect.DeepEqual(record.Key, opencdc.StructuredData{models.KeyID: resource[models.KeyID].(string)}) {
		return fmt.Errorf("key: got = %v, want %v", string(record.Key.Bytes()), resource[models.KeyID].(string))
	}

	if record.Operation != operation {
		return fmt.Errorf("action: got = %v, want %v", record.Operation, operation)
	}

	if operation == opencdc.OperationDelete {
		return nil
	}

	payload, err := json.Marshal(resource)
	if err != nil {
		return fmt.Errorf("marshal payload error = \"%s\"", err)
	}

	if !reflect.DeepEqual(record.Payload.After.Bytes(), payload) {
		return fmt.Errorf("payload: got = %v, want %v", string(record.Payload.After.Bytes()), string(payload))
	}

	return nil
}
