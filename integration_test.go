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
	"flag"
	"fmt"
	"io/ioutil"
	"path"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var dumpNew = flag.Bool("dump_new", false, "Dump new golden files.")
var objDir string

////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////

// buildPackage build a copy of the (possibly locally modified) package, and
// return a path to an include path where it can be found.
func buildPackage() (string, error) {
	// HACK(jacobsa): The gotest command reads the GCFLAGS environment variable,
	// which allows us to give a -I flag to 6g to let it find the locally built
	// package. However 6l cannot find it, and there is no linker flags
	// environment variable. So we actually install the package here, instead of
	// just building it locally.
	//
	// TODO(jacobsa): When the new 'go' tool comes out, check whether this still
	// replies to its 'test' command. If so, file an issue to get a linker flags
	// environment variable added.
	cmd := exec.Command("gomake", "install")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", errors.New(fmt.Sprintf("%v:\n%s", err, output))
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", errors.New("os.Getwd: " + err.Error())
	}

	return path.Join(wd, "_obj"), nil
}

// getCaseNames looks for integration test cases as files in the test_cases
// directory.
func getCaseNames() ([]string, error) {
	// Open the test cases directory.
	dir, err := os.Open("test_cases")
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
		// Skip golden files and hidden files.
		if strings.HasPrefix(name, "golden.") || strings.HasPrefix(name, ".") {
			continue
		}

		// Check for the right format.
		if !strings.HasSuffix(name, "_test.go") {
			return nil, errors.New(fmt.Sprintf("Unexpected file: %s", name))
		}

		// Store the name minus the extension.
		result[resultLen] = name[:len(name) - 3]
		resultLen++
	}

	return result[:resultLen], nil
}

func writeContentsToFileOrDie(contents []byte, path string) {
	if err := ioutil.WriteFile(path, contents, 0600); err != nil {
		panic("iotuil.WriteFile: " + err.Error())
	}
}

func readFileOrDie(path string) []byte {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		panic("iotuil.ReadFile: " + err.Error())
	}

	return contents
}

// cleanOutput transforms the supplied output so that it no longer contains
// information that changes from run to run, making the golden tests less
// flaky.
func cleanOutput(o []byte, tempDir string) []byte {
	// HACK(jacobsa): Replace references to the temporary directory in the
	// output, since they differ for each run.
	o = []byte(strings.Replace(string(o), tempDir, "/tmp/dir", -1))

	return o
}

// runTestCase runs the case with the supplied name (e.g. "passing_test"), and
// returns its output and exit code.
func runTestCase(name string) ([]byte, int, error) {
	// Create a temporary directory for the test files.
	tempDir, err := ioutil.TempDir("", "ogletest_integration_test")
	if err != nil {
		return nil, 0, errors.New("ioutil.TempDir: " + err.Error())
	}

	defer os.RemoveAll(tempDir)

	// Create a makefile within the directory.
	makefileContents := "include $(GOROOT)/src/Make.inc\n" +
		"TARG = github.com/jacobsa/ogletest/foobar\n" +
		"GOFILES = " + name + ".go\n" +
		"GCFLAGS += -I" + objDir + "\n" +
		"include $(GOROOT)/src/Make.pkg\n"

	writeContentsToFileOrDie([]byte(makefileContents), path.Join(tempDir, "Makefile"))

	// Create the test source file.
	sourceFile := name + ".go"
	testContents := readFileOrDie(path.Join("test_cases", sourceFile))
	writeContentsToFileOrDie(testContents, path.Join(tempDir, sourceFile))

	// Invoke gotest. Special case: pass a test filter to the filtered_test case.
	cmd := exec.Command("gotest")
	if name == "filtered_test" {
		cmd.Args = append(cmd.Args, "--ogletest.run=Test(Bar|Baz)")
	}

	cmd.Dir = tempDir
	output, err := cmd.CombinedOutput()

	// Did the process exist with zero code?
	if err == nil {
		return output, 0, nil
	}

	// Make sure the process actually exited.
	exitError, ok := err.(*exec.ExitError)
	if !ok || !exitError.Exited() {
		return nil, 0, errors.New("exec.Command.Output: " + err.Error())
	}

	output = cleanOutput(output, tempDir)
	return output, exitError.ExitStatus(), nil
}

// checkGolden file checks the supplied actual output for the named test case
// against the golden file for that case. If requested by the user, it rewrites
// the golden file on failure.
func checkAgainstGoldenFile(caseName string, output []byte) bool {
	goldenFile := path.Join("test_cases", "golden." + caseName)
	goldenContents := readFileOrDie(goldenFile)

	result := string(output) == string(goldenContents)
	if !result && *dumpNew {
		writeContentsToFileOrDie(output, goldenFile)
	}

	return result
}

////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////

func TestGoldenFiles(t *testing.T) {
	// Ensure the local package is up to date, and get a path to it. This will
	// prevent the test cases from using the installed version, which may be out
	// of date.
	var err error
	objDir, err = buildPackage()

	// We expect there to be at least one case.
	caseNames, err := getCaseNames()
	if err != nil || len(caseNames) == 0 {
		t.Fatalf("Error getting cases: %v", err)
	}

	// Run each test case.
	for _, caseName := range caseNames {
		// Run the test case.
		output, exitCode, err := runTestCase(caseName)
		if err != nil {
			t.Fatalf("Running test case %s: %v", caseName, err)
		}

		// Check the status code. We assume all test cases fail except for passing_test.
		shouldPass := caseName == "passing_test"
		didPass := exitCode == 0
		if shouldPass != didPass {
			t.Errorf("Bad exit code for test case %s: %d", caseName, exitCode)
		}

		// Check the output against the golden file.
		if !checkAgainstGoldenFile(caseName, output) {
			t.Errorf("Output for test case %s doesn't match golden file.", caseName)
		}
	}
}
