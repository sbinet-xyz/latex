// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command latex-fmt formats LaTeX documents.
package main // import "github.com/go-latex/latex/cmd/latex-fmt"

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-latex/latex/format"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("latex-fmt: ")

	var (
		doW = flag.Bool("w", false, "write result to (source) file instead of stdout")
	)

	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatalf("missing input file(s)")
	}

	err := run(*doW, flag.Args())
	if err != nil {
		log.Fatalf("could not process files: %+v", err)
	}
}

func run(inplace bool, fnames []string) error {
	for _, fname := range fnames {
		err := run1(inplace, fname)
		if err != nil {
			return fmt.Errorf("could not format %q: %w", fname, err)
		}
	}
	return nil
}

func run1(inplace bool, fname string) error {

	src, err := ioutil.ReadFile(fname)
	if err != nil {
		return fmt.Errorf("could not read input file: %w", err)
	}

	o, err := format.Source(src)
	if err != nil {
		return fmt.Errorf("could not format input: %w", err)
	}

	switch inplace {
	case true:
		st, err := os.Stat(fname)
		if err != nil {
			return fmt.Errorf("could not stat input: %w", err)
		}
		err = ioutil.WriteFile(fname, o, st.Mode())
		if err != nil {
			return fmt.Errorf("could not write back formatted file: %w", err)
		}
	default:
		_, err = os.Stdout.Write(o)
		if err != nil {
			return fmt.Errorf("could not display formatted file: %w", err)
		}
	}

	return nil
}
