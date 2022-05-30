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
	"io/ioutil"
	"net/url"
	"os"
	"reflect"
	"testing"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/conduitio/conduit-connector-stripe/config"
	"github.com/conduitio/conduit-connector-stripe/models"
	"github.com/hashicorp/go-retryablehttp"
)

const (
	resourceName = "customer"

	idKey          = "id"
	nameKey        = "name"
	descriptionKey = "description"

	methodGet    = "GET"
	methodPost   = "POST"
	methodDelete = "DELETE"

	backoffRetryErr = "backoff retry"
)

var resources []map[string]interface{}

func TestSource_Read(t *testing.T) { // nolint:gocyclo,nolintlint
	// add resCount resources
	// initialize source
	// read resCount - stop resources
	// reinitialize source
	// read rest resources
	// read when there is no resources
	t.Run("snapshot iterator", func(t *testing.T) {
		const (
			resCount = 12
			stop     = 5
		)

		var (
			ctx = context.Background()

			record sdk.Record
			pos    sdk.Position
		)

		cfg, err := prepareConfig()
		if err != nil {
			t.Log(err)
			t.Skip()
		}

		err = isEmpty(ctx, cfg)
		if err != nil {
			t.Errorf("check is empty: %s", err.Error())

			return
		}

		defer func(ctx context.Context, cfg map[string]string) {
			for i := len(resources) - 1; i >= 0; i-- {
				err = deleteResource(ctx, cfg, i)
				if err != nil {
					t.Errorf("delete resource: %s", err.Error())
				}
			}
		}(ctx, cfg)

		// add resCount resources
		err = prepareData(ctx, cfg, resCount)
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

		for i := resCount - 1; i >= stop; i-- {
			record, err = source.Read(ctx)
			if err != nil {
				t.Errorf("read: %s", err.Error())
			}

			pos = record.Position

			err = compareResult(record, resources[i], models.InsertAction)
			if err != nil {
				t.Errorf(err.Error())
			}
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

		for i := stop - 1; i >= 0; i-- {
			record, err = source.Read(ctx)
			if err != nil {
				t.Errorf("read: %s", err.Error())
			}

			err = compareResult(record, resources[i], models.InsertAction)
			if err != nil {
				t.Errorf(err.Error())
			}
		}

		// read empty source
		_, err = source.Read(ctx)
		if err != nil && err.Error() != backoffRetryErr {
			t.Errorf("read: %s", err.Error())
		}
	})

	// add resCount resources
	// initialize source
	// read all resources by snapshot iterator
	// reinitialize source
	// add resource
	// read resource by cdc iterator
	// read when there is no events
	t.Run("reinitialize source between snapshot and cdc iterators", func(t *testing.T) {
		const resCount = 3

		var (
			ctx = context.Background()

			record sdk.Record
			rp     sdk.Position
		)

		cfg, err := prepareConfig()
		if err != nil {
			t.Log(err)
			t.Skip()
		}

		err = isEmpty(ctx, cfg)
		if err != nil {
			t.Errorf("check is empty: %s", err.Error())

			return
		}

		defer func(ctx context.Context, cfg map[string]string) {
			for i := len(resources) - 1; i >= 0; i-- {
				err = deleteResource(ctx, cfg, i)
				if err != nil {
					t.Errorf("delete resource: %s", err.Error())
				}
			}
		}(ctx, cfg)

		// add resCount resources
		err = prepareData(ctx, cfg, resCount)
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

		for i := resCount - 1; i >= 0; i-- {
			record, err = source.Read(ctx)
			if err != nil {
				t.Errorf("read: %s", err.Error())
			}

			rp = record.Position

			err = compareResult(record, resources[i], models.InsertAction)
			if err != nil {
				t.Errorf(err.Error())
			}
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
		err = addResource(ctx, cfg, map[string]string{
			nameKey:        "inserted",
			descriptionKey: "inserted description",
		})
		if err != nil {
			t.Errorf("add resource: %s", err.Error())
		}

		record, err = source.Read(ctx)
		if err != nil {
			t.Errorf("read: %s", err.Error())
		}

		err = compareResult(record, resources[len(resources)-1], models.InsertAction)
		if err != nil {
			t.Errorf(err.Error())
		}

		// read empty source
		_, err = source.Read(ctx)
		if err != nil && err.Error() != backoffRetryErr {
			t.Errorf("read: %s", err.Error())
		}
	})

	// add resCount resources
	// initialize source
	// read all resources by snapshot iterator
	// read when there is no resources
	// add resource
	// update resource
	// delete resource
	// reinitialize source
	// update resource
	t.Run("reinitialize cdc", func(t *testing.T) {
		const resCount = 3

		var (
			ctx = context.Background()

			record sdk.Record
		)

		cfg, err := prepareConfig()
		if err != nil {
			t.Log(err)
			t.Skip()
		}

		err = isEmpty(ctx, cfg)
		if err != nil {
			t.Errorf("check is empty: %s", err.Error())

			return
		}

		defer func(ctx context.Context, cfg map[string]string) {
			for i := len(resources) - 1; i >= 0; i-- {
				err = deleteResource(ctx, cfg, i)
				if err != nil {
					t.Errorf("delete resource: %s", err.Error())
				}
			}
		}(ctx, cfg)

		// add resCount resources
		err = prepareData(ctx, cfg, resCount)
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

		for i := resCount - 1; i >= 0; i-- {
			record, err = source.Read(ctx)
			if err != nil {
				t.Errorf("read: %s", err.Error())
			}

			err = compareResult(record, resources[i], models.InsertAction)
			if err != nil {
				t.Errorf(err.Error())
			}
		}

		// read empty source
		_, err = source.Read(ctx)
		if err != nil && err.Error() != backoffRetryErr {
			t.Errorf("read: %s", err.Error())
		}

		// add resource
		err = addResource(ctx, cfg, map[string]string{
			nameKey:        "inserted",
			descriptionKey: "inserted description",
		})
		if err != nil {
			t.Errorf("add resource: %s", err.Error())
		}

		record, err = source.Read(ctx)
		if err != nil {
			t.Errorf("read: %s", err.Error())
		}

		err = compareResult(record, resources[len(resources)-1], models.InsertAction)
		if err != nil {
			t.Errorf(err.Error())
		}

		// update resource
		updateIndex := 1
		err = updateResource(ctx, cfg, updateIndex, "new name", "new description")
		if err != nil {
			t.Errorf("update resource: %s", err.Error())
		}

		record, err = source.Read(ctx)
		if err != nil {
			t.Errorf("read: %s", err.Error())
		}

		err = compareResult(record, resources[updateIndex], models.UpdateAction)
		if err != nil {
			t.Errorf(err.Error())
		}

		// delete resource
		deleteIndex := 2
		deletedResource := resources[deleteIndex]

		err = deleteResource(ctx, cfg, deleteIndex)
		if err != nil {
			t.Errorf("delete resource: %s", err.Error())
		}

		record, err = source.Read(ctx)
		if err != nil {
			t.Errorf("read: %s", err.Error())
		}

		err = compareResult(record, deletedResource, models.DeleteAction)
		if err != nil {
			t.Errorf(err.Error())
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
		updateIndex = 0
		err = updateResource(ctx, cfg, updateIndex, "new name after stop", "new description after stop")
		if err != nil {
			t.Errorf("update resource: %s", err.Error())
		}

		record, err = source.Read(ctx)
		if err != nil {
			t.Errorf("read: %s", err.Error())
		}

		err = compareResult(record, resources[updateIndex], models.UpdateAction)
		if err != nil {
			t.Errorf(err.Error())
		}
	})
}

func prepareConfig() (map[string]string, error) {
	secretKey := os.Getenv("STRIPE_SECRET_KEY")
	if secretKey == "" {
		return map[string]string{}, errors.New("STRIPE_SECRET_KEY env var must be set")
	}

	return map[string]string{
		config.SecretKey:    secretKey,
		config.ResourceName: resourceName,
	}, nil
}

func isEmpty(ctx context.Context, cfg map[string]string) error {
	var resource models.ResourceResponse

	data, err := makeRequest(ctx, methodGet, "", cfg, nil)
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

func prepareData(ctx context.Context, cfg map[string]string, count int) error {
	const (
		nameValue        = "client_%d"
		descriptionValue = "info about the client_%d"
	)

	var err error

	for i := 0; i < count; i++ {
		err = addResource(ctx, cfg, map[string]string{
			nameKey:        fmt.Sprintf(nameValue, i),
			descriptionKey: fmt.Sprintf(descriptionValue, i),
		})
		if err != nil {
			return fmt.Errorf("add resource: %w", err)
		}
	}

	return nil
}

func addResource(ctx context.Context, cfg, params map[string]string) error {
	var resource map[string]interface{}

	data, err := makeRequest(ctx, methodPost, "", cfg, params)
	if err != nil {
		return fmt.Errorf("make post request: %w", err)
	}

	err = json.Unmarshal(data, &resource)
	if err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	if len(resource) == 0 {
		return errors.New("response is empty")
	}

	resources = append(resources, resource)

	return nil
}

func updateResource(ctx context.Context, cfg map[string]string, index int, name, description string) error {
	_, err := makeRequest(ctx, methodPost, resources[index]["id"].(string), cfg, map[string]string{
		nameKey:        name,
		descriptionKey: description,
	})
	if err != nil {
		return fmt.Errorf("make put request: %w", err)
	}

	resources[index]["name"] = name
	resources[index]["description"] = description

	return nil
}

func deleteResource(ctx context.Context, cfg map[string]string, index int) error {
	id := resources[index]["id"].(string)

	_, err := makeRequest(ctx, methodDelete, id, cfg, nil)
	if err != nil {
		return fmt.Errorf("make delete request: %w", err)
	}

	resources = append(resources[:index], resources[index+1:]...)

	return nil
}

func makeRequest(ctx context.Context, method, path string, cfg, params map[string]string) ([]byte, error) {
	const (
		apiURL  = "https://api.stripe.com/v1"
		pathFmt = "/%s"

		headerAuthKey         = "Authorization"
		headerAuthValueFormat = "Bearer %s"
	)

	reqURL, err := url.Parse(apiURL)
	if err != nil {
		return nil, fmt.Errorf("parse api url: %w", err)
	}

	reqURL.Path += fmt.Sprintf(pathFmt, models.ResourcesMap[cfg[config.ResourceName]])

	if path != "" {
		reqURL.Path += fmt.Sprintf(pathFmt, path)
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
	req.Header.Add(headerAuthKey, fmt.Sprintf(headerAuthValueFormat, cfg[config.SecretKey]))

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

func compareResult(record sdk.Record, resource map[string]interface{}, action string) error {
	payload, err := json.Marshal(resource)
	if err != nil {
		return fmt.Errorf("marshal payload error = \"%s\"", err)
	}

	if !reflect.DeepEqual(record.Payload.Bytes(), payload) {
		return fmt.Errorf("payload: got = %v, want %v", string(record.Payload.Bytes()), string(payload))
	}

	if !reflect.DeepEqual(record.Key, sdk.StructuredData{idKey: resource["id"].(string)}) {
		return fmt.Errorf("key: got = %v, want %v", string(record.Key.Bytes()), resource["id"].(string))
	}

	if record.Metadata[models.ActionKey] != action {
		return fmt.Errorf("action: got = %v, want %v", record.Metadata[models.ActionKey], action)
	}

	return nil
}
