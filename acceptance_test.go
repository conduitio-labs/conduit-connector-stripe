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

package stripe

import (
	"os"
	"testing"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/conduitio/conduit-connector-stripe/config"
	"github.com/conduitio/conduit-connector-stripe/source"
	"go.uber.org/goleak"
)

func TestAcceptance(t *testing.T) {
	secretKey := os.Getenv("STRIPE_SECRET_KEY")
	if secretKey == "" {
		t.Skip("STRIPE_SECRET_KEY env var must be set")
	}

	resourceName := os.Getenv("STRIPE_RESOURCE_NAME")
	if resourceName == "" {
		t.Skip("STRIPE_RESOURCE_NAME env var must be set")
	}

	cfg := map[string]string{
		config.SecretKey:    secretKey,
		config.ResourceName: resourceName,
	}

	sdk.AcceptanceTest(t, sdk.ConfigurableAcceptanceTestDriver{
		Config: sdk.ConfigurableAcceptanceTestDriverConfig{
			Connector: sdk.Connector{
				NewSpecification: Specification,
				NewSource:        source.NewSource,
				NewDestination:   nil,
			},
			SourceConfig:      cfg,
			DestinationConfig: nil,
			Skip: []string{
				// the method requires NewDestination.
				"TestSource_Read*",
				// the method requires NewDestination.
				"TestSource_Open_ResumeAtPositionSnapshot",
				// the method requires NewDestination.
				"TestSource_Open_ResumeAtPositionCDC",
			},
			GoleakOptions: []goleak.Option{
				goleak.IgnoreTopFunction("database/sql.(*DB).connectionOpener"),
				goleak.IgnoreTopFunction("net/http.(*persistConn).writeLoop"),
				goleak.IgnoreTopFunction("internal/poll.runtime_pollWait"),
			},
		},
	},
	)
}
