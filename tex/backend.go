// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tex

// Backend is the interface that allows to render math expressions.
type Backend interface {
	// RenderGlyphs renders the glyph g at the reference point (x,y).
	RenderGlyph(x, y float64)

	// RenderRectFilled draws a filled black rectangle from (x1,y1) to (x2,y2).
	RenderRectFilled(x1, y1, x2, y2 float64)
}
