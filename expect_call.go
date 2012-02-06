// Copyright 2012 Aaron Jacobs. All Rights Reserved.
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
)

// ExpectCall expresses an expectation that the method of the given name
// should be called on the supplied mock object. It returns a function that
// should be called with the expected arguments, matchers for the arguments,
// or a mix of both.
//
// For example:
//
//     mockWriter := [...]
//     ogletest.ExpectCall(mockWriter, "Write")(ElementsAre(0x1))
//         .WillOnce(oglemock.Return(1, nil))
//
// This is a shortcut for calling i.MockController.ExpectCall, where i is the
// TestInfo struct for the currently-running test. Unlike that direct approach,
// this function automatically sets the correct file name and line number.
func ExpectCall(o oglemock.MockObject, methodName string) oglemock.PartialExpecation {
	// TODO
	return nil
}
