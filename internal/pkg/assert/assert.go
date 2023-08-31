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

// Package assert defines functions for making assertions in Go's standard testing framework.
package assert

import "testing"

// Nil compares got against nil.
// If they are NOT equal, t is marked as failed, and it's execution is terminated.
func Nil(t testing.TB, got any, name string, msg ...any) {
	if got != nil {
		t.Helper()

		failT(t, got, nil, name, "%s = %v, want %v", msg...)
	}
}

// NotNil compares got against nil.
// If they are equal, t is marked as failed, and it's execution is terminated.
func NotNil(t testing.TB, got any, name string, msg ...any) {
	if got == nil {
		t.Helper()

		failT(t, got, "NOT <nil>", name, "%s = %v, want %s", msg...)
	}
}

// Equal compares got against want for equality.
// If they are not equal, t is marked as failed, and it's execution is terminated.
func Equal[V comparable](t testing.TB, got, want V, name string, msg ...any) {
	if got != want {
		t.Helper()

		failT(t, got, want, name, "%s = %v, want %v", msg...)
	}
}

// EqualFn compares got against want for equality using the custom comparison function cmpFn.
// If they are not equal, t is marked as failed, and it's execution is terminated.
func EqualFn[V any](t testing.TB, got, want V, cmpFn func(got, want V) bool, name string, msg ...any) {
	if !cmpFn(got, want) {
		t.Helper()

		failT(t, got, want, name, "%s = %v, want %v", msg...)
	}
}

// Marks t as failed and terminates its execution.
func failT[V any](t testing.TB, got, want V, name, msgTemplate string, msg ...any) {
	if name != "" {
		t.Fatalf(msgTemplate, name, got, want)
	} else {
		t.Fatalf(msg[0].(string), msg[1:]...)
	}
}
