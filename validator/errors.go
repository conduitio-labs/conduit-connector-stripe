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
	"fmt"
)

// An UnexpectedIteratorTypeErr represents the message of unexpected iterator type error.
const UnexpectedIteratorTypeErr = "unexpected iterator type"

// WrongResourceNameErr returns the formatted wrong resource name error.
func WrongResourceNameErr(name string) error {
	return fmt.Errorf("%q wrong resource name", name)
}

// RequiredErr returns the formatted required field error.
func RequiredErr(name string) error {
	return fmt.Errorf("%q value must be set", name)
}

// IntegerTypeConfigErr returns the formatted integer type error.
func IntegerTypeConfigErr(name string) error {
	return fmt.Errorf("%q config value must be an integer", name)
}

// InvalidBatchSizeErr returns the formatted invalid batch size value error.
func InvalidBatchSizeErr(name string) error {
	return fmt.Errorf("%q must be greater than or equal to 1 and less than or equal to 100", name)
}
