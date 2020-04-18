// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tex

import "testing"

func TestTex2Unicode(t *testing.T) {
	for _, tc := range []struct {
		v    string
		want rune
		math bool
	}{
		{v: `a`, want: 'a'},
		{v: `a`, want: 'a', math: true},
		{v: `A`, want: 'A'},
		{v: `A`, want: 'A', math: true},
		{v: `0`, want: '0'},
		{v: `0`, want: '0', math: true},
		{v: `-`, want: '-'},
		{v: `\alpha`, want: 'α', math: true},
		{v: `\Join`, want: '⨝', math: true},
		{v: `\perp`, want: '⟂', math: true},
		{v: `\pm`, want: '±', math: true},
		{v: `\mp`, want: '∓', math: true},
		{v: `\neq`, want: '≠', math: true},
		{v: `\__sqrt__`, want: '√', math: true},
		{v: `\partial`, want: '∂', math: true},
		{v: `\hbar`, want: 'ħ', math: true},
		{v: `\hslash`, want: 'ℏ', math: true},
		{v: `\int`, want: '∫', math: true},
		{v: `\oint`, want: '∮', math: true},
		{v: `\oiint`, want: '∯', math: true},
		{v: `\infty`, want: '∞', math: true},
		{v: `\sigma`, want: 'σ', math: true},
		{v: `\varsigma`, want: 'ς', math: true},
		{v: `\Sigma`, want: 'Σ', math: true},
		{v: `\sum`, want: '∑', math: true},
		{v: `\varepsilon`, want: 'ε', math: true},
	} {
		t.Run(tc.v, func(t *testing.T) {
			got := unicodeIndex(tc.v, tc.math)
			if got != tc.want {
				t.Fatalf("error: got=%q, want=%q", got, tc.want)
			}
		})
	}
}
