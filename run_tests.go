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
	"os"
	"path"
	"regexp"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/net/context"

	"github.com/jacobsa/reqtrace"
)

////////////////////////////////////////////////////////////////////////
// Flags
////////////////////////////////////////////////////////////////////////

var fTestFilter = flag.String(
	"ogletest.run",
	"",
	"Regexp for matching tests to run.")

var fStopEarly = flag.Bool(
	"ogletest.stop_early",
	false,
	"If true, stop after the first failure.")

var fParallelism = flag.Int(
	"ogletest.parallelism",
	16,
	"The maximum number of tests to run concurrently.")

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

func isAbortError(x interface{}) bool {
	_, ok := x.(abortError)
	return ok
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
func runWithProtection(t *T, f func(*T)) (panicked bool) {
	defer func() {
		// If the test didn't panic, we're done.
		r := recover()
		if r == nil {
			return
		}

		panicked = true

		// If the function panicked (and the panic was not due to an AssertThat
		// failure), add a failure for the panic.
		if !isAbortError(r) {
			var panicRecord FailureRecord
			panicRecord.FileName, panicRecord.LineNumber = findPanicFileLine()
			panicRecord.Error = fmt.Sprintf(
				"panic: %v\n\n%s", r, formatPanicStack())

			t.AddFailureRecord(panicRecord)
		}
	}()

	f(t)
	return
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

// Filter test functions according to the user-supplied filter flag.
func filterTestFunctions(suite TestSuite) (out []TestFunction) {
	re, err := regexp.Compile(*fTestFilter)
	if err != nil {
		panic("Invalid value for --ogletest.run: " + err.Error())
	}

	for _, tf := range suite.TestFunctions {
		fullName := fmt.Sprintf("%s.%s", suite.Name, tf.Name)
		if !re.MatchString(fullName) {
			continue
		}

		out = append(out, tf)
	}

	return
}

////////////////////////////////////////////////////////////////////////
// Test functions
////////////////////////////////////////////////////////////////////////

// Run a single test function, returning an indication of whether it failed and
// its output.
func runTestFunction(tf TestFunction) (failed bool, output []byte) {
	// Start a trace.
	ctx, reportOutcome := reqtrace.Trace(context.Background(), tf.Name)

	// Create the T.
	t := newT(ctx, tf.Name)

	// Run the SetUp function, if any, paying attention to whether it panics.
	setUpPanicked := false
	if tf.SetUp != nil {
		setUpPanicked = runWithProtection(t, tf.SetUp)
	}

	// Run the test function itself, but only if the SetUp function didn't panic.
	// (This includes AssertThat errors.)
	if !setUpPanicked {
		runWithProtection(t, tf.Run)
	}

	// Run the TearDown function, if any.
	if tf.TearDown != nil {
		runWithProtection(t, tf.TearDown)
	}

	// Tell the mock controller for the tests to report any errors it's sitting
	// on.
	t.MockController.Finish()

	// Find out what happened.
	failed, output = t.result()

	// Report the outcome to reqtrace.
	if failed {
		reportOutcome(fmt.Errorf("failed"))
	} else {
		reportOutcome(nil)
	}

	return
}

////////////////////////////////////////////////////////////////////////
// Test suites
////////////////////////////////////////////////////////////////////////

// Run a single test suite, signalling failures to the supplied testing.T.
//
// TODO(jacobsa): Hoist parallelism so that we can process multiple suites in
// parallel?
func runTestSuite(
	t *testing.T,
	suite TestSuite) {
	tfs := filterTestFunctions(suite)

	// Ensure that if we exit this function early due to StopRunningTests being
	// called, we exit the program with an error. This prevents us from skipping
	// test functions but making the program look like it succeeded.
	defer func() {
		if atomic.LoadUint64(&gStopRunning) != 0 {
			fmt.Println("Exiting early due to user request.")
			os.Exit(1)
			return
		}
	}()

	// If the overall test target has already failed and we've been told to stop
	// on failure, then don't do anything.
	if t.Failed() && *fStopEarly {
		return
	}

	// Set up a channel containing indices of test functions to be run. This will
	// be used to assign work to workers.
	indices := make(chan int, len(tfs))
	for i := 0; i < len(tfs); i++ {
		indices <- i
	}
	close(indices)

	// Set up a slice of channels into which results will be written. These will
	// be used to communicate results from the workers, in order.
	type result struct {
		failed   bool
		output   []byte
		duration time.Duration
	}

	var resultChans []chan result
	for i := 0; i < len(tfs); i++ {
		resultChans = append(resultChans, make(chan result, 1))
	}

	// Start several workers processing work in parallel.
	var wg sync.WaitGroup
	for i := 0; i < *fParallelism; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range indices {
				// Special case: if the user has asked us to stop running additional
				// tests, then do so.
				if atomic.LoadUint64(&gStopRunning) != 0 {
					return
				}

				startTime := time.Now()
				failed, output := runTestFunction(tfs[i])
				duration := time.Since(startTime)

				resultChans[i] <- result{failed, output, duration}
			}
		}()
	}

	// Print results.
	fmt.Printf("[----------] Running tests from %s\n", suite.Name)
	for i, tf := range tfs {
		// If the user has asked us to stop running tests, then wait for all
		// workers to exit (in order to fulfill the guarantee made by
		// StopRunningTests that tests will finish) and then bail out.
		if atomic.LoadUint64(&gStopRunning) == 1 {
			wg.Wait()
			return
		}

		// Print a banner for the start of this test function.
		fmt.Printf("[ RUN      ] %s.%s\n", suite.Name, tf.Name)

		// Wait for the result.
		result := <-resultChans[i]

		// Mark the test as having failed if appropriate.
		if result.failed {
			t.Fail()
		}

		// Print output.
		fmt.Printf("%s", result.output)

		// Print a banner for the end of the test.
		bannerMessage := "[       OK ]"
		if result.failed {
			bannerMessage = "[  FAILED  ]"
		}

		// Print a summary of the time taken, if long enough.
		var timeMessage string
		if result.duration >= 25*time.Millisecond {
			timeMessage = fmt.Sprintf(" (%v)", result.duration)
		}

		fmt.Printf(
			"%s %s.%s%s\n",
			bannerMessage,
			suite.Name,
			tf.Name,
			timeMessage)

		// Stop printing results from this suite if we've been told to stop on
		// failure and we have failed already.
		if t.Failed() && *fStopEarly {
			break
		}
	}

	fmt.Printf("[----------] Finished with tests from %s\n", suite.Name)
}

////////////////////////////////////////////////////////////////////////
// RunTests
////////////////////////////////////////////////////////////////////////

// runTestsOnce protects RunTests from executing multiple times.
var runTestsOnce sync.Once

// Run everything registered with Register (including via the wrapper
// RegisterTestSuite).
//
// Failures are communicated to the supplied testing.T object. This is the
// bridge between ogletest and the testing package (and `go test`); you should
// ensure that it's called at least once by creating a test function compatible
// with `go test` and calling it there.
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
	for _, suite := range registeredSuites {
		runTestSuite(t, suite)
	}
}

////////////////////////////////////////////////////////////////////////
// StopRunningTests
////////////////////////////////////////////////////////////////////////

// Signalling between RunTests and StopRunningTests.
var gStopRunning uint64

// Request that RunTests stop running additional tests and cause the program to
// exit with a non-zero status when the currently running tests finish.
func StopRunningTests() {
	atomic.StoreUint64(&gStopRunning, 1)
}
