`ogletest` is a unit testing framework for Go with the following features:

 *  An extensive and extensible set of matchers for expressing expectations.
 *  Automatic failure messages; no need to say `t.Errorf("Expected %v, got
    %v"...)`.
 *  Clean, readable output that tells you exaclty what you need to know.
 *  Style and semantics similar to [Google Test][googletest] and
    [Google JS Test][google-js-test].

It integrates with Go's built-in `testing` package, so it works with the
`gotest` command, and even with other types of test within your package. Unlike
the `testing` package which offers only basic capabilities for signalling
failures, it offers ways to express expectations and get nice failure messages
automatically.


Installation
------------

First, make sure you have installed a version of the Go tools at least as new as
`weekly/weekly.2011-12-22`. See [here][golang-install] for instructions. Until
release `r61` comes out, this involes using the `weekly` tag.

Use the following command to install `ogletest` and its dependencies, and to
keep them up to date:

    goinstall -u github.com/jacobsa/ogletest


Documentation
-------------

See [here][reference] for package documentation hosted on GoPkgDoc containing an
exhaustive list of exported symbols. Alternatively, you can install the package
and then use `godoc`:

    godoc github.com/jacobsa/oglematchers

An important part of `ogletest` is its use of matchers provided by the
[`oglematchers`][matcher-reference] package. See that package's documentation
for information on the built-in matchers available, and check out the
`oglematchers.Matcher` interface if you want to define your own.


Example
-------

Let's say you have a function in your package `people` with the following
signature:

```go
// GetRandomPerson returns the name and phone number of Tony, Dennis, or Scott.
func GetRandomPerson() (name, phone string) {
  [...]
}
```

A silly function, but it will do for an example. You can write a couple of tests
for it as follows:

```go
package people

import (
  . "github.com/jacobsa/oglematchers"
  . "github.com/jacobsa/ogletest"
  "testing"
)

// Give ogletest a chance to run your tests when invoked by gotest.
func TestOgletest(t *testing.T) { RunTests(t) }

// Create a test suite, which groups together logically related test methods
// (defined below). You can share common setup and teardown code here; see the
// package docs for more info.
type PeopleTest struct {}
func init() { RegisterTestSuite(&PeopleTest{}) }

func (t *PeopleTest) ReturnsCorrectNames() {
  // Call the function a few times, and make sure it never strays from the set
  // of expected names.
  for i := 0; i < 25; i++ {
    name, _ := GetRandomPerson()
    ExpectThat(name, AnyOf("Tony", "Dennis", "Scott"))
  }
}

func (t *PeopleTest) FormatsPhoneNumbersCorrectly() {
  // Call the function a few times, and make sure it returns phone numbers in a
  // standard US format.
  for i := 0; i < 25; i++ {
    _, phone := GetRandomPerson()
    ExpectThat(phone, MatchesRegexp(`^\(\d{3}\) \d{3}-\d{4}$`))
  }
}
```

If you save this test in a file whose name ends in `_test.go` and set up a
makefile for your package as described in the [How to Write Go Code][howtowrite]
docs, you can run your tests by simply invoking the following in your package
directory:

    gotest

Here's what the failure output of ogletest looks like, if your function's
implementation is bad.

    [----------] Running tests from PeopleTest
    [ RUN      ] PeopleTest.FormatsPhoneNumbersCorrectly
    people_test.go:32:
    Expected: matches regexp "^\(\d{3}\) \d{3}-\d{4}$"
    Actual:   +1 800 555 5555
    
    [  FAILED  ] PeopleTest.FormatsPhoneNumbersCorrectly
    [ RUN      ] PeopleTest.ReturnsCorrectNames
    people_test.go:23:
    Expected: or(Tony, Dennis, Scott)
    Actual:   Bart
    
    [  FAILED  ] PeopleTest.ReturnsCorrectNames
    [----------] Finished with tests from PeopleTest

And if the test passes:

    [----------] Running tests from PeopleTest
    [ RUN      ] PeopleTest.FormatsPhoneNumbersCorrectly
    [       OK ] PeopleTest.FormatsPhoneNumbersCorrectly
    [ RUN      ] PeopleTest.ReturnsCorrectNames
    [       OK ] PeopleTest.ReturnsCorrectNames
    [----------] Finished with tests from PeopleTest


[reference]: http://gopkgdoc.appspot.com/pkg/github.com/jacobsa/ogletest
[matcher-reference]: http://gopkgdoc.appspot.com/pkg/github.com/jacobsa/oglematchers
[golang-install]: http://golang.org/doc/install.html#releases
[googletest]: http://code.google.com/p/googletest/
[google-js-test]: http://code.google.com/p/google-js-test/
[howtowrite]: http://golang.org/doc/code.html
