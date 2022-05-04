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

package iterator

import (
	"encoding/json"
	"testing"
)

func TestNext(t *testing.T) {
	underTestSnapshot := SnapshotIterator{}

	type wantData struct {
		position string
		action   string
		key      string
		payload  string
	}

	tests := []struct {
		name        string
		in          []byte
		len         int
		want        []wantData
		wantErr     bool
		expectedErr string
	}{
		{
			name: "valid data",
			in: []byte(`{
    "data": [
        {
            "id": "prod_La50",
            "created": 1651153850
        },
		{
            "id": "prod_La49",
            "created": 1651153849
        },
		{
            "id": "prod_La48",
            "created": 1651153848
        }
    ],
    "has_more": false
}`),
			len: 3,
			want: []wantData{
				{
					position: "false.",
					action:   "insert",
					key:      "prod_La50",
					payload:  `{"created":1651153850,"id":"prod_La50"}`,
				},
				{
					position: "false.prod_La50",
					action:   "insert",
					key:      "prod_La49",
					payload:  `{"created":1651153849,"id":"prod_La49"}`,
				},
				{
					position: "false.prod_La49",
					action:   "insert",
					key:      "prod_La48",
					payload:  `{"created":1651153848,"id":"prod_La48"}`,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			underTestSnapshot.index = 0

			if err := json.Unmarshal(tt.in, &underTestSnapshot.response); err != nil {
				t.Fatalf("%s: %v", tt.in, err)
			}

			for i := 0; i < tt.len; i++ {
				got, err := underTestSnapshot.Next()
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

				if string(got.Position) != tt.want[i].position {
					t.Errorf("position: got \"%s\", want \"%s\"", got.Position, tt.want[i].position)

					return
				}

				if got.Metadata["action"] != tt.want[i].action {
					t.Errorf("action: got \"%s\", want \"%s\"", got.Metadata["action"], tt.want[i].action)

					return
				}

				if string(got.Key.Bytes()) != tt.want[i].key {
					t.Errorf("key: got \"%s\", want \"%s\"", got.Key.Bytes(), tt.want[i].key)

					return
				}

				if string(got.Payload.Bytes()) != tt.want[i].payload {
					t.Errorf("payload: got \"%s\", want \"%s\"", got.Payload.Bytes(), tt.want[i].payload)

					return
				}
			}
		})
	}
}
