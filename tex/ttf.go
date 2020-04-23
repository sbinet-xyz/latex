// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tex

import (
	"fmt"
	"io/ioutil"
	"log"
	"unicode"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

type TTFBackend struct {
	glyphs map[ttfKey]ttfVal
	fonts  map[string]*truetype.Font
	sfnts  map[string]*sfnt.Font
}

func NewTTFBackend() *TTFBackend {
	ttf := &TTFBackend{
		glyphs: make(map[ttfKey]ttfVal),
		fonts:  make(map[string]*truetype.Font),
		sfnts:  make(map[string]*sfnt.Font),
	}

	ftmap := map[string]string{
		"default": "/usr/lib/python3.8/site-packages/matplotlib/mpl-data/fonts/ttf/DejaVuSans.ttf",
		"regular": "/usr/lib/python3.8/site-packages/matplotlib/mpl-data/fonts/ttf/DejaVuSans.ttf",
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

		sft, err := sfnt.Parse(raw)
		if err != nil {
			panic(err)
		}
		ttf.sfnts[k] = sft
	}

	return ttf
}

// RenderGlyphs renders the glyph g at the reference point (x,y).
func (ttf *TTFBackend) RenderGlyph(x, y float64) {
	panic("not implemented")
}

// RenderRectFilled draws a filled black rectangle from (x1,y1) to (x2,y2).
func (ttf *TTFBackend) RenderRectFilled(x1, y1, x2, y2 float64) {
	panic("not implemented")
}

// Metrics returns the metrics.
func (ttf *TTFBackend) Metrics(symbol string, font Font, dpi float64, math bool) Metrics {
	key := ttfKey{symbol, font, dpi}
	val, ok := ttf.glyphs[key]
	if ok {
		return val.metrics
	}

	hinting := hintingNone

	ft, _, rn, symbol, fontSize, slanted := ttf.getGlyph(symbol, font, math)
	idx := ft.Index(rn)

	fupe := fixed.Int26_6(ft.FUnitsPerEm())
	//	fupe = fixed.Int26_6(0.5 + (font.Size * 64))
	var glyph truetype.GlyphBuf
	err := glyph.Load(ft, fupe, idx, hinting)
	if err != nil {
		panic(err)
	}

	//	var sbuf sfnt.Buffer
	//	name, err := sft.Name(&sbuf, sfnt.NameIDPostScript)
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Printf("postscriptname: %q\n", name)
	//
	//	gi, err := sft.GlyphIndex(&sbuf, rn)
	//	if err != nil {
	//		panic(err)
	//	}
	//	const ppem = 32
	//	_, err = sft.LoadGlyph(&sbuf, gi, fixed.I(ppem), nil)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	sadv, err := sft.GlyphAdvance(&sbuf, gi, fixed.I(ppem), hinting)
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Printf("sadv: %#v\n", float64(sadv))
	//	//	sme, err := sft.Metrics(&sbuf, fixed.I(ppem), hinting)
	//	//	if err != nil {
	//	//		panic(err)
	//	//	}
	//	//	fmt.Printf("sme: %#v\n", sme)
	//
	//	sbds, err := sft.Bounds(&sbuf, fixed.I(ppem), hinting)
	//	if err != nil {
	//		panic(err)
	//	}
	//	fmt.Printf("sbds: %#v\n", sbds)

	//fbds := ft.Bounds(fupe)
	//log.Printf("font: bbox: %v,%v,%v,%v", float64(fbds.Min.X), float64(fbds.Min.Y), float64(fbds.Max.X), float64(fbds.Max.Y))
	//	log.Printf("idx: %v", idx)
	// log.Printf("glyph.bbox: %#v", glyph.Bounds)
	//	log.Printf("glyph.bbox: x=%v, y=%v", float64(glyph.Bounds.Min.X), float64(glyph.Bounds.Min.Y))
	// log.Printf("glyph.adv: %#v, %v", glyph.AdvanceWidth, float64(glyph.AdvanceWidth)/64)
	//	log.Printf("symbol: %q, %q %d", symbolName, rn, rn)
	//	log.Printf("scale: %v", float64(scale))
	//	log.Printf("fupe:  %v", float64(fupe))
	//	log.Printf("slanted: %v", slanted)
	// log.Printf("font-size: %v", fontSize)

	face := truetype.NewFace(ft, &truetype.Options{
		DPI:     72,
		Size:    fontSize,
		Hinting: hinting,
	})
	bds, _ /*adv*/, ok := face.GlyphBounds(rn)
	if !ok {
		panic(fmt.Errorf("could not load glyph bounds for %q", rn))
	}

	log.Printf("box: %#v", bds)
	dy := bds.Max.Y + bds.Min.Y
	xmin := float64(bds.Min.X) / 64
	ymin := float64(bds.Min.Y-dy) / 64
	xmax := float64(bds.Max.X) / 64
	ymax := float64(bds.Max.Y-dy) / 64
	width := xmax - xmin
	height := ymax - ymin
	adv := ft.HMetric(fupe*64*6, idx).AdvanceWidth

	// FIXME(sbinet): for certain fonts (with postscript_name == "Cmex10")
	// offset = height/2 + (font.Size/3*dpi/72)
	offset := 0.0

	me := Metrics{
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
		postscript: ft.Name(truetype.NameIDPostscriptName),
		metrics:    me,
		symbol:     symbol,
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

func (ttf *TTFBackend) getGlyph(symbol string, font Font, math bool) (*truetype.Font, *sfnt.Font, rune, string, float64, bool) {
	var (
		fontType = font.Type
		idx      = unicodeIndex(symbol, math)
	)

	// only characters in the "Letter" class should be italicized in "it" mode.
	// Greek capital letters should be roman.
	if font.Type == "it" && idx < 0x10000 {
		if !unicode.Is(unicode.L, idx) {
			fontType = "rm"
		}
	}
	slanted := (fontType == "it") || ttf.isSlanted(symbol)
	ft, sft := ttf.getFont(fontType)
	if ft == nil {
		panic("could not find TTF font for [" + fontType + "]")
	}

	// FIXME(sbinet):
	// \sigma -> sigma, A->A, \infty->infinity, \nabla->gradient
	// etc...
	symbolName := symbol
	return ft, sft, idx, symbolName, font.Size, slanted
}

func (ttf *TTFBackend) isSlanted(symbol string) bool {
	switch symbol {
	case `\int`, `\oint`:
		return true
	default:
		return false
	}
}

func (ttf *TTFBackend) getFont(fontType string) (*truetype.Font, *sfnt.Font) {
	return ttf.fonts[fontType], ttf.sfnts[fontType]
}

type ttfKey struct {
	symbol string
	font   Font
	dpi    float64
}

type ttfVal struct {
	font       *truetype.Font
	size       float64
	postscript string
	metrics    Metrics
	symbol     string
	rune       rune
	glyph      truetype.Index
	offset     float64
}

var (
	_ Backend = (*TTFBackend)(nil)
)
