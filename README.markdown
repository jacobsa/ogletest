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
`weekly/weekly.2011-12-14`. See [here][golang-install] for instructions. Until
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


[reference]: http://gopkgdoc.appspot.com/pkg/github.com/jacobsa/ogletest
[matcher-reference]: http://gopkgdoc.appspot.com/pkg/github.com/jacobsa/oglematchers
[golang-install]: http://golang.org/doc/install.html#releases
[googletest]: http://code.google.com/p/googletest/
[google-js-test]: http://code.google.com/p/google-js-test/
