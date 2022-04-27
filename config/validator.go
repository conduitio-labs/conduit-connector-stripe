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

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

// Validate validates configuration fields.
func (c Config) Validate() error {
	var err error

	validate := validator.New()

	validationErr := validate.Struct(c)
	if validationErr != nil {
		if _, ok := validationErr.(*validator.InvalidValidationError); ok {
			return errors.Wrap(validationErr, "validate config struct")
		}

		for _, e := range validationErr.(validator.ValidationErrors) {
			if e.ActualTag() == "required" {
				err = multierr.Append(err, requiredConfigErr(configName(e.Field())))
			}
		}
	}

	return err
}

func requiredConfigErr(name string) error {
	return fmt.Errorf("%q config value must be set", name)
}

func configName(fieldName string) string {
	return map[string]string{
		"SecretKey":    SecretKey,
		"ResourceName": ResourceName,
	}[fieldName]
}
