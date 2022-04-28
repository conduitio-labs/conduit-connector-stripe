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

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/multierr"
)

func TestParse(t *testing.T) {
	underTestConfig := Config{}

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
				SecretKey:          "sk_51JB",
				ResourceName:       "subscriptions",
				HTTPClientRetryMax: "5",
			},
			expectedSecretKey:    "sk_51JB",
			expectedResourceName: "subscriptions",
			expectedRetryMax:     5,
		},
		{
			name: "No secret key",
			cfg: map[string]string{
				SecretKey:    "",
				ResourceName: "subscriptions",
			},
			expectedErr: underTestConfig.RequiredConfigErr(SecretKey),
		},
		{
			name: "Empty secret key",
			cfg: map[string]string{
				SecretKey:    "",
				ResourceName: "subscriptions",
			},
			expectedErr: underTestConfig.RequiredConfigErr(SecretKey),
		},
		{
			name: "No resource name",
			cfg: map[string]string{
				SecretKey: "sk_51JB",
			},
			expectedErr: underTestConfig.RequiredConfigErr(ResourceName),
		},
		{
			name: "Empty resource name",
			cfg: map[string]string{
				SecretKey:    "sk_51JB",
				ResourceName: "",
			},
			expectedErr: underTestConfig.RequiredConfigErr(ResourceName),
		},
		{
			name: "No secret key and resource name",
			cfg: map[string]string{
				SecretKey:    "",
				ResourceName: "",
			},
			expectedErr: multierr.Combine(underTestConfig.RequiredConfigErr(SecretKey),
				underTestConfig.RequiredConfigErr(ResourceName)),
		},
		{
			name: "HTTPClientRetryMax is greater than the value of lte tag",
			cfg: map[string]string{
				SecretKey:          "sk_51JB",
				ResourceName:       "subscriptions",
				HTTPClientRetryMax: "12",
			},
			expectedErr: underTestConfig.OutOfRangeConfigErr(HTTPClientRetryMax),
		},
		{
			name: "HTTPClientRetryMax is more than the value of lte tag",
			cfg: map[string]string{
				SecretKey:          "sk_51JB",
				ResourceName:       "subscriptions",
				HTTPClientRetryMax: "0",
			},
			expectedErr: underTestConfig.OutOfRangeConfigErr(HTTPClientRetryMax),
		},
		{
			name: "Invalid HTTPClientRetryMax",
			cfg: map[string]string{
				SecretKey:          "sk_51JB",
				ResourceName:       "subscriptions",
				HTTPClientRetryMax: "test",
			},
			expectedErr: underTestConfig.IntegerTypeConfigErr(HTTPClientRetryMax),
		},
		{
			name: "HTTPClientRetryMax by default",
			cfg: map[string]string{
				SecretKey:    "sk_51JB",
				ResourceName: "subscriptions",
			},
			expectedSecretKey:    "sk_51JB",
			expectedResourceName: "subscriptions",
			expectedRetryMax:     RetryMaxDefault,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := Parse(tt.cfg)
			if err != nil {
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.Equal(t, tt.expectedSecretKey, cfg.SecretKey)
				assert.Equal(t, tt.expectedResourceName, cfg.ResourceName)
				assert.Equal(t, tt.expectedRetryMax, cfg.HTTPClientRetryMax)
			}
		})
	}
}
