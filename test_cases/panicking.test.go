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
	"log"
	"testing"

	. "github.com/jacobsa/oglematchers"
	"github.com/jacobsa/ogletest"
)

func TestPanickingTest(t *testing.T) { ogletest.RunTests(t) }

////////////////////////////////////////////////////////////////////////
// PanickingTest
////////////////////////////////////////////////////////////////////////

func someFuncThatPanics() {
	panic("Panic in someFuncThatPanics")
}

type PanickingTest struct {
}

func init() { ogletest.RegisterTestSuite(&PanickingTest{}) }

func (s *PanickingTest) TearDown(t *ogletest.T) {
	t.Logf("TearDown running.")
}

func (s *PanickingTest) ExplicitPanic(t *ogletest.T) {
	panic("Panic in ExplicitPanic")
}

func (s *PanickingTest) ExplicitPanicInHelperFunction(t *ogletest.T) {
	someFuncThatPanics()
}

func (s *PanickingTest) NilPointerDerefence(t *ogletest.T) {
	var p *int
	log.Println(*p)
}

func (s *PanickingTest) ZzzSomeOtherTest(t *ogletest.T) {
	t.ExpectThat(17, Equals(17.0))
}

////////////////////////////////////////////////////////////////////////
// SetUpPanicTest
////////////////////////////////////////////////////////////////////////

type SetUpPanicTest struct {
}

func init() { ogletest.RegisterTestSuite(&SetUpPanicTest{}) }

func (s *SetUpPanicTest) SetUp(t *ogletest.T) {
	t.Logf("SetUp about to panic.")
	panic("Panic in SetUp")
}

func (s *SetUpPanicTest) TearDown(t *ogletest.T) {
	t.Logf("TearDown running.")
}

func (s *SetUpPanicTest) SomeTestCase(t *ogletest.T) {
}

////////////////////////////////////////////////////////////////////////
// TearDownPanicTest
////////////////////////////////////////////////////////////////////////

type TearDownPanicTest struct {
}

func init() { ogletest.RegisterTestSuite(&TearDownPanicTest{}) }

func (s *TearDownPanicTest) TearDown(t *ogletest.T) {
	t.Logf("TearDown about to panic.")
	panic("Panic in TearDown")
}

func (s *TearDownPanicTest) SomeTestCase(t *ogletest.T) {
}
