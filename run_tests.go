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
	"bytes"
	"flag"
	"fmt"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"sync"
	"testing"
	"time"
)

var testFilter = flag.String("ogletest.run", "", "Regexp for matching tests to run.")

// runTestsOnce protects RunTests from executing multiple times.
var runTestsOnce sync.Once

func isAssertThatError(x interface{}) bool {
	_, ok := x.(*assertThatError)
	return ok
}

// Run a single test function, returning a slice of failure records.
func runTestFunction(tf TestFunction) (failures []*failureRecord) {
	// Set up a clean slate for this test. Make sure to reset it after everything
	// below is finished, so we don't accidentally use it elsewhere.
	currentlyRunningTest = newTestInfo()
	defer func() {
		currentlyRunningTest = nil
	}()

	// Run the SetUp function, if any, paying attention to whether it panics.
	setUpPanicked := false
	if tf.SetUp != nil {
		setUpPanicked = runWithProtection(func() { tf.SetUp(currentlyRunningTest) })
	}

	// Run the test function itself, but only if the SetUp function didn't panic.
	// (This includes AssertThat errors.)
	if !setUpPanicked {
		runWithProtection(tf.Run)
	}

	// Run the TearDown function, if any.
	if tf.TearDown != nil {
		runWithProtection(tf.TearDown)
	}

	// Tell the mock controller for the tests to report any errors it's sitting
	// on.
	currentlyRunningTest.MockController.Finish()

	return currentlyRunningTest.failureRecords
}

// RunTests runs the test suites registered with ogletest, communicating
// failures to the supplied testing.T object. This is the bridge between
// ogletest and the testing package (and gotest); you should ensure that it's
// called at least once by creating a gotest-compatible test function and
// calling it there.
//
// For example:
//
//     import (
//       "github.com/jacobsa/ogletest"
//       "testing"
//     )
//
//     func TestOgletest(t *testing.T) {
//       ogletest.RunTests(t)
//     }
//
func RunTests(t *testing.T) {
	runTestsOnce.Do(func() { runTestsInternal(t) })
}

// runTestsInternal does the real work of RunTests, which simply wraps it in a
// sync.Once.
func runTestsInternal(t *testing.T) {
	// Process each registered suite.
	for _, suite := range testSuites {
		val := reflect.ValueOf(suite)
		typ := val.Type()
		suiteName := typ.Elem().Name()

		// Grab methods for the suite, filtering them to just the ones that we
		// don't need to skip.
		testMethods := filterMethods(suiteName, getMethodsInSourceOrder(typ))

		// Is there anything left to do?
		if len(testMethods) == 0 {
			continue
		}

		fmt.Printf("[----------] Running tests from %s\n", suiteName)

		// Run the SetUpTestSuite method, if any.
		if i, ok := suite.(SetUpTestSuiteInterface); ok {
			i.SetUpTestSuite()
		}

		// Run each method.
		for _, method := range testMethods {
			// Print a banner for the start of this test.
			fmt.Printf("[ RUN      ] %s.%s\n", suiteName, method.Name)

			// Run the test.
			startTime := time.Now()
			failures := runTest(suite, method)
			runDuration := time.Since(startTime)

			// Print any failures, and mark the test as having failed if there are any.
			for _, record := range failures {
				t.Fail()
				userErrorSection := ""
				if record.UserError != "" {
					userErrorSection = record.UserError + "\n"
				}

				fmt.Printf(
					"%s:%d:\n%s\n%s\n",
					record.FileName,
					record.LineNumber,
					record.GeneratedError,
					userErrorSection)
			}

			// Print a banner for the end of the test.
			bannerMessage := "[       OK ]"
			if len(failures) != 0 {
				bannerMessage = "[  FAILED  ]"
			}

			// Print a summary of the time taken, if long enough.
			var timeMessage string
			if runDuration >= 25*time.Millisecond {
				timeMessage = fmt.Sprintf(" (%s)", runDuration.String())
			}

			fmt.Printf(
				"%s %s.%s%s\n",
				bannerMessage,
				suiteName,
				method.Name,
				timeMessage)
		}

		// Run the TearDownTestSuite method, if any.
		if i, ok := suite.(TearDownTestSuiteInterface); ok {
			i.TearDownTestSuite()
		}

		fmt.Printf("[----------] Finished with tests from %s\n", suiteName)
	}
}

// Return true iff the supplied program counter appears to lie within panic().
func isPanic(pc uintptr) bool {
	f := runtime.FuncForPC(pc)
	if f == nil {
		return false
	}

	return f.Name() == "runtime.gopanic" || f.Name() == "runtime.sigpanic"
}

// Find the deepest stack frame containing something that appears to be a
// panic. Return the 'skip' value that a caller to this function would need
// to supply to runtime.Caller for that frame, or a negative number if not found.
func findPanic() int {
	localSkip := -1
	for i := 0; ; i++ {
		// Stop if we've passed the base of the stack.
		pc, _, _, ok := runtime.Caller(i)
		if !ok {
			break
		}

		// Is this a panic?
		if isPanic(pc) {
			localSkip = i
		}
	}

	return localSkip - 1
}

// Attempt to find the file base name and line number for the ultimate source
// of a panic, on the panicking stack. Return a human-readable sentinel if
// unsuccessful.
func findPanicFileLine() (string, int) {
	panicSkip := findPanic()
	if panicSkip < 0 {
		return "(unknown)", 0
	}

	// Find the trigger of the panic.
	_, file, line, ok := runtime.Caller(panicSkip + 1)
	if !ok {
		return "(unknown)", 0
	}

	return path.Base(file), line
}

// Run the supplied function, catching panics (including AssertThat errors) and
// reporting them to the currently-running test as appropriate. Return true iff
// the function panicked.
func runWithProtection(f func()) (panicked bool) {
	defer func() {
		// If the test didn't panic, we're done.
		r := recover()
		if r == nil {
			return
		}

		panicked = true

		// We modify the currently running test below.
		currentlyRunningTest.mutex.Lock()
		defer currentlyRunningTest.mutex.Unlock()

		// If the function panicked (and the panic was not due to an AssertThat
		// failure), add a failure for the panic.
		if !isAssertThatError(r) {
			var panicRecord failureRecord
			panicRecord.FileName, panicRecord.LineNumber = findPanicFileLine()
			panicRecord.GeneratedError = fmt.Sprintf(
				"panic: %v\n\n%s", r, formatPanicStack())

			currentlyRunningTest.failureRecords = append(
				currentlyRunningTest.failureRecords,
				&panicRecord)
		}
	}()

	f()
	return
}

func runTestMethod(suite reflect.Value, method reflect.Method) {
	if method.Func.Type().NumIn() != 1 {
		panic(fmt.Sprintf(
			"%s: expected 1 args, actually %d.",
			method.Name,
			method.Func.Type().NumIn()))
	}

	method.Func.Call([]reflect.Value{suite})
}

func formatPanicStack() string {
	buf := new(bytes.Buffer)

	// Find the panic. If successful, we'll skip to below it. Otherwise, we'll
	// format everything.
	var initialSkip int
	if panicSkip := findPanic(); panicSkip >= 0 {
		initialSkip = panicSkip + 1
	}

	for i := initialSkip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		// Choose a function name to display.
		funcName := "(unknown)"
		if f := runtime.FuncForPC(pc); f != nil {
			funcName = f.Name()
		}

		// Stop if we've gotten as far as the test runner code.
		if funcName == "github.com/jacobsa/ogletest.runTestMethod" ||
			funcName == "github.com/jacobsa/ogletest.runWithProtection" {
			break
		}

		// Add an entry for this frame.
		fmt.Fprintf(buf, "%s\n\t%s:%d\n", funcName, file, line)
	}

	return buf.String()
}

func filterMethods(suiteName string, in []reflect.Method) (out []reflect.Method) {
	for _, m := range in {
		// Skip set up, tear down, and unexported methods.
		if isSpecialMethod(m.Name) || !isExportedMethod(m.Name) {
			continue
		}

		// Has the user told us to skip this method?
		fullName := fmt.Sprintf("%s.%s", suiteName, m.Name)
		matched, err := regexp.MatchString(*testFilter, fullName)
		if err != nil {
			panic("Invalid value for --ogletest.run: " + err.Error())
		}

		if !matched {
			continue
		}

		out = append(out, m)
	}

	return
}

func isSpecialMethod(name string) bool {
	return (name == "SetUpTestSuite") ||
		(name == "TearDownTestSuite") ||
		(name == "SetUp") ||
		(name == "TearDown")
}

func isExportedMethod(name string) bool {
	return len(name) > 0 && name[0] >= 'A' && name[0] <= 'Z'
}
