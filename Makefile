include $(GOROOT)/src/Make.inc

TARG = github.com/jacobsa/ogletest
GOFILES = \
	expect_that.go \
	register_test_suite.go \
	run_tests.go \

include $(GOROOT)/src/Make.pkg
