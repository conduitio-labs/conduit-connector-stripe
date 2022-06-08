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

package validator

import (
	"sync"

	"github.com/conduitio/conduit-connector-stripe/models"
	v "github.com/go-playground/validator/v10"
)

var (
	validatorInstance *v.Validate

	once sync.Once
)

// Get initializes and registers validation tags once,
// and returns validator instance.
func Get() *v.Validate {
	once.Do(func() {
		validatorInstance = v.New()

		err := validatorInstance.RegisterValidation("resource_name", validateResourceName)
		if err != nil {
			return
		}

		err = validatorInstance.RegisterValidation("iterator_type", validateIteratorType)
		if err != nil {
			return
		}
	})

	return validatorInstance
}

func validateResourceName(fl v.FieldLevel) bool {
	_, ok := models.ResourcesMap[fl.Field().String()]

	return ok
}

func validateIteratorType(fl v.FieldLevel) bool {
	_, ok := models.IteratorTypeMap[models.IteratorType(fl.Field().String())]

	return ok
}
