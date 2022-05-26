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
	"errors"
	"fmt"
)

// WrongResourceNameErr returns the formatted wrong resource name error.
func WrongResourceNameErr(name string) error {
	return fmt.Errorf("%q wrong resource name", name)
}

// RequiredErr returns the formatted required field error.
func RequiredErr(name string) error {
	return fmt.Errorf("%q value must be set", name)
}

// UnexpectedIteratorTypeErr returns the unexpected iterator type error.
func UnexpectedIteratorTypeErr() error {
	return errors.New("unexpected iterator type")
}
