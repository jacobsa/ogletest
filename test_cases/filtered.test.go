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

func TestFiltered(t *testing.T) { ogletest.RunTests(t) }

////////////////////////////////////////////////////////////////////////
// Partially filtered out
////////////////////////////////////////////////////////////////////////

type PartiallyFilteredTest struct {
}

func init() { ogletest.RegisterTestSuite(&PartiallyFilteredTest{}) }

func (s *PartiallyFilteredTest) PassingTestFoo(t *ogletest.T) {
	t.ExpectThat(19, Equals(19))
}

func (s *PartiallyFilteredTest) PassingTestBar(t *ogletest.T) {
	t.ExpectThat(17, Equals(17))
}

func (s *PartiallyFilteredTest) PartiallyFilteredTestFoo(t *ogletest.T) {
	t.ExpectThat(18, LessThan(17))
}

func (s *PartiallyFilteredTest) PartiallyFilteredTestBar(t *ogletest.T) {
	t.ExpectThat("taco", HasSubstr("blah"))
}

func (s *PartiallyFilteredTest) PartiallyFilteredTestBaz(t *ogletest.T) {
	t.ExpectThat(18, LessThan(17))
}

////////////////////////////////////////////////////////////////////////
// Completely filtered out
////////////////////////////////////////////////////////////////////////

type CompletelyFilteredTest struct {
}

func init() { ogletest.RegisterTestSuite(&CompletelyFilteredTest{}) }

func (s *PartiallyFilteredTest) SomePassingTest(t *ogletest.T) {
	t.ExpectThat(19, Equals(19))
}

func (s *PartiallyFilteredTest) SomeFailingTest(t *ogletest.T) {
	t.ExpectThat(19, Equals(17))
}
