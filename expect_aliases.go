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

import "github.com/jacobsa/oglematchers"

// ExpectEq(e, a) is equivalent to ExpectThat(a, oglematchers.Equals(e)).
func (t *T) ExpectEq(expected, actual interface{}, errorParts ...interface{}) {
	t.expectThat(actual, oglematchers.Equals(expected), 1, errorParts)
}

// ExpectNe(e, a) is equivalent to
// ExpectThat(a, oglematchers.Not(oglematchers.Equals(e))).
func (t *T) ExpectNe(expected, actual interface{}, errorParts ...interface{}) {
	t.expectThat(
		actual,
		oglematchers.Not(oglematchers.Equals(expected)),
		1,
		errorParts)
}

// ExpectLt(x, y) is equivalent to ExpectThat(x, oglematchers.LessThan(y)).
func (t *T) ExpectLt(x, y interface{}, errorParts ...interface{}) {
	t.expectThat(x, oglematchers.LessThan(y), 1, errorParts)
}

// ExpectLe(x, y) is equivalent to ExpectThat(x, oglematchers.LessOrEqual(y)).
func (t *T) ExpectLe(x, y interface{}, errorParts ...interface{}) {
	t.expectThat(x, oglematchers.LessOrEqual(y), 1, errorParts)
}

// ExpectGt(x, y) is equivalent to ExpectThat(x, oglematchers.GreaterThan(y)).
func (t *T) ExpectGt(x, y interface{}, errorParts ...interface{}) {
	t.expectThat(x, oglematchers.GreaterThan(y), 1, errorParts)
}

// ExpectGe(x, y) is equivalent to
// ExpectThat(x, oglematchers.GreaterOrEqual(y)).
func (t *T) ExpectGe(x, y interface{}, errorParts ...interface{}) {
	t.expectThat(x, oglematchers.GreaterOrEqual(y), 1, errorParts)
}

// ExpectTrue(b) is equivalent to ExpectThat(b, oglematchers.Equals(true)).
func (t *T) ExpectTrue(b interface{}, errorParts ...interface{}) {
	t.expectThat(b, oglematchers.Equals(true), 1, errorParts)
}

// ExpectFalse(b) is equivalent to ExpectThat(b, oglematchers.Equals(false)).
func (t *T) ExpectFalse(b interface{}, errorParts ...interface{}) {
	t.expectThat(b, oglematchers.Equals(false), 1, errorParts)
}
