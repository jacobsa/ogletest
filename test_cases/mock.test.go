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

package oglematchers_test

import (
	"image/color"
	"testing"

	. "github.com/jacobsa/oglematchers"
	"github.com/jacobsa/oglemock"
	"github.com/jacobsa/ogletest"
	"github.com/jacobsa/ogletest/test_cases/mock_image"
)

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

type MockTest struct {
	image mock_image.MockImage
}

func init()                     { ogletest.RegisterTestSuite(&MockTest{}) }
func TestMockTest(t *testing.T) { ogletest.RunTests(t) }

func (s *MockTest) SetUp(t *ogletest.T) {
	s.image = mock_image.NewMockImage(t.MockController, "some mock image")
}

////////////////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////////////////

func (s *MockTest) ExpectationSatisfied(t *ogletest.T) {
	t.ExpectCall(s.image, "At")(11, GreaterThan(19)).
		WillOnce(oglemock.Return(color.Gray{0}))

	t.ExpectThat(s.image.At(11, 23), IdenticalTo(color.Gray{0}))
}

func (s *MockTest) MockExpectationNotSatisfied(t *ogletest.T) {
	t.ExpectCall(s.image, "At")(11, GreaterThan(19)).
		WillOnce(oglemock.Return(color.Gray{0}))
}

func (s *MockTest) ExpectCallForUnknownMethod(t *ogletest.T) {
	t.ExpectCall(s.image, "FooBar")(11)
}

func (s *MockTest) UnexpectedCall(t *ogletest.T) {
	s.image.At(11, 23)
}

func (s *MockTest) InvokeFunction(t *ogletest.T) {
	var suppliedX, suppliedY int
	f := func(x, y int) color.Color {
		suppliedX = x
		suppliedY = y
		return color.Gray{17}
	}

	t.ExpectCall(s.image, "At")(Any(), Any()).
		WillOnce(oglemock.Invoke(f))

	t.ExpectThat(s.image.At(-1, 12), IdenticalTo(color.Gray{17}))
	t.ExpectEq(-1, suppliedX)
	t.ExpectEq(12, suppliedY)
}
