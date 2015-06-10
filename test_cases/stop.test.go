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

package oglematchers_test

import (
	"testing"

	"github.com/jacobsa/ogletest"
)

func TestStop(t *testing.T) { ogletest.RunTests(t) }

////////////////////////////////////////////////////////////////////////
// Boilerplate
////////////////////////////////////////////////////////////////////////

type StopTest struct {
}

var _ ogletest.TearDownInterface = &StopTest{}

func init() { ogletest.RegisterTestSuite(&StopTest{}) }

func (s *StopTest) TearDown(t *ogletest.T) {
	t.Logf("TearDown running.")
}

////////////////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////////////////

func (s *StopTest) First(t *ogletest.T) {
}

func (s *StopTest) Second(t *ogletest.T) {
	t.Logf("About to call StopRunningTests.")
	ogletest.StopRunningTests()
	t.Logf("Called StopRunningTests.")
}

func (s *StopTest) Third(t *ogletest.T) {
}
