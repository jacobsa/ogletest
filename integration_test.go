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
	"go/build"
	"io/ioutil"
	"path"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

var dumpNew = flag.Bool("dump_new", false, "Dump new golden files.")
var objDir string

////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////

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
		panic("ioutil.WriteFile: " + err.Error())
	}
}

func readFileOrDie(path string) []byte {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		panic("ioutil.ReadFile: " + err.Error())
	}

	return contents
}

// cleanOutput transforms the supplied output so that it no longer contains
// information that changes from run to run, making the golden tests less
// flaky.
func cleanOutput(o []byte, testDir string) []byte {
	// Replace references to the test directory.
	o = []byte(strings.Replace(string(o), testDir, "/tmp/dir", -1))

	// Replace things that look like line numbers and process counters in stack
	// traces.
	stackFrameRe := regexp.MustCompile(`\S+\.go:\d+ \(0x[0-9a-f]+\)`)
	o = stackFrameRe.ReplaceAll(o, []byte("some_file.go:0 (0x00000)"))

	// Replace unstable timings in gotest fail messages.
	timingRe := regexp.MustCompile(`--- FAIL: .* \(\d\.\d{2} seconds\)`)
	o = timingRe.ReplaceAll(o, []byte("--- FAIL: sometest (0.00 seconds)"))

	return o
}

// Create a temporary package directory somewhere that 'go test' can find, and
// return the directory and package name.
func createTempPackageDir(caseName string) (dir, pkg string) {
	const ogletestPkg = "github.com/jacobsa/ogletest"

	// Figure out where the local source code for ogletest is.
	tree, _, err := build.FindTree(ogletestPkg)
	if err != nil { panic("Finding ogletest tree: " + err.Error()) }

	// Create a temporary directory underneath this.
	ogletestPkgDir := path.Join(tree.Path, "src", ogletestPkg)
	prefix := fmt.Sprintf("tmp-%s-", caseName)

	dir, err = ioutil.TempDir(ogletestPkgDir, prefix)
	if err != nil { panic("ioutil.TempDir: " + err.Error()) }

	pkg = path.Join("github.com/jacobsa/ogletest", dir[len(ogletestPkgDir):])
	return
}

// runTestCase runs the case with the supplied name (e.g. "passing_test"), and
// returns its output and exit code.
func runTestCase(name string) ([]byte, int, error) {
	// Create a temporary directory for the test files.
	testDir, testPkg := createTempPackageDir(name)
	 defer os.RemoveAll(testDir)

	// Create the test source file.
	sourceFile := name + ".go"
	testContents := readFileOrDie(path.Join("test_cases", sourceFile))
	writeContentsToFileOrDie(testContents, path.Join(testDir, sourceFile))

	// Invoke 'go test'. Special case: pass a test filter to the filtered_test
	// case.
	cmd := exec.Command("go", "test", testPkg)
	if name == "filtered_test" {
		cmd.Args = append(cmd.Args, "--ogletest.run=Test(Bar|Baz)")
	}

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

	output = cleanOutput(output, testDir)
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

		// Check the status code. We assume all test cases fail except for
		// passing_test.
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
