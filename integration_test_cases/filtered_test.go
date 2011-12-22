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

package oglematchers_test

import (
	. "github.com/jacobsa/oglematchers"
	. "github.com/jacobsa/ogletest"
	"testing"
)

////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////

type FilteredTest struct {
}

func init()                     { RegisterTestSuite(&FilteredTest{}) }
func TestOgletest(t *testing.T) { RunTests(t) }

////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////

func (t *FilteredTest) PassingTestFoo() {
	ExpectThat(19, Equals(19))
}

func (t *FilteredTest) PassingTestBar() {
	ExpectThat(17, Equals(17))
}

func (t *FilteredTest) FilteredTestFoo() {
	ExpectThat(18, LessThan(17))
}

func (t *FilteredTest) FilteredTestBar() {
	ExpectThat("taco", HasSubstr("blah"))
}

func (t *FilteredTest) FilteredTestBaz() {
	ExpectThat(18, LessThan(17))
}
