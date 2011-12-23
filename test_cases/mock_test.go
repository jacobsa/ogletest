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
	"github.com/jacobsa/oglemock"
	"reflect"
	"testing"
	"unsafe"
)

////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////

type MockTest struct {
	controller oglemock.Controller
}

func init()                     { RegisterTestSuite(&MockTest{}) }
func TestOgletest(t *testing.T) { RunTests(t) }

func (t *MockTest) SetUp(i *TestInfo) {
	t.controller = i.MockController
}

// TODO(jacobsa): Replace this with an auto-generated mock class when oglemock
// supoorts it.
type mockFooer struct {
	controller oglemock.Controller
}

func (f *mockFooer) Oglemock_Id() uintptr {
	return uintptr(unsafe.Pointer(f))
}

func (f *mockFooer) Oglemock_Description() string {
	return "some mockFooer"
}

func (f *mockFooer) DoFoo(s string) int {
	retVals := f.controller.HandleMethodCall(
		f,
		"DoFoo",
		"blah.go",
		112,
		[]interface{}{s})

	return int(reflect.ValueOf(retVals[0]).Int())
}

////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////

func (t *MockTest) ExpectationSatisfied() {
	f := &mockFooer{t.controller}

	// TODO(jacobsa): Replace this hand-spun expectation with a call to a more
	// convenient ExpectCall function when one is available. See issue #8.
	t.controller.ExpectCall(
		f,
		"DoFoo",
		"blah_test.go",
		117)(HasSubstr("taco")).WillOnce(oglemock.Return(17))

	ExpectEq(17, f.DoFoo("burritos and tacos"))
}

func (t *MockTest) MockExpectationNotSatisfied() {
	f := &mockFooer{t.controller}

	// TODO(jacobsa): Replace this hand-spun expectation with a call to a more
	// convenient ExpectCall function when one is available. See issue #8.
	t.controller.ExpectCall(
		f,
		"DoFoo",
		"blah_test.go",
		117)(HasSubstr("taco"))
}

func (t *MockTest) UnexpectedCall() {
	f := &mockFooer{t.controller}

	// TODO(jacobsa): Replace this hand-spun expectation with a call to a more
	// convenient ExpectCall function when one is available. See issue #8.
	t.controller.ExpectCall(
		f,
		"DoFoo",
		"blah_test.go",
		117)(HasSubstr("taco")).WillOnce(oglemock.Return(17))

	f.DoFoo("blah")
}
