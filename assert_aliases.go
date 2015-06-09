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
)

// AssertEq(e, a) is equivalent to AssertThat(a, oglematchers.Equals(e)).
func (t *T) AssertEq(expected, actual interface{}, errorParts ...interface{}) {
	t.assertThat(
		actual,
		oglematchers.Equals(expected),
		1,
		errorParts)
}

// AssertNe(e, a) is equivalent to
// AssertThat(a, oglematchers.Not(oglematchers.Equals(e))).
func (t *T) AssertNe(expected, actual interface{}, errorParts ...interface{}) {
	t.assertThat(
		actual,
		oglematchers.Not(oglematchers.Equals(expected)),
		1,
		errorParts)
}

// AssertLt(x, y) is equivalent to AssertThat(x, oglematchers.LessThan(y)).
func (t *T) AssertLt(x, y interface{}, errorParts ...interface{}) {
	t.assertThat(x, oglematchers.LessThan(y), 1, errorParts)
}

// AssertLe(x, y) is equivalent to AssertThat(x, oglematchers.LessOrEqual(y)).
func (t *T) AssertLe(x, y interface{}, errorParts ...interface{}) {
	t.assertThat(x, oglematchers.LessOrEqual(y), 1, errorParts)
}

// AssertGt(x, y) is equivalent to AssertThat(x, oglematchers.GreaterThan(y)).
func (t *T) AssertGt(x, y interface{}, errorParts ...interface{}) {
	t.assertThat(x, oglematchers.GreaterThan(y), 1, errorParts)
}

// AssertGe(x, y) is equivalent to
// AssertThat(x, oglematchers.GreaterOrEqual(y)).
func (t *T) AssertGe(x, y interface{}, errorParts ...interface{}) {
	t.assertThat(x, oglematchers.GreaterOrEqual(y), 1, errorParts)
}

// AssertTrue(b) is equivalent to AssertThat(b, oglematchers.Equals(true)).
func (t *T) AssertTrue(b interface{}, errorParts ...interface{}) {
	t.assertThat(b, oglematchers.Equals(true), 1, errorParts)
}

// AssertFalse(b) is equivalent to AssertThat(b, oglematchers.Equals(false)).
func (t *T) AssertFalse(b interface{}, errorParts ...interface{}) {
	t.assertThat(b, oglematchers.Equals(false), 1, errorParts)
}
