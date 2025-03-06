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
	"fmt"
	"testing"

	"github.com/conduitio-labs/conduit-connector-stripe/models/resources"
	"github.com/matryer/is"
)

const (
	testSecretKey = "sk_test_123456789"
)

func TestValidateConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      *Config
		wantErr error
	}{
		{
			name: "success_valid_config",
			in: &Config{
				SecretKey:    testSecretKey,
				ResourceName: resources.CreditNoteResource,
				BatchSize:    10,
				Snapshot:     true,
			},
			wantErr: nil,
		},
		{
			name: "failure_invalid_resource_name",
			in: &Config{
				SecretKey:    testSecretKey,
				ResourceName: "invalid_resource",
				BatchSize:    10,
				Snapshot:     true,
			},
			wantErr: fmt.Errorf("\"invalid_resource\" wrong resource name"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)

			err := tt.in.Validate()
			if tt.wantErr == nil {
				is.NoErr(err)
			} else {
				is.True(err != nil)
				is.Equal(err.Error(), tt.wantErr.Error())
			}
		})
	}
}
