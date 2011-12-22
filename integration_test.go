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

package ogletest_test

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
)

////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////

func getCaseNames() ([]string, error) {
	// Open the test cases directory.
	dir, err := os.Open("integration_test_cases")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Opening dir: %v", err))
	}

	// Get a list of the names in the directory.
	names, err := dir.Readdirnames(0)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Readdirnames: %v", err))
	}

	// Filter the names.
	result := make([]string, len(names))
	resultLen := 0
	for _, name := range names {
		// Skip golden files.
		if strings.HasPrefix(name, "golden.") {
			continue
		}

		// Check for the right format.
		if !strings.HasSuffix(name, "_test.go") {
			return nil, errors.New(fmt.Sprintf("Unexpected file: %s", name))
		}

		// Store the name minus the extension.
		result[resultLen] = name[len(name) - 2:]
	}

	return result, nil
}

////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////

func TestGoldenFiles(t *testing.T) {
	// We expect there to be at least one case.
	caseNames, err := getCaseNames()
	if err != nil || len(caseNames) == 0 {
		t.Fatalf("Error getting cases: %v", err)
	}
}
