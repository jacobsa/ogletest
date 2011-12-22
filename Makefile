include $(GOROOT)/src/Make.inc

TARG = github.com/jacobsa/ogletest
GOFILES = \
	expect_aliases.go \
	expect_that.go \
	register_test_suite.go \
	run_tests.go \
	test_state.go \

include $(GOROOT)/src/Make.pkg
