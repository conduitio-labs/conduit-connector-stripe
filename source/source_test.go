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
	"reflect"
	"testing"

	"github.com/conduitio-labs/conduit-connector-stripe/config"
)

func TestSource_Configure(t *testing.T) {
	source := new(Source)

	tests := []struct {
		name        string
		in          map[string]string
		want        Source
		wantErr     bool
		expectedErr string
	}{
		{
			name: "valid config",
			in: map[string]string{
				config.ConfigSecretKey:    "sk_51JB",
				config.ConfigResourceName: "subscription",
			},
			want: Source{
				cfg: config.Config{
					SecretKey:    "sk_51JB",
					ResourceName: "subscription",
					Snapshot:     true,
					BatchSize:    10,
				},
			},
		},
		{
			name: "no secret key",
			in: map[string]string{
				config.ConfigSecretKey:    "",
				config.ConfigResourceName: "subscription",
			},
			wantErr:     true,
			expectedErr: `config invalid: error validating "secretKey": required parameter is not provided`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := source.Configure(context.Background(), tt.in)
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

			if !reflect.DeepEqual(source.cfg, tt.want.cfg) {
				t.Errorf("parse = %v, want %v", source.cfg, tt.want.cfg)

				return
			}
		})
	}
}
