// Copyright 2012 Aaron Jacobs. All Rights Reserved.
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
	. "github.com/jacobsa/oglematchers"
	"reflect"
	"testing"
)

func TestRegisterMethodsTest(t *testing.T) { RunTests(t) }

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

type MethodsTest struct {
}

func init() { RegisterTestSuite(&MethodsTest{}) }

type OneMethodType int
func (x OneMethodType) Foo() {}

type MultipleMethodsType int
func (x MultipleMethodsType) Foo() {}
func (x MultipleMethodsType) Bar() {}
func (x MultipleMethodsType) Baz() {}

type SingleLineType int
func (x SingleLineType) Foo() {}; func (x SingleLineType) Bar() {}
func (x SingleLineType) Baz() {}; func (x SingleLineType) Qux() {}

////////////////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////////////////

func (t *MethodsTest) NoMethods() {
	type foo int

	methods := getMethodsInSourceOrder(reflect.TypeOf(foo(17)))
	ExpectThat(methods, ElementsAre())
}

func (t *MethodsTest) OneMethod() {
	methods := getMethodsInSourceOrder(reflect.TypeOf(OneMethodType(17)))
	AssertThat(methods, ElementsAre(Any()))

	ExpectEq("Foo", methods[0].Name)
}

func (t *MethodsTest) MultipleMethods() {
	methods := getMethodsInSourceOrder(reflect.TypeOf(MultipleMethodsType(17)))
	AssertThat(methods, ElementsAre(Any(), Any(), Any()))

	ExpectEq("Foo", methods[0].Name)
	ExpectEq("Bar", methods[1].Name)
	ExpectEq("Baz", methods[2].Name)
}

func (t *MethodsTest) MultipleMethodsOnSingleLine() {
	methods := getMethodsInSourceOrder(reflect.TypeOf(SingleLineType(17)))
	AssertThat(methods, ElementsAre(Any(), Any(), Any(), Any()))

	// TODO(jacobsa): Delete this block of code when the following issue is
	// resolved:
	//     http://code.google.com/p/go/issues/detail?id=4174
	ExpectThat(methods, Contains("Foo"))
	ExpectThat(methods, Contains("Bar"))
	ExpectThat(methods, Contains("Baz"))
	ExpectThat(methods, Contains("Qux"))
	return

	ExpectThat(methods[0].Name, AnyOf("Foo", "Bar"))
	ExpectThat(methods[1].Name, AnyOf("Foo", "Bar"))
	ExpectNe(methods[0].Name, methods[1].Name)

	ExpectThat(methods[2].Name, AnyOf("Baz", "Qux"))
	ExpectThat(methods[3].Name, AnyOf("Baz", "Qux"))
	ExpectNe(methods[2].Name, methods[3].Name)
}
