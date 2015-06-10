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
	"fmt"
	"testing"
	"time"

	. "github.com/jacobsa/oglematchers"
	"github.com/jacobsa/ogletest"
)

func TestPassingTest(t *testing.T) { ogletest.RunTests(t) }

////////////////////////////////////////////////////////////////////////
// PassingTest
////////////////////////////////////////////////////////////////////////

type PassingTest struct {
}

func init() { ogletest.RegisterTestSuite(&PassingTest{}) }

func (s *PassingTest) EmptyTestMethod(t *ogletest.T) {
}

func (s *PassingTest) SuccessfullMatches(t *ogletest.T) {
	t.ExpectThat(17, Equals(17.0))
	t.ExpectThat(16.9, LessThan(17))
	t.ExpectThat("taco", HasSubstr("ac"))

	t.AssertThat(17, Equals(17.0))
	t.AssertThat(16.9, LessThan(17))
	t.AssertThat("taco", HasSubstr("ac"))
}

func (s *PassingTest) ExpectAliases(t *ogletest.T) {
	t.ExpectEq(17, 17.0)

	t.ExpectLe(17, 17.0)
	t.ExpectLe(17, 18.0)
	t.ExpectLt(17, 18.0)

	t.ExpectGe(17, 17.0)
	t.ExpectGe(17, 16.0)
	t.ExpectGt(17, 16.0)

	t.ExpectNe(17, 18.0)

	t.ExpectTrue(true)
	t.ExpectFalse(false)
}

func (s *PassingTest) AssertAliases(t *ogletest.T) {
	t.AssertEq(17, 17.0)

	t.AssertLe(17, 17.0)
	t.AssertLe(17, 18.0)
	t.AssertLt(17, 18.0)

	t.AssertGe(17, 17.0)
	t.AssertGe(17, 16.0)
	t.AssertGt(17, 16.0)

	t.AssertNe(17, 18.0)

	t.AssertTrue(true)
	t.AssertFalse(false)
}

func (s *PassingTest) SlowTest(t *ogletest.T) {
	time.Sleep(37 * time.Millisecond)
}

////////////////////////////////////////////////////////////////////////
// PassingTestWithHelpers
////////////////////////////////////////////////////////////////////////

type PassingTestWithHelpers struct {
}

var _ ogletest.SetUpInterface = &PassingTestWithHelpers{}
var _ ogletest.TearDownInterface = &PassingTestWithHelpers{}

func init() { ogletest.RegisterTestSuite(&PassingTestWithHelpers{}) }

func (s *PassingTestWithHelpers) SetUp(t *ogletest.T) {
	fmt.Println("SetUp ran.")
}

func (s *PassingTestWithHelpers) TearDown(t *ogletest.T) {
	fmt.Println("TearDown ran.")
}

func (s *PassingTestWithHelpers) EmptyTestMethod(t *ogletest.T) {
}
