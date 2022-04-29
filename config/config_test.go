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
	"reflect"
	"testing"

	"go.uber.org/multierr"
)

func TestParse(t *testing.T) {
	underTestConfig := Config{}

	tests := []struct {
		name        string
		in          map[string]string
		want        Config
		wantErr     bool
		expectedErr string
	}{
		{
			name: "Valid config",
			in: map[string]string{
				SecretKey:          "sk_51JB",
				ResourceName:       "subscriptions",
				HTTPClientRetryMax: "5",
			},
			want: Config{
				SecretKey:          "sk_51JB",
				ResourceName:       "subscriptions",
				HTTPClientRetryMax: 5,
			},
		},
		{
			name: "HTTPClientRetryMax by default",
			in: map[string]string{
				SecretKey:    "sk_51JB",
				ResourceName: "subscriptions",
			},
			want: Config{
				SecretKey:          "sk_51JB",
				ResourceName:       "subscriptions",
				HTTPClientRetryMax: RetryMaxDefault,
			},
		},
		{
			name: "No secret key",
			in: map[string]string{
				SecretKey:    "",
				ResourceName: "subscriptions",
			},
			wantErr:     true,
			expectedErr: underTestConfig.RequiredConfigErr(SecretKey).Error(),
		},
		{
			name: "Empty secret key",
			in: map[string]string{
				SecretKey:    "",
				ResourceName: "subscriptions",
			},
			wantErr:     true,
			expectedErr: underTestConfig.RequiredConfigErr(SecretKey).Error(),
		},
		{
			name: "No resource name",
			in: map[string]string{
				SecretKey: "sk_51JB",
			},
			wantErr:     true,
			expectedErr: underTestConfig.RequiredConfigErr(ResourceName).Error(),
		},
		{
			name: "Empty resource name",
			in: map[string]string{
				SecretKey:    "sk_51JB",
				ResourceName: "",
			},
			wantErr:     true,
			expectedErr: underTestConfig.RequiredConfigErr(ResourceName).Error(),
		},
		{
			name: "No secret key and resource name",
			in: map[string]string{
				SecretKey:    "",
				ResourceName: "",
			},
			wantErr: true,
			expectedErr: multierr.Combine(underTestConfig.RequiredConfigErr(SecretKey),
				underTestConfig.RequiredConfigErr(ResourceName)).Error(),
		},
		{
			name: "HTTPClientRetryMax is greater than the value of lte tag",
			in: map[string]string{
				SecretKey:          "sk_51JB",
				ResourceName:       "subscriptions",
				HTTPClientRetryMax: "12",
			},
			wantErr:     true,
			expectedErr: underTestConfig.OutOfRangeConfigErr(HTTPClientRetryMax).Error(),
		},
		{
			name: "HTTPClientRetryMax is more than the value of lte tag",
			in: map[string]string{
				SecretKey:          "sk_51JB",
				ResourceName:       "subscriptions",
				HTTPClientRetryMax: "0",
			},
			wantErr:     true,
			expectedErr: underTestConfig.OutOfRangeConfigErr(HTTPClientRetryMax).Error(),
		},
		{
			name: "Invalid HTTPClientRetryMax",
			in: map[string]string{
				SecretKey:          "sk_51JB",
				ResourceName:       "subscriptions",
				HTTPClientRetryMax: "test",
			},
			wantErr:     true,
			expectedErr: underTestConfig.IntegerTypeConfigErr(HTTPClientRetryMax).Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.in)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("parse error = \"%s\", wantErr %t", err.Error(), tt.wantErr)

					return
				}

				if err.Error() != tt.expectedErr {
					t.Errorf("expected error \"%s\", got \"%s\"", tt.expectedErr, err.Error())

					return
				}

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse = %v, want %v", got, tt.want)
			}
		})
	}
}
