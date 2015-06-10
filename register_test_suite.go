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

import (
	"fmt"
	"reflect"

	"github.com/jacobsa/ogletest/srcutil"
)

// Test suites that implement this interface have special meaning to
// RegisterTestSuite.
type SetUpInterface interface {
	// This method is called before each test method is invoked, with the same
	// receiver as that test method. At the time this method is invoked, the
	// receiver is a zero value for the test suite type. Use this method for
	// common setup code that works on data not shared across tests.
	SetUp(t *T)
}

// Test suites that implement this interface have special meaning to
// RegisterTestSuite.
type TearDownInterface interface {
	// This method is called after each test method is invoked, with the same
	// receiver as that test method. Use this method for common cleanup code that
	// works on data not shared across tests.
	TearDown(t *T)
}

// RegisterTestSuite tells ogletest about a test suite containing tests that it
// should run. Any exported method on the type pointed to by the supplied
// prototype value that accepts a single argument of type *T will be treated as
// a test method, with the exception of the methods defined by the following
// interfaces, which when present are treated as described in the documentation
// for those interfaces:
//
//  *  SetUpInterface
//  *  TearDownInterface
//
// Each test method is invoked on a different receiver, which is initially a
// zero value of the test suite type. Methods may be run concurrently.
//
// Example:
//
//     type FooTest struct {
//       // Path to a temporary file used by the tests. Each test gets a
//       // different temporary file.
//       tempFile string
//     }
//     func init() { ogletest.RegisterTestSuite(&FooTest{}) }
//
//     func (s *FooTest) SetUp(t *ogletest.T) {
//       var err error
//       s.tempFile, err = CreateTempFile()
//       t.AssertEq(nil, err)
//     }
//
//     func (s *FooTest) TearDown(t *ogletest.T) {
//       err := DeleteTempFile(s.tempFile)
//       t.AssertEq(nil, err)
//     }
//
//     func (s *FooTest) FrobinicatorIsSuccessfullyTweaked(t *ogletest.T) {
//       res := DoSomething(s.tempFile)
//       t.ExpectTrue(res)
//     }
//
func RegisterTestSuite(p interface{}) {
	if p == nil {
		panic("RegisterTestSuite called with nil suite.")
	}

	val := reflect.ValueOf(p)
	typ := val.Type()

	// We will transform to a TestSuite struct.
	suite := TestSuite{}
	suite.Name = typ.Elem().Name()

	// Transform a list of test methods for the suite, filtering them to just the
	// ones that we don't need to skip.
	methods := filterMethods(suite.Name, srcutil.GetMethodsInSourceOrder(typ))
	for _, method := range methods {
		var tf TestFunction
		tf.Name = method.Name

		// Create an instance to be operated on by all of the TestFunction's
		// internal functions.
		instance := reflect.New(typ.Elem())

		// Bind the functions to the instance.
		if i, ok := instance.Interface().(SetUpInterface); ok {
			tf.SetUp = i.SetUp
		}

		methodCopy := method
		tf.Run = func(t *T) { runTestMethod(t, instance, methodCopy) }

		if i, ok := instance.Interface().(TearDownInterface); ok {
			tf.TearDown = i.TearDown
		}

		// Save the TestFunction.
		suite.TestFunctions = append(suite.TestFunctions, tf)
	}

	// Register the suite.
	Register(suite)
}

func runTestMethod(t *T, suite reflect.Value, method reflect.Method) {
	// TODO(jacobsa): Put type checking logic in filterMethods, and just leave
	// out non-conforming methods. Then delete this check. Make sure to add an
	// integration test that shows methods being left out.
	if method.Func.Type().NumIn() != 2 {
		panic(fmt.Sprintf(
			"%s: expected 2 args, actually %d.",
			method.Name,
			method.Func.Type().NumIn()))
	}

	method.Func.Call([]reflect.Value{suite, reflect.ValueOf(t)})
}

func filterMethods(
	suiteName string,
	in []reflect.Method) (out []reflect.Method) {
	for _, m := range in {
		// Skip set up, tear down, and unexported methods.
		if isSpecialMethod(m.Name) || !isExportedMethod(m.Name) {
			continue
		}

		out = append(out, m)
	}

	return
}

func isSpecialMethod(name string) bool {
	return (name == "SetUp") ||
		(name == "TearDown")
}

func isExportedMethod(name string) bool {
	return len(name) > 0 && name[0] >= 'A' && name[0] <= 'Z'
}
