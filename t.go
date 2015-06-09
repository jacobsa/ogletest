// Copyright 2015 Aaron Jacobs. All Rights Reserved.
// Author: aaronjjacobs@gmail.com (Aaron Jacobs)
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

package ogletest

import (
	"github.com/jacobsa/oglemock"
	"golang.org/x/net/context"
)

// T is a type that contains dependencies for test functions and whose methods
// can be used by test functions to control test execution, including adding
// failure messages.
type T struct {
	// A context that the test should use for any context-aware code it calls.
	// May carry tracing information or be used for cancellation.
	Ctx context.Context

	// A mock controller that can be used for creating mocks for use by the test.
	// Mock errors will be associated with the test. The Finish method should not
	// be run by the user; ogletest will do that automatically after the test
	// finishes.
	MockController oglemock.Controller
}
