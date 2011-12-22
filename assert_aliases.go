// Copyright 2011 Aaron Jacobs. All Rights Reserved.
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
	"github.com/jacobsa/oglematchers"
	"path"
	"runtime"
)

// AssertEq(e, a) is equivalent to AssertThat(a, oglematchers.Equals(e)).
func AssertEq(expected, actual interface{}) AssertationResult {
	res := AssertThat(actual, oglematchers.Equals(expected))
	res.SetCaller(getCallerForAlias())
	return res
}

// AssertNe(e, a) is equivalent to AssertThat(a, oglematchers.Not(oglematchers.Equals(e))).
func AssertNe(expected, actual interface{}) AssertationResult {
	res := AssertThat(actual, oglematchers.Not(oglematchers.Equals(expected)))
	res.SetCaller(getCallerForAlias())
	return res
}

// AssertLt(x, y) is equivalent to AssertThat(x, oglematchers.LessThan(y)).
func AssertLt(x, y interface{}) AssertationResult {
	res := AssertThat(x, oglematchers.LessThan(y))
	res.SetCaller(getCallerForAlias())
	return res
}

// AssertLe(x, y) is equivalent to AssertThat(x, oglematchers.LessOrEqual(y)).
func AssertLe(x, y interface{}) AssertationResult {
	res := AssertThat(x, oglematchers.LessOrEqual(y))
	res.SetCaller(getCallerForAlias())
	return res
}

// AssertGt(x, y) is equivalent to AssertThat(x, oglematchers.GreaterThan(y)).
func AssertGt(x, y interface{}) AssertationResult {
	res := AssertThat(x, oglematchers.GreaterThan(y))
	res.SetCaller(getCallerForAlias())
	return res
}

// AssertGe(x, y) is equivalent to AssertThat(x, oglematchers.GreaterOrEqual(y)).
func AssertGe(x, y interface{}) AssertationResult {
	res := AssertThat(x, oglematchers.GreaterOrEqual(y))
	res.SetCaller(getCallerForAlias())
	return res
}

// AssertTrue(b) is equivalent to AssertThat(b, oglematchers.Equals(true)).
func AssertTrue(b interface{}) AssertationResult {
	res := AssertThat(b, oglematchers.Equals(true))
	res.SetCaller(getCallerForAlias())
	return res
}

// AssertFalse(b) is equivalent to AssertThat(b, oglematchers.Equals(false)).
func AssertFalse(b interface{}) AssertationResult {
	res := AssertThat(b, oglematchers.Equals(false))
	res.SetCaller(getCallerForAlias())
	return res
}
