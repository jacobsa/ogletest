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
	"github.com/jacobsa/oglematchers"
	"path"
	"reflect"
	"runtime"
)

// ExpectationModifier is an interface returned by ExpectThat that allows
// callers to set custom information on the expectation.
type ExpectationModifier interface {
	// SetCaller updates the file name and line number associated with the
	// expectation. This allows, for example, a utility function to express that
	// *its* caller should have its line number printed if the expectation fails,
	// instead of the line number of the ExpectThat call within the utility
	// function.
	SetCaller(fileName string, lineNumber int) ExpectationModifier
}

// ExpectThat confirms that the supplied matcher matches the value x, adding a
// failure record to the currently running test if it does not. If additional
// parameters are supplied, the first will be used as a format string for the
// later ones, and the user-supplied error message will be added to the test
// output in the event of a failure.
//
// For example:
//
//     ExpectThat(userName, Equals("jacobsa"))
//     ExpectThat(users[i], Equals("jacobsa"), "while processing user %d", i)
//
func ExpectThat(x interface{}, m oglematchers.Matcher, errorParts ...interface{}) ExpectationModifier {
	// Get information about the call site.
	_, file, lineNumber, ok := runtime.Caller(1)
	if !ok {
		panic("ExpectThat: runtime.Caller")
	}

	// Assemble the user error, if any.
	userError := ""
	if len(errorParts) != 0 {
		v := reflect.ValueOf(errorParts[0])
		if v.Kind() != reflect.String {
			panic(fmt.Sprintf("ExpectThat: invalid format string type %v", v.Kind()))
		}

		userError = fmt.Sprintf(v.String(), errorParts[1:]...)
	}

	// Grab the current test state.
	state := currentlyRunningTest
	if state == nil {
		panic("ExpectThat: no test state.")
	}

	// Check whether the value matches.
	res, matcherErr := m.Matches(x)
	switch res {
	// Return immediately on success.
	case oglematchers.MATCH_TRUE:
		return &expectationModifierImpl{}

	// Handle errors below.
	case oglematchers.MATCH_FALSE:
	case oglematchers.MATCH_UNDEFINED:

	// Panic for invalid results.
	default:
		panic(fmt.Sprintf("ExpectThat: invalid matcher result %v.", res))
	}

	// Form an appropriate failure message. Make sure that the expected and
	// actual values align properly.
	var record failureRecord
	relativeClause := ""
	if matcherErr != nil {
		relativeClause = fmt.Sprintf(", %s", matcherErr.Error())
	}

	record.GeneratedError = fmt.Sprintf(
		"Expected: %s\nActual:   %v%s",
		m.Description(),
		x,
		relativeClause)

	// Record additional failure info.
	record.FileName = path.Base(file)
	record.LineNumber = lineNumber
	record.UserError = userError

	// Store the failure.
	state.FailureRecords = append(state.FailureRecords, &record)

	return &expectationModifierImpl{&record}
}

type expectationModifierImpl struct {
	// The failure record created by the expectation, or nil if none.
	failureRecord *failureRecord
}

func (m *expectationModifierImpl) SetCaller(fileName string, lineNumber int) ExpectationModifier {
	if m.failureRecord == nil {
		// Do nothing.
		return m
	}

	m.failureRecord.FileName = fileName
	m.failureRecord.LineNumber = lineNumber
	return m
}
