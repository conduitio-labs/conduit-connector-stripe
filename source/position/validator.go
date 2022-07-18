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

package position

import (
	"errors"
	"fmt"

	"github.com/conduitio-labs/conduit-connector-stripe/validator"
	v "github.com/go-playground/validator/v10"
	"go.uber.org/multierr"
)

// Validate validates position fields.
func (p Position) Validate() error {
	var err error

	validationErr := validator.Get().Struct(p)
	if validationErr != nil {
		if _, ok := validationErr.(*v.InvalidValidationError); ok {
			return fmt.Errorf("validate config struct: %w", validationErr)
		}

		for _, e := range validationErr.(v.ValidationErrors) {
			switch e.ActualTag() {
			case "required":
				err = multierr.Append(err, validator.RequiredErr(e.Field()))
			case "iterator_type":
				err = multierr.Append(err, errors.New(validator.UnexpectedIteratorTypeErr))
			}
		}
	}

	return err
}
