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
	"golang.org/x/image/font/opentype"
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

	var buf sfnt.Buffer

	ft, sft, rn, symbol, fontSize, slanted := ttf.getGlyph(symbol, font, math)
	idx := ft.Index(rn)

	gidx, err := sft.GlyphIndex(&buf, rn)
	if err != nil {
		panic(err)
	}
	symName, err := sft.GlyphName(&buf, gidx)
	if err != nil {
		panic(err)
	}
	log.Printf("symbol-name: %q", symName)

	fupe := fixed.Int26_6(ft.FUnitsPerEm())
	log.Printf("fupe: %v|%v | %v", fupe, fupe/2, scale(fixed.I(12), 2048))
	//	fupe = fixed.Int26_6(0.5 + (font.Size * 64))
	var glyph truetype.GlyphBuf
	err = glyph.Load(ft, fupe, idx, hinting)
	if err != nil {
		panic(err)
	}

	psname, err := sft.Name(&buf, sfnt.NameIDPostScript)
	if err != nil {
		panic(err)
	}
	var ppem = int(sft.UnitsPerEm() * 6)
	_, err = sft.LoadGlyph(&buf, gidx, fixed.I(ppem), nil)
	if err != nil {
		panic(err)
	}

	sadv, err := sft.GlyphAdvance(&buf, gidx, fixed.I(ppem), hinting)
	if err != nil {
		panic(err)
	}
	fmt.Printf("sadv: %v -> %v -> %v\n", sadv, float64(sadv), float64(sadv)/65536*fontSize/12)
	//	//	sme, err := sft.Metrics(&sbuf, fixed.I(ppem), hinting)
	//	//	if err != nil {
	//	//		panic(err)
	//	//	}
	//	//	fmt.Printf("sme: %#v\n", sme)
	//
	sfupe := fixed.Int26_6(sft.UnitsPerEm())
	_, err = sft.LoadGlyph(&buf, gidx, sfupe, nil)
	if err != nil {
		panic(err)
	}

	sbds, err := sft.Bounds(&buf, sfupe, hinting)
	if err != nil {
		panic(err)
	}
	fmt.Printf("sbds: %#v\n", sbds)

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

	fmt.Printf("box:  %#v\n", bds)
	dy := bds.Max.Y + bds.Min.Y
	xmin := float64(bds.Min.X) / 64
	ymin := float64(bds.Min.Y-dy) / 64
	xmax := float64(bds.Max.X) / 64
	ymax := float64(bds.Max.Y-dy) / 64
	width := xmax - xmin
	height := ymax - ymin
	adv := ft.HMetric(fupe*64*6, idx).AdvanceWidth

	//	_, err = sft.LoadGlyph(&buf, gidx, fixed.Int26_6(sft.UnitsPerEm()), nil)
	//	if err != nil {
	//		panic(err)
	//	}

	sface, err := opentype.NewFace(sft, &opentype.FaceOptions{
		DPI:     72,
		Size:    fontSize,
		Hinting: hinting,
	})
	if err != nil {
		panic(err)
	}
	defer sface.Close()
	{
		bds, _ /*adv*/, ok := sface.GlyphBounds(rn)
		if !ok {
			panic(fmt.Errorf("could not load glyph bounds for %q", rn))
		}
		fmt.Printf("sbox: %#v\n", bds)
		met, err := sft.GlyphMetrics(&buf, gidx, fixed.I(ppem), hinting)
		if err != nil {
			panic(err)
		}
		//fmt.Printf("smet: %#v\n", met)
		scale := font.Size / 12
		fmt.Printf("smet: adv=%v, h=%v, w=%v, x=%v,%v y=%v,%v\n",
			scale*float64(met.AdvanceX)/65536.,
			scale*float64(met.Height)/65536.,
			scale*float64(met.Width)/65536.,
			scale*float64(met.Bounds.Min.X)/65536.,
			scale*float64(met.Bounds.Max.X)/65536.,
			scale*float64(met.Bounds.Min.Y)/65536.,
			scale*float64(met.Bounds.Max.Y)/65536.,
		)
		{
			met, err := sft.GlyphMetrics(&buf, gidx, fixed.I(12), hinting)
			if err != nil {
				panic(err)
			}
			bnds := met.Bounds
			fmt.Printf("sbnds: %#v\n", bnds)
			fmt.Printf("smet: adv=%v, h=%v, w=%v, x=%v,%v y=%v,%v\n",
				scale*float64(met.AdvanceX)/64,
				scale*float64(met.Height)/64,
				scale*float64(met.Width)/64,
				scale*float64(met.Bounds.Min.X)/64,
				scale*float64(met.Bounds.Max.X)/64,
				scale*float64(met.Bounds.Min.Y)/64,
				scale*float64(met.Bounds.Max.Y)/64,
			)
		}
	}

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
		postscript: psname,
		metrics:    me,
		symbolName: symName,
		rune:       rn,
		glyph:      gidx,
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
	symbolName string
	rune       rune
	glyph      sfnt.GlyphIndex
	offset     float64
}

var (
	_ Backend = (*TTFBackend)(nil)
)
