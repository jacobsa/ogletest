// Copyright 2015 Aaron Jacobs. All Rights Reserved.
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
	"path"
	"runtime"
	"sync"

	"github.com/jacobsa/oglemock"
	"golang.org/x/net/context"
)

// T is a type that contains dependencies for test functions and whose methods
// can be used by test functions to control test execution, including adding
// failure messages.
type T struct {
	// A context that the test should use for any context-aware code it calls.
	// May carry tracing information or be used for cancellation.
	Ctx context.Context

	// A mock controller that can be used for creating mocks for use by the test.
	// Mock errors will be associated with the test. The Finish method should not
	// be run by the user; ogletest will do that automatically after the test
	// finishes.
	MockController oglemock.Controller

	/////////////////////////
	// Constant data
	/////////////////////////

	testName string

	/////////////////////////
	// Mutable state
	/////////////////////////

	mu sync.Mutex

	// Failure records accumulated so far. Only ever appended to.
	//
	// GUARDED_BY(mu)
	records []FailureRecord
}

func newT(
	ctx context.Context,
	name string) (t *T) {
	t = &T{
		Ctx:      ctx,
		testName: name,
	}

	t.MockController = oglemock.NewController(tErrorReporter{t})

	return
}

func (t *T) name() string {
	return t.testName
}

func (t *T) failureRecords() []FailureRecord {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.records
}

// FailureRecord represents a single failed expectation or assertion for a
// test. Most users don't want to interact with these directly; they are
// generated implicitly using ExpectThat, AssertThat, ExpectLt, etc.
type FailureRecord struct {
	// The file name within which the expectation failed, e.g. "foo_test.go".
	FileName string

	// The line number at which the expectation failed.
	LineNumber int

	// The error associated with the file:line pair above. For example, the
	// following expectation:
	//
	//     ExpectEq(17, "taco")"
	//
	// May cause this error:
	//
	//     Expected: 17
	//     Actual:   "taco", which is not numeric
	//
	Error string
}

// Record a failure (and continue running the test).
//
// Most users will want to use ExpectThat, ExpectEq, etc. instead of this
// function. Those that do want to report arbitrary errors will probably be
// satisfied with AddFailure, which is easier to use.
func (t *T) AddFailureRecord(r FailureRecord) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.records = append(t.records, r)
}

// Call AddFailureRecord with a record whose file name and line number come
// from the caller of this function, and whose error string is created by
// calling fmt.Sprintf using the arguments to this function.
func (t *T) AddFailure(format string, a ...interface{}) {
	r := FailureRecord{
		Error: fmt.Sprintf(format, a...),
	}

	// Get information about the call site.
	var ok bool
	if _, r.FileName, r.LineNumber, ok = runtime.Caller(1); !ok {
		panic("Can't find caller")
	}

	r.FileName = path.Base(r.FileName)

	t.AddFailureRecord(r)
}

// A sentinel type that is used in a conspiracy between AbortTest and runTests.
// If runTests sees an abortError as the value given to a panic() call, it will
// avoid printing the panic error.
type abortError struct {
}

// Immediately stop executing the test, causing it to fail with the failures
// previously recorded. Behavior is undefined if no failures have been
// recorded.
//
// This function must only be called from the goroutine on which the test was
// initially started.
func (t *T) AbortTest() {
	panic(abortError{})
}

// tErrorReporter is an oglemock.ErrorReporter that writes failure records into
// a T.
type tErrorReporter struct {
	t *T
}

func (r tErrorReporter) ReportError(
	fileName string,
	lineNumber int,
	err error) {
	record := FailureRecord{
		FileName:   fileName,
		LineNumber: lineNumber,
		Error:      err.Error(),
	}

	r.t.AddFailureRecord(record)
}

func (r tErrorReporter) ReportFatalError(
	fileName string,
	lineNumber int,
	err error) {
	r.ReportError(fileName, lineNumber, err)
	r.t.AbortTest()
}
