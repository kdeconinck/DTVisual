// =====================================================================================================================
// = LICENSE:       Copyright (c) 2023 Kevin De Coninck
// =
// =                Permission is hereby granted, free of charge, to any person
// =                obtaining a copy of this software and associated documentation
// =                files (the "Software"), to deal in the Software without
// =                restriction, including without limitation the rights to use,
// =                copy, modify, merge, publish, distribute, sublicense, and/or sell
// =                copies of the Software, and to permit persons to whom the
// =                Software is furnished to do so, subject to the following
// =                conditions:
// =
// =                The above copyright notice and this permission notice shall be
// =                included in all copies or substantial portions of the Software.
// =
// =                THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// =                EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// =                OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// =                NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// =                HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// =                WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// =                FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// =                OTHER DEALINGS IN THE SOFTWARE.
// =====================================================================================================================

// Quality assurance: Verify (and measure the performance) of the public API of the "assert" package.
package assert_test

import (
	"fmt"
	"testing"

	"github.com/kdeconinck/dtvisual/internal/pkg/assert"
)

// The testableT wraps the testing.T struct and adds a field for storing the failure message.
type testableT struct {
	testing.TB
	failureMsg string
}

// Fatal formats args using fmt.Sprintf and stores the result in t.
func (t *testableT) Fatalf(format string, args ...any) {
	t.failureMsg = fmt.Sprintf(format, args...)
}

// UT: Compare a value against nil.
func TestNil(t *testing.T) {
	t.Parallel() // Enable parallel execution.

	for _, tc := range []struct {
		g    any
		name string
		want string
	}{
		{
			g:    true,
			name: "ValueOf(true)",
			want: "ValueOf(true) = true, want <nil>",
		},
		{
			name: "ValueOf(nil)",
		},
	} {
		// ARRANGE.
		testingT := &testableT{TB: t}

		// ACT.
		assert.Nil(testingT, tc.g, tc.name)

		// ASSERT.
		if testingT.failureMsg != tc.want {
			t.Fatalf("Failure message = \"%s\", want \"%s\"", testingT.failureMsg, tc.want)
		}
	}
}

// UT: Compare a value against nil (with a custom message).
func TestNilWithCustomMessage(t *testing.T) {
	t.Parallel() // Enable parallel execution.

	for _, tc := range []struct {
		g    any
		name string
		msg  []any
		want string
	}{
		{
			g:    true,
			name: "",
			msg:  []any{"UT Failed: `ValueOf(true)` - got %t, want <nil>.", true},
			want: "UT Failed: `ValueOf(true)` - got true, want <nil>.",
		},
	} {
		// ARRANGE.
		testingT := &testableT{TB: t}

		// ACT.
		assert.Nil(testingT, tc.g, tc.name, tc.msg...)

		// ASSERT.
		if testingT.failureMsg != tc.want {
			t.Fatalf("Failure message = \"%s\", want \"%s\"", testingT.failureMsg, tc.want)
		}
	}
}

// UT: Compare a value against nil.
func TestNotNil(t *testing.T) {
	t.Parallel() // Enable parallel execution.

	for _, tc := range []struct {
		g    any
		name string
		want string
	}{
		{
			g:    nil,
			name: "ValueOf(nil)",
			want: "ValueOf(nil) = <nil>, want NOT <nil>",
		},
		{
			g:    true,
			name: "ValueOf(true)",
		},
	} {
		// ARRANGE.
		testingT := &testableT{TB: t}

		// ACT.
		assert.NotNil(testingT, tc.g, tc.name)

		// ASSERT.
		if testingT.failureMsg != tc.want {
			t.Fatalf("Failure message = \"%s\", want \"%s\"", testingT.failureMsg, tc.want)
		}
	}
}

// UT: Compare a value against nil (with a custom message).
func TestNotNilWithCustomMessage(t *testing.T) {
	t.Parallel() // Enable parallel execution.

	for _, tc := range []struct {
		g    any
		name string
		msg  []any
		want string
	}{
		{
			g:    nil,
			name: "",
			msg:  []any{"UT Failed: `ValueOf(nil)` - got <nil>, want NOT <nil>."},
			want: "UT Failed: `ValueOf(nil)` - got <nil>, want NOT <nil>.",
		},
	} {
		// ARRANGE.
		testingT := &testableT{TB: t}

		// ACT.
		assert.NotNil(testingT, tc.g, tc.name, tc.msg...)

		// ASSERT.
		if testingT.failureMsg != tc.want {
			t.Fatalf("Failure message = \"%s\", want \"%s\"", testingT.failureMsg, tc.want)
		}
	}
}

// UT: Compare 2 values for equality.
func TestEqual(t *testing.T) {
	t.Parallel() // Enable parallel execution.

	for _, tc := range []struct {
		g, w bool
		name string
		want string
	}{
		{
			g: true, w: false,
			name: "IsDigit(\"0\")",
			want: "IsDigit(\"0\") = true, want false",
		},
		{
			g: true, w: true,
			name: "IsDigit(\"0\")",
		},
	} {
		// ARRANGE.
		testingT := &testableT{TB: t}

		// ACT.
		assert.Equal(testingT, tc.g, tc.w, tc.name)

		// ASSERT.
		if testingT.failureMsg != tc.want {
			t.Fatalf("Failure message = \"%s\", want \"%s\"", testingT.failureMsg, tc.want)
		}
	}
}

// UT: Compare 2 values for equality (with a custom message).
func TestEqualWithCustomMessage(t *testing.T) {
	t.Parallel() // Enable parallel execution.

	for _, tc := range []struct {
		g, w bool
		name string
		msg  []any
		want string
	}{
		{
			g: true, w: false,
			name: "",
			msg:  []any{"UT Failed: `IsDigit(\"0\")` - got %t, want %t.", true, false},
			want: "UT Failed: `IsDigit(\"0\")` - got true, want false.",
		},
	} {
		// ARRANGE.
		testingT := &testableT{TB: t}

		// ACT.
		assert.Equal(testingT, tc.g, tc.w, tc.name, tc.msg...)

		// ASSERT.
		if testingT.failureMsg != tc.want {
			t.Fatalf("Failure message = \"%s\", want \"%s\"", testingT.failureMsg, tc.want)
		}
	}
}

// UT: Compare 2 values for equality using a custom comparison function.
func TestEqualFn(t *testing.T) {
	t.Parallel() // Enable parallel execution.

	for _, tc := range []struct {
		g, w bool
		name string
		want string
	}{
		{
			g: true, w: false,
			name: "IsDigit(\"0\")",
			want: "IsDigit(\"0\") = true, want false",
		},
		{
			g: true, w: true,
			name: "IsDigit(\"0\")",
		},
	} {
		// ARRANGE.
		testingT := &testableT{TB: t}

		// ACT.
		assert.EqualFn(testingT, tc.g, tc.w, func(got, want bool) bool { return got == want }, tc.name)

		// ASSERT.
		if testingT.failureMsg != tc.want {
			t.Fatalf("Failure message = \"%s\", want \"%s\"", testingT.failureMsg, tc.want)
		}
	}
}

// UT: Compare 2 values for equality (with a custom message) using a custom comparison function.
func TestEqualFnWithCustomMessage(t *testing.T) {
	t.Parallel() // Enable parallel execution.

	for _, tc := range []struct {
		g, w bool
		name string
		msg  []any
		want string
	}{
		{
			g: true, w: false,
			name: "",
			msg:  []any{"UT Failed: `IsDigit(\"0\")` - got %t, want %t.", true, false},
			want: "UT Failed: `IsDigit(\"0\")` - got true, want false.",
		},
	} {
		// ARRANGE.
		testingT := &testableT{TB: t}

		// ACT.
		assert.EqualFn(testingT, tc.g, tc.w, func(got, want bool) bool { return got == want }, tc.name, tc.msg...)

		// ASSERT.
		if testingT.failureMsg != tc.want {
			t.Fatalf("Failure message = \"%s\", want \"%s\"", testingT.failureMsg, tc.want)
		}
	}
}
