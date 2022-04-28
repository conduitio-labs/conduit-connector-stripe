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
	"testing"

	"github.com/ConduitIO/conduit-connector-stripe/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"
)

func TestSource_Configure(t *testing.T) {
	underTestConfig := config.Config{}
	underTestSource := Source{}

	tests := []struct {
		name                 string
		cfg                  map[string]string
		expectedSecretKey    string
		expectedResourceName string
		expectedRetryMax     int
		expectedErr          error
	}{
		{
			name: "Valid config",
			cfg: map[string]string{
				config.SecretKey:    "sk_51JB",
				config.ResourceName: "subscriptions",
			},
			expectedSecretKey:    "sk_51JB",
			expectedResourceName: "subscriptions",
			expectedRetryMax:     config.RetryMaxDefault,
		},
		{
			name: "No secret key and resource name",
			cfg: map[string]string{
				config.SecretKey:    "",
				config.ResourceName: "",
			},
			expectedErr: multierr.Combine(underTestConfig.RequiredConfigErr(config.SecretKey),
				underTestConfig.RequiredConfigErr(config.ResourceName)),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := underTestSource.Configure(context.Background(), tt.cfg)
			if err != nil {
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, underTestSource.config)
				assert.NotNil(t, underTestSource.httpClient)
			}
		})
	}
}
