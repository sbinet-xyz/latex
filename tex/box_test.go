// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tex

import (
	"io/ioutil"
	"testing"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
)

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
		{
			node: HCentered([]Node{
				VBox(10, 20),
				HBox(30),
				NewKern(15),
				HBox(40),
				VBox(20, 10),
			}),
			w: 85,
			h: 20,
			d: 20,
		},
		{
			node: VCentered([]Node{
				VBox(10, 20),
				HBox(30),
				NewKern(15),
				HBox(40),
				VBox(20, 10),
			}),
			w: 40,
			h: 75,
			d: 0,
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

func TestFont(t *testing.T) {
	//fname := "/usr/share/fonts/TTF/DejaVuSans.ttf"
	fname := "/usr/lib/python3.8/site-packages/matplotlib/mpl-data/fonts/ttf/DejaVuSans.ttf"
	//fname := "/usr/lib/python3.8/site-packages/matplotlib/mpl-data/fonts/ttf/DejaVuSansDisplay.ttf"
	ttf, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	ft, err := truetype.Parse(ttf)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	face := truetype.NewFace(ft, &truetype.Options{
		Size: 10,
		DPI:  72,
		//Hinting: font.HintingFull,
	})
	defer face.Close()

	t.Logf("metrics: %#v", face.Metrics())
	bounds, advance, ok := face.GlyphBounds('A')
	if !ok {
		t.Fatalf("could not find glyph bounds")
	}
	t.Logf("box: %#v, advance=%#v, ok=%v", bounds, advance, ok)
	dr, _, _, a, ok := face.Glyph(fixed.P(0, 0), 'A')
	t.Logf("dr: %#v, adv=%v, ok=%v", dr, float64(a), ok)
	b2 := ft.Bounds(fixed.Int26_6(ft.FUnitsPerEm()))
	scale := 10 / float64(ft.FUnitsPerEm())
	t.Logf("scale: %v", scale)
	t.Logf("fu/em: %v", ft.FUnitsPerEm())
	t.Logf("b2: %#v w:%v h:%v", b2, float64(b2.Max.X-b2.Min.X)*scale, scale*float64(b2.Max.Y-b2.Min.Y))

	adv, ok := face.GlyphAdvance('A')
	t.Logf("adv=%#v, ok=%v", adv, ok)
	t.Fatalf("bounds=%#v, advance=%#v", bounds, advance)
}
