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
	"log"
	"testing"

	. "github.com/jacobsa/oglematchers"
	. "github.com/jacobsa/ogletest"
)

func TestOgletest(t *testing.T) { RunTests(t) }

////////////////////////////////////////////////////////////////////////
// PanickingTest
////////////////////////////////////////////////////////////////////////

type PanickingTest struct {
}

func init() { RegisterTestSuite(&PanickingTest{}) }

func (t *PanickingTest) TearDown() {
	fmt.Println("TearDown running.")
}

func (t *PanickingTest) ExplicitPanic() {
	panic("foobar")
}

func (t *PanickingTest) NilPointerDerefence() {
	var p *int
	log.Println(*p)
}

func (t *PanickingTest) ZzzSomeOtherTest() {
	ExpectThat(17, Equals(17.0))
}

////////////////////////////////////////////////////////////////////////
// SetUpPanicTest
////////////////////////////////////////////////////////////////////////

type SetUpPanicTest struct {
}

func init() { RegisterTestSuite(&SetUpPanicTest{}) }

func (t *SetUpPanicTest) SetUp(ti *TestInfo) {
	fmt.Println("SetUp about to panic.")
	panic("baz")
}

func (t *SetUpPanicTest) TearDown() {
	fmt.Println("TearDown running.")
}

func (t *SetUpPanicTest) SomeTestCase() {
	panic("foobar")
}
