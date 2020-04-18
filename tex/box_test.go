// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tex

import "testing"

func TestBox(t *testing.T) {
	for _, tc := range []struct {
		node    Node
		w, h, d float64
	}{
		{
			node: HBox(10),
			w:    10,
			h:    0,
			d:    0,
		},
		{
			node: VBox(10, 20),
			w:    0,
			h:    10,
			d:    20,
		},
		{
			node: HListOf([]Node{VBox(10, 20), HBox(30)}, false),
			w:    30,
			h:    10,
			d:    20,
		},
		{
			node: HListOf([]Node{VBox(10, 20), HBox(30)}, true),
			w:    30,
			h:    10,
			d:    20,
		},
		{
			node: VListOf([]Node{VBox(10, 20), HBox(30)}),
			w:    30,
			h:    30,
			d:    0,
		},
	} {
		t.Run("", func(t *testing.T) {
			var (
				w = tc.node.Width()
				h = tc.node.Height()
				d = tc.node.Depth()
			)

			if got, want := w, tc.w; got != want {
				t.Fatalf("invalid width: got=%g, want=%g", got, want)
			}

			if got, want := h, tc.h; got != want {
				t.Fatalf("invalid height: got=%g, want=%g", got, want)
			}

			if got, want := d, tc.d; got != want {
				t.Fatalf("invalid depth: got=%g, want=%g", got, want)
			}
		})
	}
}