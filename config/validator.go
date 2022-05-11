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
	"go.uber.org/multierr"

	"github.com/conduitio/conduit-connector-stripe/models"
)

// Validate validates configuration fields.
func (c Config) Validate() error {
	var err error

	validatorInstance := validator.New()

	err = validatorInstance.RegisterValidation("resource_name", c.validateResourceName)
	if err != nil {
		return fmt.Errorf("register resource_name validation: %w", err)
	}

	validationErr := validatorInstance.Struct(c)
	if validationErr != nil {
		if _, ok := validationErr.(*validator.InvalidValidationError); ok {
			return fmt.Errorf("validate config struct: %w", validationErr)
		}

		for _, e := range validationErr.(validator.ValidationErrors) {
			switch e.ActualTag() {
			case "required":
				err = multierr.Append(err, c.RequiredConfigErr(c.configName(e.Field())))
			case "gte", "lte":
				err = multierr.Append(err, c.OutOfRangeConfigErr(c.configName(e.Field())))
			case "resource_name":
				err = multierr.Append(err, c.WrongResourceNameConfigErr(c.configName(e.Field())))
			}
		}
	}

	return err
}

// RequiredConfigErr returns the formatted required field config error.
func (c Config) RequiredConfigErr(name string) error {
	return fmt.Errorf("%q config value must be set", name)
}

// OutOfRangeConfigErr returns the formatted out of range error.
func (c Config) OutOfRangeConfigErr(name string) error {
	return fmt.Errorf("%q is out of range", name)
}

// IntegerTypeConfigErr returns the formatted integer type error.
func (c Config) IntegerTypeConfigErr(name string) error {
	return fmt.Errorf("%q config value must be an integer", name)
}

// PollingPeriodIsNotDurationErr returns the formatted polling period duration error.
func (c Config) PollingPeriodIsNotDurationErr(name string) error {
	return fmt.Errorf("%q config value must be a valid duration", name)
}

// PollingPeriodPositiveErr returns the formatted negative polling period error.
func (c Config) PollingPeriodPositiveErr(name string) error {
	return fmt.Errorf("%q config value must be a positive", name)
}

// WrongResourceNameConfigErr returns the formatted wrong resource name error.
func (c Config) WrongResourceNameConfigErr(name string) error {
	return fmt.Errorf("%q wrong resource name", name)
}

func (c Config) configName(fieldName string) string {
	return map[string]string{
		"SecretKey":          SecretKey,
		"ResourceName":       ResourceName,
		"HTTPClientRetryMax": HTTPClientRetryMax,
		"Limit":              Limit,
		"PollingPeriod":      PollingPeriod,
	}[fieldName]
}

func (c Config) validateResourceName(fl validator.FieldLevel) bool {
	_, ok := models.ResourcesMap[fl.Field().String()]

	return ok
}
