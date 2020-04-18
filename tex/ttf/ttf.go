// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ttf provides a truetype font tex Backend
package ttf

import (
	"fmt"
	"io/ioutil"
	"unicode"

	"github.com/go-latex/latex/internal/tex2unicode"
	"github.com/go-latex/latex/tex"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Backend struct {
	glyphs map[ttfKey]ttfVal
	fonts  map[string]*truetype.Font
}

func NewBackend() *Backend {
	ttf := &Backend{
		glyphs: make(map[ttfKey]ttfVal),
		fonts:  make(map[string]*truetype.Font),
	}

	ftmap := map[string]string{
		"default": "/usr/lib/python3.8/site-packages/matplotlib/mpl-data/fonts/ttf/DejaVuSans.ttf",
		"regular": "/usr/lib/python3.8/site-packages/matplotlib/mpl-data/fonts/ttf/DejaVuSans.ttf",
		"rm":      "/usr/lib/python3.8/site-packages/matplotlib/mpl-data/fonts/ttf/DejaVuSans.ttf",
		"it":      "/usr/lib/python3.8/site-packages/matplotlib/mpl-data/fonts/ttf/DejaVuSans-Oblique.ttf",
	}

	for k, fname := range ftmap {
		raw, err := ioutil.ReadFile(fname)
		if err != nil {
			panic(err)
		}
		ft, err := truetype.Parse(raw)
		if err != nil {
			panic(err)
		}
		ttf.fonts[k] = ft
	}

	return ttf
}

// RenderGlyphs renders the glyph g at the reference point (x,y).
func (ttf *Backend) RenderGlyph(x, y float64, font tex.Font, symbol string, dpi float64) {
	panic("not implemented")
}

// RenderRectFilled draws a filled black rectangle from (x1,y1) to (x2,y2).
func (ttf *Backend) RenderRectFilled(x1, y1, x2, y2 float64) {
	panic("not implemented")
}

// Metrics returns the metrics.
func (ttf *Backend) Metrics(symbol string, font tex.Font, dpi float64, math bool) tex.Metrics {
	key := ttfKey{symbol, font, dpi}
	val, ok := ttf.glyphs[key]
	if ok {
		return val.metrics
	}

	hinting := hintingNone
	ft, rn, symbol, fontSize, slanted := ttf.getGlyph(symbol, font, math)

	var (
		postscript = ft.Name(truetype.NameIDPostscriptName)
		idx        = ft.Index(rn)
	)

	fupe := fixed.Int26_6(ft.FUnitsPerEm())
	var glyph truetype.GlyphBuf
	err := glyph.Load(ft, fupe, idx, hinting)
	if err != nil {
		panic(err)
	}

	face := truetype.NewFace(ft, &truetype.Options{
		DPI:     72,
		Size:    fontSize,
		Hinting: hinting,
	})
	defer face.Close()

	bds, _ /*adv*/, ok := face.GlyphBounds(rn)
	if !ok {
		panic(fmt.Errorf("could not load glyph bounds for %q", rn))
	}

	dy := bds.Max.Y + bds.Min.Y
	xmin := float64(bds.Min.X) / 64
	ymin := float64(bds.Min.Y-dy) / 64
	xmax := float64(bds.Max.X) / 64
	ymax := float64(bds.Max.Y-dy) / 64
	width := xmax - xmin
	height := ymax - ymin
	hme := ft.HMetric(fupe*64*6, idx)
	adv := hme.AdvanceWidth

	offset := 0.0
	if postscript == "Cmex10" {
		offset = height/2 + (font.Size / 3 * dpi / 72)
	}

	me := tex.Metrics{
		Advance: float64(adv) / 65536 * font.Size / 12,
		Height:  height,
		Width:   width,
		XMin:    xmin,
		XMax:    xmax,
		YMin:    ymin + offset,
		YMax:    ymax + offset,
		Iceberg: ymax + offset,
		Slanted: slanted,
	}

	ttf.glyphs[key] = ttfVal{
		font:       ft,
		size:       font.Size,
		postscript: postscript,
		metrics:    me,
		symbolName: symbol,
		rune:       rn,
		glyph:      idx,
		offset:     offset,
	}
	return me
}

const (
	hintingNone = font.HintingNone
	hintingFull = font.HintingFull
)

func (ttf *Backend) getGlyph(symbol string, font tex.Font, math bool) (*truetype.Font, rune, string, float64, bool) {
	var (
		fontType = font.Type
		idx      = tex2unicode.Index(symbol, math)
	)

	// only characters in the "Letter" class should be italicized in "it" mode.
	// Greek capital letters should be roman.
	if font.Type == "it" && idx < 0x10000 {
		if !unicode.Is(unicode.L, idx) {
			fontType = "rm"
		}
	}
	slanted := (fontType == "it") || ttf.isSlanted(symbol)
	ft := ttf.getFont(fontType)
	if ft == nil {
		panic("could not find TTF font for [" + fontType + "]")
	}

	// FIXME(sbinet):
	// \sigma -> sigma, A->A, \infty->infinity, \nabla->gradient
	// etc...
	symbolName := symbol
	return ft, idx, symbolName, font.Size, slanted
}

func (ttf *Backend) isSlanted(symbol string) bool {
	switch symbol {
	case `\int`, `\oint`:
		return true
	default:
		return false
	}
}

func (ttf *Backend) getFont(fontType string) *truetype.Font {
	return ttf.fonts[fontType]
}

// UnderlineThickness returns the line thickness that matches the given font.
// It is used as a base unit for drawing lines such as in a fraction or radical.
func (ttf *Backend) UnderlineThickness(font tex.Font, dpi float64) float64 {
	// theoretically, we could grab the underline thickness from the font
	// metrics.
	// but that information is just too un-reliable.
	// so, it is hardcoded.
	return (0.75 / 12 * font.Size * dpi) / 72
}

// Kern returns the kerning distance between two symbols.
func (ttf *Backend) Kern(ft1 tex.Font, sym1 string, ft2 tex.Font, sym2 string, dpi float64) float64 {
	panic("not implemented")
}

type ttfKey struct {
	symbol string
	font   tex.Font
	dpi    float64
}

type ttfVal struct {
	font       *truetype.Font
	size       float64
	postscript string
	metrics    tex.Metrics
	symbolName string
	rune       rune
	glyph      truetype.Index
	offset     float64
}

var (
	_ tex.Backend = (*Backend)(nil)
)