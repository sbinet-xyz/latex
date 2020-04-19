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
		{
			node: HListOf([]Node{
				VBox(10, 20), HBox(30),
				HListOf([]Node{HBox(11), HBox(22)}, false),
			}, false),
			w: 63,
			h: 10,
			d: 20,
		},
		{
			node: HListOf([]Node{
				VBox(10, 20), HBox(30),
				HListOf([]Node{HBox(11), HBox(22)}, false),
				VListOf([]Node{HBox(15), VBox(11, 22)}),
			}, false),
			w: 78,
			h: 11,
			d: 22,
		},
		{
			node: HListOf([]Node{VBox(10, 20), NewKern(15), HBox(30)}, true),
			w:    45,
			h:    10,
			d:    20,
		},
		{
			node: VListOf([]Node{
				VBox(10, 20),
				VListOf([]Node{
					VBox(11, 22),
					NewKern(10),
					HBox(40),
				}),
				HListOf([]Node{VBox(10, 20), NewKern(15), HBox(30)}, true),
				HBox(30),
			}),
			w: 45,
			h: 103,
			d: 0,
		},
		{
			node: NewKern(10),
			w:    10,
			h:    0,
			d:    0,
		},
		{
			node: NewGlue("fil"),
		},
		{
			node: NewGlue("fill"),
		},
		{
			node: NewGlue("filll"),
		},
		{
			node: NewGlue("neg_fil"),
		},
		{
			node: NewGlue("neg_fill"),
		},
		{
			node: NewGlue("neg_filll"),
		},
		{
			node: NewGlue("empty"),
		},
		{
			node: NewGlue("ss"),
		},
		{
			node: VListOf([]Node{
				NewKern(10),
				VBox(10, 20),
				NewKern(10),
				VListOf([]Node{
					NewKern(10),
					VBox(11, 22),
					NewKern(10),
					HBox(40),
					NewKern(10),
				}),
				NewKern(10),
				HListOf([]Node{
					NewKern(10), VBox(10, 20),
					NewKern(15), HBox(30),
					NewKern(10),
				}, true),
				NewKern(10),
				HBox(30),
				NewKern(10),
			}),
			w: 65,
			h: 173,
			d: 0,
		},
		{
			node: VListOf([]Node{
				NewKern(10),
				VBox(10, 20),
				NewGlue("fill"),
				NewKern(10),
				VListOf([]Node{
					NewKern(10),
					VBox(11, 22),
					NewKern(10),
					NewGlue("neg_fill"),
					HBox(40),
					NewKern(10),
				}),
				NewKern(10),
				HListOf([]Node{
					NewKern(10), VBox(10, 20),
					NewGlue("empty"),
					NewKern(15), HBox(30),
					NewKern(10),
				}, true),
				NewKern(10),
				NewGlue("ss"),
				HBox(30),
				NewKern(10),
			}),
			w: 65,
			h: 173,
			d: 0,
		},
		{
			node: HListOf([]Node{
				NewKern(10),
				NewGlue("fil"),
				VBox(10, 20), HBox(30),
				NewGlue("fil"),
				HListOf([]Node{HBox(11), NewGlue("filll"), HBox(22)}, true),
				VListOf([]Node{HBox(15), NewGlue("neg_filll"), VBox(11, 22)}),
			}, true),
			w: 88,
			h: 11,
			d: 22,
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
