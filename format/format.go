// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package format implements standard formatting of LaTeX documents.
package format // import "github.com/go-latex/latex/format"

import (
	"go/token"
)

// Source formats src in canonical LaTeX style and returns the result or an
// (I/O or syntax) error.
// src is expected to be a syntactically correct LaTeX file.
func Source(src []byte) ([]byte, error) {
	var (
		fs = token.NewFileSet()
		f  = fs.AddFile("<input>", -1, len(src))
	)

	f.SetLinesForContent(src)

	return format(fs, f, src)
}

func format(fset *token.FileSet, file *token.File, src []byte) ([]byte, error) {
	return src, nil
}
