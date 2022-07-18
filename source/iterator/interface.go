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
	"github.com/conduitio-labs/conduit-connector-stripe/models"
	sdk "github.com/conduitio/conduit-connector-sdk"
)

// An Interface defines the interface to iterator methods.
type Interface interface {
	Next() (sdk.Record, error)
}

// A Stripe defines the interface of methods.
type Stripe interface {
	GetResource(string) (models.ResourceResponse, error)
	GetEvent(createdAt int64, startingAfter, endingBefore string) (models.EventResponse, error)
}
