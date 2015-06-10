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

type workItem struct {
	// The test function to be run.
	tf TestFunction

	// Results of executing the test function. Valid only once complete has been
	// closed.
	failed   bool
	output   []byte
	duration time.Duration

	complete chan struct{}
}

// Run test functions from the channel, closing 'complete' channels when
// finished with them. Return early without draining the input channel if
// gStopRunning is closed.
func processWork(
	work <-chan *workItem) {
	for wi := range work {
		// Check whether we should return early.
		select {
		default:
		case <-gStopRunning:
			return
		}

		// Run the test function.
		startTime := time.Now()
		wi.failed, wi.output = runTestFunction(wi.tf)
		wi.duration = time.Since(startTime)

		close(wi.complete)
	}
}

////////////////////////////////////////////////////////////////////////
// Test suites
////////////////////////////////////////////////////////////////////////

// Wait for results for each test function in the suite of the given name,
// signalling failures to the supplied testing.T.
func processSuiteResults(
	t *testing.T,
	suiteName string,
	work []workItem,
	workers *sync.WaitGroup) {
	// If the overall test target has already failed and we've been told to stop
	// on failure, then don't do anything.
	if t.Failed() && *fStopEarly {
		return
	}

	// Print results.
	fmt.Printf("[----------] Running tests from %s\n", suiteName)
	for i := range work {
		wi := &work[i]

		// Print a banner for the start of this test function.
		fmt.Printf("[ RUN      ] %s.%s\n", suiteName, wi.tf.Name)

		// Wait for the result. Special case: if the user has told us to finish up,
		// join the workers and then exit in error.
		select {
		case <-gStopRunning:
			workers.Wait()
			fmt.Println("Exiting early due to user request.")
			os.Exit(1)

		case <-wi.complete:
		}

		// Mark the test as having failed if appropriate.
		if wi.failed {
			t.Fail()
		}

		// Print output.
		fmt.Printf("%s", wi.output)

		// Print a banner for the end of the test.
		bannerMessage := "[       OK ]"
		if wi.failed {
			bannerMessage = "[  FAILED  ]"
		}

		// Print a summary of the time taken, if long enough.
		var timeMessage string
		if wi.duration >= 25*time.Millisecond {
			timeMessage = fmt.Sprintf(" (%v)", wi.duration)
		}

		fmt.Printf(
			"%s %s.%s%s\n",
			bannerMessage,
			suiteName,
			wi.tf.Name,
			timeMessage)

		// Stop printing results from this suite if we've been told to stop on
		// failure and we have failed already.
		if t.Failed() && *fStopEarly {
			break
		}
	}

	fmt.Printf("[----------] Finished with tests from %s\n", suiteName)
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
	// For each suite:
	//
	//  *  Filter to the set of test functions the user has asked us to run.
	//  *  Set up a slice containing the work to be handed off below.
	//
	type suiteAndWork struct {
		suiteName string
		work      []workItem
	}

	var suites []suiteAndWork
	var totalWorkItems int
	for _, s := range registeredSuites {
		suite := suiteAndWork{suiteName: s.Name}
		for _, tf := range filterTestFunctions(s) {
			wi := workItem{
				tf:       tf,
				complete: make(chan struct{}),
			}

			suite.work = append(suite.work, wi)
			totalWorkItems++
		}

		suites = append(suites, suite)
	}

	// Set up a channel containing all of the work to be divvied out to the
	// workers below.
	workChan := make(chan *workItem, totalWorkItems)
	for _, suite := range suites {
		for i := range suite.work {
			workChan <- &suite.work[i]
		}
	}
	close(workChan)

	// Start several workers processing work concurrently.
	var wg sync.WaitGroup
	for i := 0; i < *fParallelism; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			processWork(workChan)
		}()
	}

	// Process results for each suite.
	for _, suite := range suites {
		processSuiteResults(t, suite.suiteName, suite.work, &wg)
	}
}

////////////////////////////////////////////////////////////////////////
// StopRunningTests
////////////////////////////////////////////////////////////////////////

// A channel that is closed when the user wants us to halt after waiting for
// the currently running tests to complete.
var gStopRunning = make(chan struct{})

// Protect StopRunningTests from closing gStopRunning twice.
var gCloseStopRunning sync.Once

// Request that RunTests stop running additional tests and cause the program to
// exit with a non-zero status when the currently running tests finish.
func StopRunningTests() {
	gCloseStopRunning.Do(func() {
		close(gStopRunning)
	})
}
