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

type FailingTest struct {
}

func init()                     { RegisterTestSuite(&FailingTest{}) }
func TestOgletest(t *testing.T) { RunTests(t) }

////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////

func (t *FailingTest) PassingMethod() {
}

func (t *FailingTest) Equals() {
	ExpectThat(17, Equals(17.5))
	ExpectThat(17, Equals("taco"))
}

func (t *FailingTest) LessThan() {
	ExpectThat(18, LessThan(17))
	ExpectThat(18, LessThan("taco"))
}

func (t *FailingTest) HasSubstr() {
	ExpectThat("taco", HasSubstr("ac"))
	ExpectThat(17, HasSubstr("ac"))
}

func (t *FailingTest) UserErrorMessage() {
	ExpectThat(17, Equals(19), "foo bar")
}

func (t *FailingTest) ModifiedExpectation() {
	ExpectThat(17, HasSubstr("ac")).SetCaller("foo.go", 112)
	ExpectEq(17, 19).SetCaller("bar.go", 117)
}

func (t FailingTest) ExpectationAliases() {
	ExpectEq(17, 17.5)
	ExpectEq("taco", 17.5)

	ExpectLe(17, 16.9)
	ExpectLt(17, 16.9)
	ExpectLt(17, "taco")

	ExpectGe(17, 17.1)
	ExpectGt(17, 17.1)
	ExpectGt(17, "taco")

	ExpectNe(17, 17.0)
	ExpectNe(17, "taco")

	ExpectFalse(true)
	ExpectFalse("taco")

	ExpectTrue(false)
	ExpectTrue("taco")
}
