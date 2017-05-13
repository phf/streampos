// Copyright 2016 Peter H. Froehlich. All rights reserved.
// Use of this source code is governed by the MIT license,
// see the LICENSE.md file.

package streampos

import (
	"testing"
)

type want struct {
	offset       int64 // for this
	line, column int64 // we want this
}

type testcase struct {
	data  []string // stream contents
	wants []want   // correct results
}

var testcases = []testcase{
	{
		[]string{"peter\npaul\nand mary"},
		[]want{
			{0, 1, 1},
			{4, 1, 5},
			{5, 1, 6},
			{6, 2, 1},
			{10, 2, 5},
			{11, 3, 1},
			{18, 3, 8},
			{19, -1, -1},
		},
	},
	{
		[]string{""},
		[]want{
			{0, -1, -1},
		},
	},
	{
		[]string{"\n\n\n"},
		[]want{
			{0, 1, 1},
			{1, 2, 1},
			{2, 3, 1},
			{3, -1, -1},
		},
	},
	{
		[]string{"\n", "\n", "\n"},
		[]want{
			{0, 1, 1},
			{1, 2, 1},
			{2, 3, 1},
			{3, -1, -1},
		},
	},
}

func TestAll(t *testing.T) {
	for _, tc := range testcases {
		w := &Writer{}
		total := int64(0)
		for _, d := range tc.data {
			if total != w.Length() {
				t.Errorf("%v, want %v", w.Length(), total)
			}
			n, err := w.Write([]byte(d))
			if n != len(d) || err != nil {
				t.Errorf("(%v, %v), want (%v, %v)", n, err, len(d), nil)
			}
			total += int64(n)
		}
		for _, wa := range tc.wants {
			line, _ := w.Line(wa.offset)
			column, _ := w.Column(wa.offset)
			if line != wa.line || column != wa.column {
				t.Errorf("(%v, %v), want (%v, %v)", line, column, wa.line, wa.column)
			}
		}
	}
}
