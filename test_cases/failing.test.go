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
	"testing"

	. "github.com/jacobsa/oglematchers"
	"github.com/jacobsa/ogletest"
)

func TestFailingTest(t *testing.T) { ogletest.RunTests(t) }

////////////////////////////////////////////////////////////////////////
// Usual failures
////////////////////////////////////////////////////////////////////////

type FailingTest struct {
}

var _ ogletest.TearDownInterface = &FailingTest{}

func init() { ogletest.RegisterTestSuite(&FailingTest{}) }

func (s *FailingTest) TearDown(t *ogletest.T) {
	t.Logf("TearDown running.")
}

func (s *FailingTest) PassingMethod(t *ogletest.T) {
}

func (s *FailingTest) Equals(t *ogletest.T) {
	t.ExpectThat(17, Equals(17.5))
	t.ExpectThat(17, Equals("taco"))
}

func (s *FailingTest) LessThan(t *ogletest.T) {
	t.ExpectThat(18, LessThan(17))
	t.ExpectThat(18, LessThan("taco"))
}

func (s *FailingTest) HasSubstr(t *ogletest.T) {
	t.ExpectThat("taco", HasSubstr("ac"))
	t.ExpectThat(17, HasSubstr("ac"))
}

func (s *FailingTest) ExpectWithUserErrorMessages(t *ogletest.T) {
	t.ExpectThat(17, Equals(19), "foo bar: %d", 112)
	t.ExpectEq(17, 17.5, "foo bar: %d", 112)
	t.ExpectLe(17, 16.9, "foo bar: %d", 112)
	t.ExpectLt(17, 16.9, "foo bar: %d", 112)
	t.ExpectGe(17, 17.1, "foo bar: %d", 112)
	t.ExpectGt(17, "taco", "foo bar: %d", 112)
	t.ExpectNe(17, 17.0, "foo bar: %d", 112)
	t.ExpectFalse(true, "foo bar: %d", 112)
	t.ExpectTrue(false, "foo bar: %d", 112)
}

func (s *FailingTest) AssertWithUserErrorMessages(t *ogletest.T) {
	t.AssertThat(17, Equals(19), "foo bar: %d", 112)
}

func (s *FailingTest) ExpectationAliases(t *ogletest.T) {
	t.ExpectEq(17, 17.5)
	t.ExpectEq("taco", 17.5)

	t.ExpectLe(17, 16.9)
	t.ExpectLt(17, 16.9)
	t.ExpectLt(17, "taco")

	t.ExpectGe(17, 17.1)
	t.ExpectGt(17, 17.1)
	t.ExpectGt(17, "taco")

	t.ExpectNe(17, 17.0)
	t.ExpectNe(17, "taco")

	t.ExpectFalse(true)
	t.ExpectFalse("taco")

	t.ExpectTrue(false)
	t.ExpectTrue("taco")
}

func (s *FailingTest) AssertThatFailure(t *ogletest.T) {
	t.AssertThat(17, Equals(19))
	panic("Shouldn's get here.")
}

func (s *FailingTest) AssertEqFailure(t *ogletest.T) {
	t.AssertEq(19, 17)
	panic("Shouldn's get here.")
}

func (s *FailingTest) AssertNeFailure(t *ogletest.T) {
	t.AssertNe(19, 19)
	panic("Shouldn's get here.")
}

func (s *FailingTest) AssertLeFailure(t *ogletest.T) {
	t.AssertLe(19, 17)
	panic("Shouldn's get here.")
}

func (s *FailingTest) AssertLtFailure(t *ogletest.T) {
	t.AssertLt(19, 17)
	panic("Shouldn's get here.")
}

func (s *FailingTest) AssertGeFailure(t *ogletest.T) {
	t.AssertGe(17, 19)
	panic("Shouldn's get here.")
}

func (s *FailingTest) AssertGtFailure(t *ogletest.T) {
	t.AssertGt(17, 19)
	panic("Shouldn's get here.")
}

func (s *FailingTest) AssertTrueFailure(t *ogletest.T) {
	t.AssertTrue("taco")
	panic("Shouldn's get here.")
}

func (s *FailingTest) AssertFalseFailure(t *ogletest.T) {
	t.AssertFalse("taco")
	panic("Shouldn's get here.")
}

func (s *FailingTest) AddFailureRecord(t *ogletest.T) {
	r := ogletest.FailureRecord{
		FileName:   "foo.go",
		LineNumber: 17,
		Error:      "taco\nburrito",
	}

	t.AddFailureRecord(r)
}

func (s *FailingTest) AddFailure(t *ogletest.T) {
	t.AddFailure("taco")
	t.AddFailure("burrito: %d", 17)
}

func (s *FailingTest) AddFailureThenAbortTest(t *ogletest.T) {
	t.AddFailure("enchilada")
	t.AbortTest()
	t.Logf("Shouldn't get here.")
}

////////////////////////////////////////////////////////////////////////
// Expectation failure during SetUp
////////////////////////////////////////////////////////////////////////

type ExpectFailDuringSetUpTest struct {
}

func init() { ogletest.RegisterTestSuite(&ExpectFailDuringSetUpTest{}) }

func (s *ExpectFailDuringSetUpTest) SetUp(t *ogletest.T) {
	t.ExpectFalse(true)
}

func (s *ExpectFailDuringSetUpTest) TearDown(t *ogletest.T) {
	t.Logf("TearDown running.")
}

func (s *ExpectFailDuringSetUpTest) PassingMethod(t *ogletest.T) {
	t.Logf("Method running.")
}

////////////////////////////////////////////////////////////////////////
// Assertion failure during SetUp
////////////////////////////////////////////////////////////////////////

type AssertFailDuringSetUpTest struct {
}

func init() { ogletest.RegisterTestSuite(&AssertFailDuringSetUpTest{}) }

func (s *AssertFailDuringSetUpTest) SetUp(t *ogletest.T) {
	t.AssertFalse(true)
}

func (s *AssertFailDuringSetUpTest) TearDown(t *ogletest.T) {
	t.Logf("TearDown running.")
}

func (s *AssertFailDuringSetUpTest) PassingMethod(t *ogletest.T) {
	t.Logf("Method running.")
}

////////////////////////////////////////////////////////////////////////
// Expectation failure during TearDown
////////////////////////////////////////////////////////////////////////

type ExpectFailDuringTearDownTest struct {
}

func init() { ogletest.RegisterTestSuite(&ExpectFailDuringTearDownTest{}) }

func (s *ExpectFailDuringTearDownTest) SetUp(t *ogletest.T) {
	t.Logf("SetUp running.")
}

func (s *ExpectFailDuringTearDownTest) TearDown(t *ogletest.T) {
	t.ExpectFalse(true)
}

func (s *ExpectFailDuringTearDownTest) PassingMethod(t *ogletest.T) {
	t.Logf("Method running.")
}

////////////////////////////////////////////////////////////////////////
// Assertion failure during TearDown
////////////////////////////////////////////////////////////////////////

type AssertFailDuringTearDownTest struct {
}

func init() { ogletest.RegisterTestSuite(&AssertFailDuringTearDownTest{}) }

func (s *AssertFailDuringTearDownTest) SetUp(t *ogletest.T) {
	t.Logf("SetUp running.")
}

func (s *AssertFailDuringTearDownTest) TearDown(t *ogletest.T) {
	t.AssertFalse(true)
}

func (s *AssertFailDuringTearDownTest) PassingMethod(t *ogletest.T) {
	t.Logf("Method running.")
}
