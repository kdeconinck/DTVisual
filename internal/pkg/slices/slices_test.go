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

// Quality assurance: Verify (and measure the performance) of the public API of the "slices" package.
package slices_test

import (
	"testing"

	"github.com/kdeconinck/dtvisual/internal/pkg/assert"
	"github.com/kdeconinck/dtvisual/internal/pkg/slices"
)

// UT: Verify if a slice contains an element.
func TestContains(t *testing.T) {
	for _, tc := range []struct {
		input []int
		v     int
		want  bool
	}{
		{
			input: []int{1, 2, 3},
			v:     1,
			want:  true,
		},
		{
			input: []int{1, 2, 3},
			v:     2,
			want:  true,
		},
		{
			input: []int{1, 2, 3},
			v:     3,
			want:  true,
		},
		{
			input: []int{1, 2, 3},
			v:     0,
			want:  false,
		},
	} {
		// ACT.
		got := slices.Contains(tc.input, tc.v)

		// ASSERT.
		assert.Equal(t, got, tc.want, "", "\n\n"+
			"UT Name:    Verify if a slice contains an element.\n"+
			"Input:      %v\n"+
			"\033[32mExpected:   %t\033[0m\n"+
			"\033[31mActual:     %t\033[0m\n\n", tc.input, tc.want, got)
	}
}

// Benchmark: Verify if a slice contains an element.
func BenchmarkContains_LargeSet_PositivePath(b *testing.B) {
	benchmarkContains(make([][4 * 1024]byte, 4096), [4 * 1024]byte{0}, b)
}

// Benchmark: Verify if a slice contains an element.
func BenchmarkContains_LargeSet_NegativePath(b *testing.B) {
	benchmarkContains(make([][4 * 1024]byte, 4096), [4 * 1024]byte{1}, b)
}

// Benchmark: Verify if a slice contains an element.
func benchmarkContains(input [][4096]byte, v [4096]byte, b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = slices.Contains(input, v)
	}
}
