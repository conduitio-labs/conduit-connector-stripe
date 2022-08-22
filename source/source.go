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

	"github.com/conduitio-labs/conduit-connector-stripe/clients/http"
	"github.com/conduitio-labs/conduit-connector-stripe/config"
	"github.com/conduitio-labs/conduit-connector-stripe/source/iterator"
	"github.com/conduitio-labs/conduit-connector-stripe/stripe"
	sdk "github.com/conduitio/conduit-connector-sdk"
)

// An Iterator defines the interface to iterator methods.
type Iterator interface {
	Next() (sdk.Record, error)
}

// A Source represents the source connector.
type Source struct {
	sdk.UnimplementedSource
	cfg      config.Config
	iterator Iterator
	httpCli  http.Client
}

// New initialises a new source.
func New() sdk.Source {
	return sdk.SourceWithMiddleware(new(Source), sdk.DefaultSourceMiddleware()...)
}

// Parameters returns a map of named Parameters that describe how to configure the Source.
func (s *Source) Parameters() map[string]sdk.Parameter {
	return map[string]sdk.Parameter{
		config.SecretKey: {
			Default:     "",
			Required:    true,
			Description: "Stripe secret key.",
		},
		config.ResourceName: {
			Default:     "",
			Required:    true,
			Description: "Stripe resource name.",
		},
		config.BatchSize: {
			Default:     "",
			Required:    false,
			Description: "Number of Stripe objects in the batch.",
		},
	}
}

// Configure parses and stores configurations, returns an error in case of invalid configuration.
func (s *Source) Configure(_ context.Context, cfgRaw map[string]string) (err error) {
	s.cfg, err = config.Parse(cfgRaw)
	if err != nil {
		return err
	}

	return nil
}

// Open parses sdk.Position and initializes a SnapshotIterator iterator.
func (s *Source) Open(ctx context.Context, position sdk.Position) error {
	pos, err := iterator.ParseSDKPosition(position)
	if err != nil {
		return err
	}

	s.httpCli = http.NewClient(ctx)

	s.iterator = iterator.New(stripe.New(s.cfg, s.httpCli), pos)

	return nil
}

// Read returns the next sdk.Record.
func (s *Source) Read(_ context.Context) (sdk.Record, error) {
	record, err := s.iterator.Next()
	if err != nil {
		return sdk.Record{}, err
	}

	return record, nil
}

// Ack does nothing.
func (s *Source) Ack(ctx context.Context, position sdk.Position) error {
	sdk.Logger(ctx).Debug().Str("position", string(position)).Msg("got ack")

	return nil
}

// Teardown closes any connections which were previously connected from previous requests.
func (s *Source) Teardown(ctx context.Context) error {
	sdk.Logger(ctx).Info().Msg("tearing down a stripe source")

	s.httpCli.Close()

	return nil
}
